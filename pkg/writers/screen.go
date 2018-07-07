package writers

import (
	"fmt"
)

func TickerToScreenWriter() chan<- Ticker {
	input := make(chan Ticker)

	go func() {
		for {
			fmt.Printf("Ticker: %+v\n", <-input)
		}
	}()

	return input
}
