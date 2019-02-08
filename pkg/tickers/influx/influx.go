package influx

import (
	"log"
	"time"
	"strconv"

	"cryptocurrencies/pkg/tickers"
	"cryptocurrencies/pkg/types"
	"cryptocurrencies/pkg/utils"

	"github.com/influxdata/influxdb/client/v2"
	"encoding/json"
)



func createQuery(measurement string, timeLapse time.Duration) string {
	const (
		Q                = "\""
		queryBase        = "SELECT ask, bid, high, low, volume, vwpa, value FROM "
		queryConditional = " WHERE time > "
	)

	dataSource := Q + measurement + Q
	since := time.Now().Add(-timeLapse).Unix()

	return queryBase + dataSource + queryConditional + strconv.FormatInt(since, 10)
}

// TODO we should use the marshaller interface here instead of all this mess
func translateRecordsToTicker(res []client.Result) ([]types.Ticker, error) {
	if len(res) == 0 || len(res[0].Series) == 0{
		return nil, nil
	}

	recordset := res[0].Series[0].Values
	dataSet := make([]types.Ticker, len(recordset))

	for i, row := range recordset {
		dataSet[i] = types.Ticker{}
		var value float64
		var timestamp time.Time
		var err error

		if timestamp, err = time.Parse(time.RFC3339, row[0].(string)); err != nil {
			return nil, err
		}
		dataSet[i].Timestamp = timestamp

		if value, err = row[1].(json.Number).Float64(); err != nil {
			return nil, err
		}
		dataSet[i].Ask = value

		if value, err = row[2].(json.Number).Float64(); err != nil {
			return nil, err
		}
		dataSet[i].Bid = value

		if value, err = row[3].(json.Number).Float64(); err != nil {
			return nil, err
		}
		dataSet[i].High = value

		if value, err = row[4].(json.Number).Float64(); err != nil {
			return nil, err
		}
		dataSet[i].Low = value

		if value, err = row[5].(json.Number).Float64(); err != nil {
			return nil, err
		}
		dataSet[i].Volume = value

		if value, err = row[6].(json.Number).Float64(); err != nil {
			return nil, err
		}
		dataSet[i].Vwap = value

		if value, err = row[1].(json.Number).Float64(); err != nil {
			return nil, err
		}
		dataSet[i].Value = value
	}

	return dataSet, nil
}

func retrieveTickerSet(dbName string, measurement string, timeLapse time.Duration) ([]types.Ticker, error) {
	cmd := createQuery(measurement, timeLapse)
	result, err := utils.QueryDB(cmd, dbName)
	if err != nil {
		return nil, err
	}

	return translateRecordsToTicker(result)
}

// TODO change the interface for allowing the specification of a time frame
func NewInfluxTicker(dbName string, measurement string, timeLapse time.Duration) (tickers.TickerChan, types.StopChannel) {
	input := make(tickers.TickerChan)
	stop := make(types.StopChannel)

	go func() {
		tickerSet, err := retrieveTickerSet(dbName, measurement, timeLapse)
		if err != nil {
			log.Fatal("An error ocurred while retrieving the tickers from DB")
		}

		for _, ticker := range tickerSet {
			input <- ticker
		}
		utils.TerminateChannels(stop)
	}()

	return input, stop
}
