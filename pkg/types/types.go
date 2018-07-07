package types

import "time"

type Ticker struct {
	Ask       float64
	Bid       float64
	High      float64
	Low       float64
	Volume    float64
	Vwap      float64
	Value     float64
	Timestamp time.Time
}

type StopChannel = chan interface{}