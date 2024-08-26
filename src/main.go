package main

import (
	"fmt"
	"github.com/pgwhalen/conflate-talk/src/solutions"
	"time"
)

// initialize time to enabled "normalized" output
var initTime = time.Now()
var normalizedInit = time.Date(initTime.Year(), initTime.Month(), initTime.Day(), 0, 0, 0, 0, initTime.Location())
var normalizeDuration = initTime.Sub(normalizedInit)

const timestampFormat = "05.0s"

var currDollar float64 = 0
var currTimestamp = initTime

var totalSales float64

type Sale struct {
	dollars   float64
	timestamp time.Time
}

func (s Sale) ConflateWith(s2 Sale) Sale {
	return Sale{
		dollars:   s.dollars + s2.dollars,
		timestamp: maxTime(s.timestamp, s2.timestamp),
	}
}

func (s Sale) ZeroValue() Sale {
	return Sale{}
}

func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func (s Sale) String() string {
	timestamp := s.timestamp
	normalizedTimestamp := normalizeTime(timestamp)

	format := normalizedTimestamp.Format(timestampFormat)
	return fmt.Sprintf("$%.2f at %s", s.dollars, format)
}

func main() {
	// setup
	//salesCh := make(chan Sale)
	//inCh := salesCh
	//outCh := salesCh

	// solution 0
	//salesCh := make(chan Sale, 3)
	//inCh := salesCh
	//outCh := salesCh

	// solution 1
	//inCh, outCh := solutions.ConflateV1[Sale]()

	// solution 2
	//inCh, outCh := solutions.ConflateV2[Sale](time.Millisecond)

	// solution 3
	inCh, outCh := solutions.ConflateV3[Sale](time.Millisecond)

	go producer(inCh)
	go consumer(outCh)

	select {}
}

func producer(salesCh chan<- Sale) {
	for {
		salesCh <- makeSale()

		// solution 0b
		//select {
		//case salesCh <- makeSale():
		//default:
		//}
	}
}

func consumer(salesCh <-chan Sale) {
	for s := range salesCh {
		redrawUI(s)
	}
}

func makeSale() Sale {
	const saleInterval = 500 * time.Millisecond
	time.Sleep(saleInterval)
	currDollar = currDollar + 1
	currTimestamp = currTimestamp.Add(saleInterval)
	sale := Sale{dollars: float64(currDollar), timestamp: currTimestamp}
	nLog("Made sale: %s", sale)
	return sale
}

func redrawUI(sale Sale) {
	totalSales += sale.dollars
	lagMillis := time.Now().Sub(sale.timestamp).Milliseconds()
	const redrawUITime = 700 * time.Millisecond
	time.Sleep(redrawUITime)
	nLog("Sale received after %dms: %s. Total sales: $%.2f", lagMillis, sale, totalSales)
}

func normalizeTime(timestamp time.Time) time.Time {
	return timestamp.Add(-normalizeDuration)
}

func nLog(format string, v ...any) {
	realTs := time.Now()
	normalizedTs := normalizeTime(realTs)
	fmtTs := normalizedTs.Format(timestampFormat)
	fmt.Printf(fmtTs+" | "+format+"\n", v...)
}
