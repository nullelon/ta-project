package app

import (
	"encoding/json"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"ta-project-go/internal/model"
	"time"
)

type Subscription struct {
	symbol    string
	timeframe string
	time      time.Time
}

type MarkerService struct {
	ws         *websocket.Conn
	wsMutex    *sync.Mutex
	updateChan chan BinanceUpdate

	subs *queue.Queue

	cache map[string]map[string][]model.Candle // hash-table + ArrayList of
}

func NewMarketService() *MarkerService {
	return &MarkerService{
		wsMutex:    new(sync.Mutex),
		updateChan: make(chan BinanceUpdate),
		cache:      make(map[string]map[string][]model.Candle, 20),
		subs:       queue.New(),
	}
}

func (s *MarkerService) OpenConnection() error {
	c, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/stream1", nil)
	if err != nil {
		return err
	}
	s.ws = c
	go func() {
		for {
			var upd BinanceUpdate
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				continue
			}
			json.Unmarshal(msg, &upd)
			s.updateChan <- upd
		}
	}()

	go func() {
		for {
			update := <-s.updateChan

			switch update.EventType {
			case "kline":
				candles := s.cache[update.Symbol][update.Kline.Interval]
				lastCandle := candles[len(candles)-1]
				if lastCandle.OpenTime == update.Kline.OpenTime {
					candles[len(candles)-1] = klineToCandle(update.Kline)
				} else {
					candles = append(candles, klineToCandle(update.Kline))
				}
			}
		}
	}()

	go func() {
		for {
			err := s.sendToWebSocket(WebsocketRequest{
				Method: "LIST_SUBSCRIPTIONS",
				Id:     55,
			})
			if err != nil {
				return
			}
			<-time.After(time.Second)
		}
	}()

	return nil
}

type WebsocketRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

func (s MarkerService) sendToWebSocket(request WebsocketRequest) error {
	marshal, _ := json.Marshal(request)
	s.wsMutex.Lock()
	defer s.wsMutex.Unlock()
	return s.ws.WriteMessage(websocket.TextMessage, marshal)
}

func (s *MarkerService) Subscribe(symbol string, timeframe string) {
	fmt.Printf("Subscribing to %s:%s\n", strings.ToUpper(symbol), timeframe)
	err := s.sendToWebSocket(WebsocketRequest{
		Method: "SUBSCRIBE",
		Params: []string{fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), timeframe)},
		Id:     1,
	})
	if err != nil {
		return
	}

	s.subs.Enqueue(Subscription{
		symbol:    symbol,
		timeframe: timeframe,
		time:      time.Now(),
	})
}

func (s *MarkerService) Remove(symbol string, timeframe string) {
	s.Unsubscribe(symbol, timeframe)
	delete(s.cache[symbol], timeframe)

	if len(s.cache[symbol]) == 0 {
		delete(s.cache, symbol)
	}
}

func (s *MarkerService) Unsubscribe(symbol string, timeframe string) {
	fmt.Printf("Unsubscribing from %s:%s\n", strings.ToUpper(symbol), timeframe)
	err := s.sendToWebSocket(WebsocketRequest{
		Method: "UNSUBSCRIBE",
		Params: []string{fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), timeframe)},
		Id:     1,
	})
	if err != nil {
		return
	}
}

func (s *MarkerService) Get(symbol, timeframe string, limit int) ([]model.Candle, error) {
	timeframes, ok := s.cache[symbol]
	if !ok {
		klines, err := binanceKlines(symbol, timeframe, 1000)
		if err != nil {
			return nil, err
		}
		go func() {
			s.cache[symbol] = map[string][]model.Candle{
				timeframe: klines,
			}
			s.Subscribe(symbol, timeframe)
		}()
		return klines[:limit], nil
	}
	candles, ok := timeframes[timeframe]
	if !ok {
		klines, err := binanceKlines(symbol, timeframe, 1000)
		if err != nil {
			return nil, err
		}
		go func() {
			s.cache[symbol][timeframe] = klines
			s.Subscribe(symbol, timeframe)
		}()
		return klines[:limit], nil
	}
	return candles[:limit], nil
}

func (s *MarkerService) Start() error {
	err := s.OpenConnection()
	if err != nil {
		return err
	}

	go func() {
		for {
			if s.subs.Len() != 0 {
				sub := s.subs.Peek().(Subscription)
				if time.Since(sub.time).Seconds() > 20 {
					s.subs.Dequeue()
					s.Remove(sub.symbol, sub.timeframe)
				} else {
					<-time.After(time.Second)
				}
			} else {
				<-time.After(time.Second)
			}
		}
	}()
	return nil
}
