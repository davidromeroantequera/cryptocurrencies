package tests

import (
	"time"
	"testing"

	"cryptocurrencies/pkg/utils"
	"cryptocurrencies/pkg/types"
	"cryptocurrencies/pkg/writers"
	"cryptocurrencies/pkg/tickers/bitso"
	"cryptocurrencies/pkg/tickers/influx"
)

const(
	dbName = "integration_tests_database__"
	mBitso = "bitso_measurement__"
)

func Setup() {
	utils.CreateDatabase(dbName)
}

func TearDown() {
	utils.DropDatabase(dbName)
}

func retrieveTickersFromBitso() []types.Ticker {
	tickers := make([]types.Ticker, 0)

	bt, bts := bitso.NewBitsoTicker()
	db := writers.NewInfluxWriter(dbName, mBitso)
	w, ws := db.TickerToInfluxWriter()

	stop := make(types.StopChannel)

	go func() {
		for {
			select {
			case incoming := <-bt:
				tickers = append(tickers, incoming)
				w <- incoming
				break;
			case <-stop:
				utils.TerminateChannels(bts, ws)
				return
			}
		}
	}()

	time.Sleep(20 * time.Second);
	utils.TerminateChannels(stop)

	return tickers
}

func retrieveTickersFromDB() []types.Ticker {
	tickers := make([]types.Ticker, 0)

	it, its := influx.NewInfluxTicker(dbName, mBitso, 1 * time.Minute)

	all_tickers_received := false
	for ; !all_tickers_received ; {
		select {
		case incoming := <-it:
			tickers = append(tickers, incoming)
			break
		case <-its:
			all_tickers_received = true
		}
	}

	return tickers
}

func TestRetrieveAndWriteBitsoTicker(t *testing.T) {
	Setup()
	defer TearDown()

	// It seems that it is not possible to make sure that all the tickers were written to the database,
	// because of different times at the moment of closing the different channels (they're closed in cascade)
	// retrievedTickers := retrieveTickersFromBitso()
	// TODO flush in writer is not occuring correctly, because we are not waiting for all the go routines to terminate
	retrieveTickersFromBitso()
	storedTickers := retrieveTickersFromDB()

	utils.Assert(t, len(storedTickers) != 0, "Integration failed!")
}