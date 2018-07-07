package utils

import "cryptocurrencies/types"


func SendTickerToAllChannels(t types.Ticker, d... chan<- types.Ticker) {
	for _, item := range d {
		item <- t
	}
}
