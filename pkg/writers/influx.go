package writers

import (
	"log"
	"cryptocurrencies/pkg/types"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	highWaterMark = 100
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

func (w InfluxWriter) flushToDatabase(points []*client.Point) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  w.db,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoints(points)

	if err := w.c.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func (w InfluxWriter) TickerToInfluxWriter() (chan<- Ticker, types.StopChannel) {
	input := make(chan Ticker)
	stop:= make(types.StopChannel)

	tags := w.Tags

	go func() {
		count := 0

		points := []*client.Point{}

		for {
			var t types.Ticker
			select {
				case t = <-input:
					break
				case <-stop:
					w.flushToDatabase(points)
					return
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
				log.Printf("TickerToInfluxWriter: %v\n", err)
				continue
			}
			points = append(points, pt)

			count++

			if count > highWaterMark {
				w.flushToDatabase(points)
				points = []*client.Point{}
				count = 0
			}
		}
	}()

	return input, stop
}
