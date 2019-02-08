package utils

import (
	"testing"

	"cryptocurrencies/pkg/types"
)


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

func Assert(t *testing.T, condition bool, message string) {
	if !condition {
		t.Fatal(message)
	}
}