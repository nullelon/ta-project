package model

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Candle struct {
	OpenTime  int64   `json:"open_time"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	CloseTime int64   `json:"close_time"`
}

func (c *Candle) UnmarshalJSON(b []byte) error {
	var stuff []interface{}
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	if len(stuff) != 12 {
		return errors.New("wrong data, expected 12 elements")
	}
	c.OpenTime = int64(stuff[0].(float64))
	c.Open, _ = strconv.ParseFloat(stuff[1].(string), 64)
	c.High, _ = strconv.ParseFloat(stuff[2].(string), 64)
	c.Low, _ = strconv.ParseFloat(stuff[3].(string), 64)
	c.Close, _ = strconv.ParseFloat(stuff[4].(string), 64)
	c.Volume, _ = strconv.ParseFloat(stuff[5].(string), 64)
	c.CloseTime = int64(stuff[6].(float64))
	return nil
}
