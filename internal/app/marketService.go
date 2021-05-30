package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"ta-project-go/internal/model"
)

type Subscription struct {
	symbol    string
	timeframe string
}

type MarkerService struct {
	ws         *websocket.Conn
	updateChan chan BinanceUpdate

	subs    []Subscription
	subsIds map[Subscription]int

	cache map[string]map[string][]model.Candle
}

func NewMarketService() *MarkerService {
	return &MarkerService{
		updateChan: make(chan BinanceUpdate),
		cache:      make(map[string]map[string][]model.Candle, 20),
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
			err := c.ReadJSON(&upd)
			if err != nil {
				fmt.Println(err)
				continue
			}
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
	return nil
}

func klineToCandle(kline Kline) model.Candle {
	open, _ := strconv.ParseFloat(kline.Open, 64)
	low, _ := strconv.ParseFloat(kline.Low, 64)
	high, _ := strconv.ParseFloat(kline.High, 64)
	closeFloat, _ := strconv.ParseFloat(kline.Close, 64)
	return model.Candle{
		OpenTime:  kline.OpenTime,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     closeFloat,
		CloseTime: kline.CloseTime,
	}
}

type WebsocketRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

func (s *MarkerService) Subscribe(symbol string, timeframe string) {
	marshal, _ := json.Marshal(WebsocketRequest{
		Method: "SUBSCRIBE",
		Params: []string{fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), timeframe)},
		Id:     1,
	})
	err := s.ws.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		return
	}
}

func (s *MarkerService) Unsubscribe(symbol string, timeframe string) {
	marshal, _ := json.Marshal(WebsocketRequest{
		Method: "UNSUBSCRIBE",
		Params: []string{fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), timeframe)},
		Id:     1,
	})
	err := s.ws.WriteMessage(websocket.TextMessage, marshal)
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
	return nil
}
