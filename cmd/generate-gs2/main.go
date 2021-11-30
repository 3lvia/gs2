package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/3lvia/gs2"
)

const outputFile = "testdata/testFile.gs2"
const maxVal = 1000

func main() {
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())

	startMessage := getStartMessage()
	meterReadings := generateMeterReadings(0)
	timeSeries := generateTimeSeries(1)
	endMessage := getEndMessage(len(meterReadings) + len(timeSeries) + 2)

	g := gs2.GS2{
		StartMessage:  startMessage,
		MeterReadings: meterReadings,
		TimeSeries:    timeSeries,
		EndMessage:    endMessage,
	}

	if err := gs2.NewEncoder(file).Encode(&g); err != nil {
		log.Fatal(err)
	}
}

func getStartMessage() gs2.StartMessage {
	now := time.Now()
	_, offset := now.Zone()
	return gs2.StartMessage{
		ID:           "someId",
		MessageType:  "Settlement-data",
		Version:      "1.2",
		Time:         now,
		To:           "recipient",
		From:         "sender",
		GMTReference: offset / 60 / 60,
	}
}

func getEndMessage(noOfObjects int) gs2.EndMessage {
	return gs2.EndMessage{
		ID:              "someId",
		NumberOfObjects: noOfObjects,
	}
}

func generateMeterReadings(numberOfReadings int) []gs2.MeterReading {
	var meterReadings []gs2.MeterReading

	for i := 0; i < numberOfReadings; i++ {
		meterReadings = append(meterReadings, generateMeterReading(i))
	}

	return meterReadings
}

func generateMeterReading(n int) gs2.MeterReading {
	itoa := strconv.Itoa(n)
	return gs2.MeterReading{
		Reference:     "meterpoint" + itoa,
		Time:          time.Time{},
		Unit:          "kWh",
		Value:         generateTriplet(),
		MeterLocation: "location" + itoa,
		Meter:         "meter" + itoa,
		Description:   "description for entry " + itoa,
	}
}

func generateTimeSeries(numberofSeries int) []gs2.TimeSeries {
	var timeSeries []gs2.TimeSeries

	for i := 0; i < numberofSeries; i++ {
		timeSeries = append(timeSeries, generateTimeSerie(i))
	}

	return timeSeries
}

func generateTimeSerie(n int) gs2.TimeSeries {
	itoa := strconv.Itoa(n)
	numberOfValues := 24
	triplets := generateTriplets(numberOfValues)

	var sum float64
	for _, triplet := range triplets {
		sum += triplet.Value
	}

	return gs2.TimeSeries{
		Reference:       "meterpoint" + itoa,
		Start:           time.Time{},
		Stop:            time.Time{},
		Step:            time.Hour,
		Unit:            "kWh",
		TypeOfValue:     "interval",
		DirectionOfFlow: "out",
		Value:           triplets,
		NoOfValues:      numberOfValues,
		Sum:             sum,
		MeterLocation:   "location" + itoa,
		Meter:           "meter" + itoa,
		Description:     "description for entry " + itoa,
	}
}

func generateTriplets(numberOfTriplet int) []gs2.Triplet {
	var triplets []gs2.Triplet

	for i := 0; i < numberOfTriplet; i++ {
		triplets = append(triplets, generateTriplet())
	}

	return triplets
}

func generateTriplet() gs2.Triplet {
	return gs2.Triplet{
		Value: math.Trunc(rand.Float64()*maxVal*10000) / 10000,
	}
}
