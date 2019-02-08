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
	"cryptocurrencies/pkg/writers"
	"cryptocurrencies/pkg/types"
)

func NewTickerRetriever(ctor tickers.Ctor,	dbName string,	measurement string,	wg *sync.WaitGroup) types.StopChannel {
	wg.Add(1)
	t,s := ctor()

	w := writers.NewInfluxWriter(dbName, measurement)
	db, dbs := w.TickerToInfluxWriter()
	stop := make(types.StopChannel)

	log.Printf("Initializing ticker retriever for %s (%s)\n", measurement, dbName)

	go func() {
		defer wg.Done()
		log.Printf("Sampling ticker %s\n", measurement)
		for {
			select {
			case tick := <-t:
				db <- tick
			case <-stop:
				log.Printf("Stopping ticker retriever %s (%s)\n", measurement, dbName)
				utils.TerminateChannels(s, dbs)
				return
			}
		}
	}()

	return stop
}

func main() {
	// TODO factory for the ticker constructors
	// TODO load tickers from configuration file
	// TODO arguments for
	var wg sync.WaitGroup
	stopKraken := NewTickerRetriever(kraken.NewKrakenTicker, "kraken", "eth_usd", &wg)
	stopBitso:= NewTickerRetriever(bitso.NewBitsoTicker, "bitso", "eth_mxn", &wg)

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
