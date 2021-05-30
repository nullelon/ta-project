package app

import (
	"strconv"
	"ta-project-go/internal/model"
)

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
