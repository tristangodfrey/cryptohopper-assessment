package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"os"
)

// Some constants that serve as the signal return values
const SignalSell = "SELL"
const SignalBuy = "BUY"
const SignalNeutral = "NEUTRAL"

const MinPriceDataLength = 55

// SMA - Simple interface for structs that are able to generate an SMA(n), with an offset to determine which timepoint
// should be considered the end of the sample. (e.g. if offset = 3, the SMA will run up until t-3)
type SMA interface {
	sma(n int, offset int) float64
	isValid() bool
}

// Struct to hold data we receive from the Cryptohopper API (other fields are ignored for the purposes of this assessment)
type PriceData struct {
	Close float64 `json:"Close"`
}

type PriceDataList []PriceData

// Ensures that we have enough data (such that we don't run into index out of bounds issues later on)
func (pdl PriceDataList) isValid() bool {
	return len(pdl) >= MinPriceDataLength
}

// Note that we implement the sma method on the PriceDataList, as opposed to some "generic" value list, since we
// would then have to do a conversion first (at which point we have access to each value needed to calculate
// the SMA), before doing the same thing again to perform the calculation
func (pdl PriceDataList) sma(n int, offset int) float64 {
	start := len(pdl) - n - offset
	end := len(pdl) - offset

	slice := pdl[start:end]
	sum := float64(0)

	for _, pd := range slice {
		sum += pd.Close
	}

	return sum / float64(n)
}

// Parses a period string (5m, 2h, 1d, etc) and converts it to a time.Duration
func parsePeriod(period string) (time.Duration, error) {
	switch period {
	case "1m":
		return time.Minute, nil
	case "5m":
		return time.Minute * 5, nil
	case "15m":
		return time.Minute * 15, nil
	case "30m":
		return time.Minute * 30, nil
	case "1h":
		return time.Hour, nil
	case "2h":
		return time.Hour * 2, nil
	case "4h":
		return time.Hour * 4, nil
	case "1d":
		return time.Hour * 24, nil
	default:
		return time.Duration(0), errors.New(fmt.Sprintf("Invalid period %s", period))
	}
}

// Fetches price data from various APIs through the cryptohopper API
func fetchPriceData(exchange string, pair string, period string) (PriceDataList, error) {
	periodDuration, err := parsePeriod(period)
	pdl := PriceDataList{}

	if err != nil {
		return pdl, err
	}

	//multiply the period x 55 to get a sufficient sample size for SMA(55)
	endTime := time.Now()
	startTime := endTime.Add(-(periodDuration * 55))

	url := fmt.Sprintf(
		"http://cryptohopper-ticker-frontend.us-east-1.elasticbeanstalk.com/v1/%s/candles?pair=%s&start=%s&end=%s&period=%s",
		exchange,
		pair,
		strconv.FormatInt(startTime.Unix(), 10),
		strconv.FormatInt(endTime.Unix(), 10),
		period,
	)

	res, err := http.Get(url)

	if err != nil {
		return pdl, err
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&pdl)

	if err != nil {
		return pdl, err
	}

	return pdl, nil
}

// Compares SMA(8) at t and t-1 with SMA(55) to determine which signal should be sent
func generateSignal(s SMA) string {
	currentSma8 := s.sma(8, 0)
	prevSma8 := s.sma(8, 1)
	currentSma55 := s.sma(55, 0)

	prevHigherOrSame := prevSma8 >= currentSma55
	prevLowerOrSame := prevSma8 <= currentSma55

	currentLower := currentSma8 < currentSma55
	currentHigher := currentSma8 > currentSma55

	if currentLower && prevHigherOrSame {
		return SignalSell
	}

	if currentHigher && prevLowerOrSame {
		return SignalBuy
	}

	return SignalNeutral
}

// Main application route
func signal(w http.ResponseWriter, r *http.Request) {
	exchange := r.URL.Query().Get("exchange")
	pair := r.URL.Query().Get("pair")
	period := r.URL.Query().Get("period")

	pdl, err := fetchPriceData(exchange, pair, period)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	if ! pdl.isValid() {
		w.WriteHeader(500)
		fmt.Fprint(w, "Insufficient data available for the requested time range")
		return
	}

	signal := generateSignal(pdl)

	fmt.Fprint(w, signal)
}

func main() {
	http.HandleFunc("/", signal)

	port := os.Getenv("PORT")

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		fmt.Println(err)
	}
}
