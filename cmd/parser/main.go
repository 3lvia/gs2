package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/3lvia/gs2"
)

var filename = flag.String("file", "", "File to parse")

func validateTime(g *gs2.GS2) error {
	for _, ts := range g.TimeSeries {
		if ts.Start.Add(time.Duration(ts.NoOfValues)*ts.Step) != ts.Stop {
			return fmt.Errorf("start %q, stop %q step %q doesnt match", ts.Start.Format(time.RFC3339), ts.Stop.Format(time.RFC3339), ts.Step)
		}
	}
	return nil
}

func validateMeterReadings(g *gs2.GS2) bool {
	for _, mr := range g.MeterReadings {
		if mr.DirectionOfFlow == "out" || mr.DirectionOfFlow == "in" {
			return true
		}
	}

	return false
}

func main() {
	flag.Parse()

	if *filename == "" {
		log.Fatal("Need to speficy filename with -file=<filename>")
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	start := time.Now()

	options := []gs2.DecoderOption{
		gs2.DecodeValidators(
			gs2.ValidateNoOfObjects,
			gs2.ValidateTimeSeriesValues,
			validateTime,
		),
	}

	var result = &gs2.GS2{}

	result, err = gs2.NewDecoder(file, options...).Decode()

	if validateMeterReadings(result) {
		fmt.Println("Meterreading has direction of flow")
	}

	if err != nil {
		log.Fatal(err)
	}

	took := time.Since(start)

	fmt.Printf("Parsing took: %v\n", took)
}
