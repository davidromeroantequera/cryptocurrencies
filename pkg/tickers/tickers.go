package tickers

import (
	"cryptocurrencies/pkg/types"
)

type TickerChan = chan types.Ticker
type Ctor func() (TickerChan, types.StopChannel)
