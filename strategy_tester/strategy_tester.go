package main

import (
	"time"
	"cryptocurrencies/pkg/tickers/influx"
	"cryptocurrencies/pkg/writers"
	"cryptocurrencies/pkg/utils"
)

func main() {
	const sourceDB = "kraken"
	const sourceMeasurement = "eth_usd"

	const targetDB = "strategies"
	const targetMeasurement = "eth_usd"

	it, _ := influx.NewInfluxTicker(sourceDB, sourceMeasurement, 12*time.Hour)

	utils.DropDatabase(targetDB)
	utils.CreateDatabase(targetDB)

	w := writers.NewInfluxWriter(targetDB, targetMeasurement)
	tw, tws := w.TickerToInfluxWriter()

	go func() {
		for {
			t := <-it
			tw <- t
		}
	}()

	time.Sleep(8*time.Second)
	utils.TerminateChannels(tws)

	return
}
