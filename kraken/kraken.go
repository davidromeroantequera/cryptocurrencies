package kraken

import (
	"time"
	"log"
	"strconv"
	"net/http"
	"encoding/json"
	"cryptocurrencies/types"
)

const sleep_time = 1000

type krakenTicker struct {
	Error []interface{}
	Result struct {
		Xethzusd struct {
			// This is the name of the currency pair
			A [3]string
			B [3]string
			C [3]string
			H []string
			L []string
			O string
			P []string
			T []int
			V []string
		}
	}
}

func retreiveAndUnmarshallTicker() (krakenTicker, error) {
	// TODO change default http client, because it doesn't handles timeouts
	resp, err := http.Get("https://api.kraken.com/0/public/Ticker?pair=ETHUSD")
	if err != nil {
		return krakenTicker{}, err
	}
	defer resp.Body.Close()

	t := new(krakenTicker)
	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		return krakenTicker{}, err
	}

	return *t, nil
}

func translateTicker(data krakenTicker) (types.Ticker, error) {
	t := types.Ticker{}
	var e error
	var f float64

	f, e = strconv.ParseFloat(data.Result.Xethzusd.A[0], 64) // Last
	if e != nil {
		return types.Ticker{}, e
	}
	t.Ask = f

	f, e = strconv.ParseFloat(data.Result.Xethzusd.B[0], 64) //Last
	if e != nil {
		return types.Ticker{}, e
	}
	t.Bid = f

	f, e = strconv.ParseFloat(data.Result.Xethzusd.H[1], 64) //Last 24 hours
	if e != nil {
		return types.Ticker{}, e
	}
	t.High = f

	f, e = strconv.ParseFloat(data.Result.Xethzusd.C[0], 64) //Last
	if e != nil {
		return types.Ticker{}, e
	}
	t.Value = f

	f, e = strconv.ParseFloat(data.Result.Xethzusd.L[1], 64) //Last 24 hours
	if e != nil {
		return types.Ticker{}, e
	}
	t.Low = f

	f, e = strconv.ParseFloat(data.Result.Xethzusd.V[1], 64) //Last 24 hours
	if e != nil {
		return types.Ticker{}, e
	}
	t.Volume = f

	f, e = strconv.ParseFloat(data.Result.Xethzusd.P[1], 64) //Last 24 hours
	if e != nil {
		return types.Ticker{}, e
	}
	t.Vwap = f

	t.Timestamp = time.Now()

	return t, nil
}

func NewKrakenTicker() chan types.Ticker {
	input := make(chan types.Ticker)
	go func() {
		for {
			time.Sleep(sleep_time * time.Millisecond)
			bt, err := retreiveAndUnmarshallTicker()
			if err != nil {
				log.Fatal(err)
				time.Sleep(sleep_time * time.Millisecond)
				continue
			}

			ticker, err := translateTicker(bt)
			if err != nil {
				log.Fatal(err)
				time.Sleep(sleep_time * time.Millisecond)
				continue
			}

			input <- ticker
		}
	}()
	return input
}
