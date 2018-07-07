package bitso

import (
	"cryptocurrencies/types"
	"time"
	"net/http"
	"encoding/json"
	"strconv"
	"log"
)

type Ticker = types.Ticker

type bitsoTicker struct {
	Success bool
	Payload struct {
		Volume     string
		High       string
		Last       string
		Low        string
		Vwap       string
		Ask        string
		Bid        string
		Created_at string
	}
}

func retreiveAndUnmarshallTicker() (bitsoTicker, error) {
	resp, err := http.Get("https://api.bitso.com/v3/ticker/?book=eth_mxn")
	if err != nil {
		return bitsoTicker{}, err
	}
	defer resp.Body.Close()

	t := new(bitsoTicker)
	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		return bitsoTicker{}, err
	}

	return *t, nil
}

func bitsoTickerToTicker(bt bitsoTicker) (Ticker, error) {
	t := Ticker{}
	var e error
	var f float64

	f, e = strconv.ParseFloat(bt.Payload.Ask, 64)
	if e != nil {
		return Ticker{}, e
	}
	t.Ask = f

	f, e = strconv.ParseFloat(bt.Payload.Bid, 64)
	if e != nil {
		return Ticker{}, e
	}
	t.Bid = f

	f, e = strconv.ParseFloat(bt.Payload.High, 64)
	if e != nil {
		return Ticker{}, e
	}
	t.High = f

	f, e = strconv.ParseFloat(bt.Payload.Last, 64)
	if e != nil {
		return Ticker{}, e
	}
	t.Value = f

	f, e = strconv.ParseFloat(bt.Payload.Low, 64)
	if e != nil {
		return Ticker{}, e
	}
	t.Low = f

	f, e = strconv.ParseFloat(bt.Payload.Volume, 64)
	if e != nil {
		return Ticker{}, e
	}
	t.Volume = f

	f, e = strconv.ParseFloat(bt.Payload.Vwap, 64)
	if e != nil {
		return Ticker{}, e
	}
	t.Vwap = f

	var timestamp time.Time
	timestamp, e = time.Parse("2006-01-02T15:04:05+00:00", bt.Payload.Created_at)
	if e != nil {
		return Ticker{}, e
	}
	t.Timestamp = timestamp

	return t, nil
}

func NewBitsoTicker() chan Ticker {
	input := make(chan Ticker)
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			bt, err := retreiveAndUnmarshallTicker()
			if err != nil {
				log.Fatal("An error ocurred while retrieving ticker from cryptocurrencies: ", err)
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			ticker, err := bitsoTickerToTicker(bt)
			if err != nil {
				log.Fatal("An error ocurred while unmarshalling ticker from cryptocurrencies: ", err)
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			input <- ticker
		}
	}()
	return input
}

