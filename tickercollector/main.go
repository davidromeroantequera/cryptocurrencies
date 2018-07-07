package main

import (
	"os"
	"log"
	"sync"
	"syscall"

	"os/signal"

	"cryptocurrencies/pkg/utils"
	"cryptocurrencies/pkg/tickers"
	"cryptocurrencies/pkg/tickers/kraken"
	"cryptocurrencies/pkg/tickers/bitso"
)


func main() {
	// TODO factory for the ticker constructors
	// TODO load tickers from configuration file
	// TODO arguments for
	var wg sync.WaitGroup
	stopKraken := tickers.NewTickerRetriever(kraken.NewKrakenTicker, "kraken", "eth_usd", &wg)
	stopBitso:= tickers.NewTickerRetriever(bitso.NewBitsoTicker, "bitso", "eth_mxn", &wg)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v\n", sig)
		log.Println("Waiting for ticker retrievers to finalize...")
		utils.TerminateChannels(stopKraken, stopBitso)
		wg.Wait()
		os.Exit(0)
	}()

	select{}
}
