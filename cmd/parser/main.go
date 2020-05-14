package main

import (
	"flag"
	"fmt"
	"github.com/rejlersembriq/gs2"
	"log"
	"os"
	"time"
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

	_, err = gs2.NewDecoder(file, options...).Decode()
	if err != nil {
		log.Fatal(err)
	}

	took := time.Since(start)

	fmt.Printf("Parsing took: %v\n", took)
}
