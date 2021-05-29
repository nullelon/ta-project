package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"ta-project-go/internal/model"
)

type MarkerService struct {
	cache map[string]map[string][]model.Candle
}

func NewMarkerService() *MarkerService {
	return &MarkerService{}
}

func (MarkerService) Get(symbol, timeframe string, limit int) ([]model.Candle, error) {
	return binanceKlines(symbol, timeframe, limit)
}

func binanceKlines(symbol, timeframe string, limit int) ([]model.Candle, error) {
	var candles []model.Candle

	return candles, binanceMethod(&candles, "klines", map[string]string{
		"symbol":   symbol,
		"interval": timeframe,
		"limit":    strconv.Itoa(limit),
	})
}

func binanceMethod(dest interface{}, method string, args map[string]string) error {
	const apiUrl = "https://api.binance.com/api/v3/%s?%s"

	urlValues := url.Values{}
	for key, val := range args {
		urlValues.Add(key, val)
	}

	resp, err := new(http.Client).Get(fmt.Sprintf(apiUrl, method, urlValues.Encode()))
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(dest)
	if err != nil {
		return err
	}
	return nil
}
