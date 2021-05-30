package app

type BinanceUpdate struct {
	EventType string `json:"e"`
	EventTime int    `json:"E"`
	Symbol    string `json:"s"`
	Kline     Kline  `json:"k"`
}

type Kline struct {
	OpenTime       int64  `json:"t"`
	CloseTime      int64  `json:"T"`
	Symbol         string `json:"s"`
	Interval       string `json:"i"`
	FirstTradeId   int    `json:"f"`
	LastTradeId    int    `json:"L"`
	Open           string `json:"o"`
	Close          string `json:"c"`
	High           string `json:"h"`
	Low            string `json:"l"`
	V              string `json:"v"`
	NumberOfTrades int    `json:"n"`
	IsClosed       bool   `json:"x"`
	Volume         string `json:"q"`
	V1             string `json:"V"`
	Q1             string `json:"Q"`
	B              string `json:"B"`
}
