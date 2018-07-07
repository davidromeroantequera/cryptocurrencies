package main

import (
	"cryptocurrencies/writers"
	"cryptocurrencies/utils"
	"cryptocurrencies/kraken"
	"cryptocurrencies/bitso"
)

func main() {
	k_data := kraken.NewKrakenTicker()
	b_data := bitso.NewBitsoTicker()

	w_kraken := writers.NewInfluxWriter("kraken", "eth_usd")
	db_kraken := w_kraken.TickerToInfluxWriter()

	w_bitso := writers.NewInfluxWriter("bitso", "eth_mxn")
	db_bitso := w_bitso.TickerToInfluxWriter()

	screen := writers.TickerToScreenWriter()

	go func() {
		for {
			s := <-k_data
			utils.SendTickerToAllChannels(s, db_kraken, screen)
		}
	}()

	for {
		db_bitso<- <-b_data
	}
}
