package utils

import "cryptocurrencies/pkg/types"


func SendTickerToAllChannels(t types.Ticker, d... chan<- types.Ticker) {
	for _, item := range d {
		item <- t
	}
}

func TerminateChannels(channels... types.StopChannel) {
	for _, channel := range channels {
		channel <- 0
	}
}