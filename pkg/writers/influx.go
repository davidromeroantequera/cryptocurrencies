package writers

import (
	"log"
	"cryptocurrencies/pkg/types"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	highWaterMark = 10
)

type Ticker = types.Ticker

type InfluxWriter struct {
	c client.Client
	db string
	measurement string
	Tags map[string]string
}

func NewInfluxWriter(dbName string, measurement string) InfluxWriter {
	w := InfluxWriter{}
	var err error

	w.c, err = client.NewHTTPClient(client.HTTPConfig{Addr: "http://localhost:8086"})
	if err != nil {
		log.Fatal(err)
		return InfluxWriter{}
	}

	w.db = dbName
	w.measurement = measurement

	return w
}

func (w InfluxWriter) TickerToInfluxWriter() chan<- Ticker {
	input := make(chan Ticker)
	tags := w.Tags

	go func() {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  w.db,
			Precision: "s",
		})
		if err != nil {
			log.Fatal(err)
		}

		count := 0

		for {
			t := <-input

			if err != nil {
				log.Fatal(err)
			}

			fields := map[string]interface{}{
				"ask":    t.Ask,
				"bid":    t.Bid,
				"high":   t.High,
				"low":    t.Low,
				"volume": t.Volume,
				"value":  t.Value,
				"vwpa":   t.Vwap,
			}

			pt, err := client.NewPoint(w.measurement, tags, fields, t.Timestamp)
			if err != nil {
				log.Fatal(err)
			}
			bp.AddPoint(pt)

			count++

			if count > highWaterMark {
				if err := w.c.Write(bp); err != nil {
					log.Fatal(err)
				}

				bp, err = client.NewBatchPoints(client.BatchPointsConfig{
					Database:  w.db,
					Precision: "s",
				})
				if err != nil {
					log.Fatal(err)
				}

				count = 0
			}
		}
	}()

	return input
}
