package tickers

import (
	"cryptocurrencies/pkg/types"
	"sync"
	"cryptocurrencies/pkg/writers"
	"log"
)

type TickerChan = chan types.Ticker
type Ctor func() TickerChan

func NewTickerRetriever(ctor Ctor,	dbName string,	measurement string,	wg *sync.WaitGroup) types.StopChannel {
	wg.Add(1)
	t := ctor()
	w := writers.NewInfluxWriter(dbName, measurement)
	db := w.TickerToInfluxWriter()
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
				log.Printf("Stopping ticker retriever %s (%s)", measurement, dbName)
				return
			}
		}
	}()

	return stop
}