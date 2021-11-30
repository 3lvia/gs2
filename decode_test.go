package gs2

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestDecoder_Decode(t *testing.T) {
	for i, test := range decodeTestTable {
		file, err := os.Open(test.inputFile)
		if err != nil {
			t.Fatal(err)
		}

		result, err := NewDecoder(file).Decode()
		if err != nil {
			t.Fatalf("unexpected error when decoding: %v", err)
		}

		if !reflect.DeepEqual(*result, test.expected) {
			t.Errorf("result does not equal expected in test %d from table", i)
		}
	}
}

func TestDecoder_ReadFile(t *testing.T) {
	file, err := os.Open("testdata/timeseries_noNewlines.gs2")
	if err != nil {
		t.Fatal(err)
	}

	result, err := NewDecoder(file).Decode()
	if err != nil {
		t.Fatalf("unexpected error when decoding: %v", err)
	}

	fmt.Printf("%v\n", result)
}

// This is probably stupid since we're also reading the file and stuff. Maybe there is some other way to benchmark?
func benchmarkDecoderDecode(fileName string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, err := os.Open(fileName)
		if err != nil {
			b.Fatal(err)
		}

		if _, err := NewDecoder(file, DecodeValidators()).Decode(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecoder_SmallMeterReadingFile(b *testing.B) {
	benchmarkDecoderDecode("testdata/meterreading.gs2", b)
}

func BenchmarkDecoder_LargeMeterReadingFile(b *testing.B) {
	benchmarkDecoderDecode("testdata/meterreading_large.gs2", b)
}

func BenchmarkDecoder_SmallTimeSeriesFile(b *testing.B) {
	benchmarkDecoderDecode("testdata/timeseries.gs2", b)
}

func BenchmarkDecoder_LargeTimeSeriesFile(b *testing.B) {
	benchmarkDecoderDecode("testdata/timeseries_large.gs2", b)
}

var decodeTestTable = []struct {
	inputFile string
	expected  GS2
}{
	{
		"testdata/meterreading.gs2",
		GS2{
			StartMessage: StartMessage{
				ID:           "0",
				MessageType:  "Settlement-data",
				Version:      "1.2",
				Time:         getTime("2019-07-22T05:37:40Z"),
				To:           "MDM",
				From:         "Sender",
				GMTReference: 1,
			},
			MeterReadings: []MeterReading{
				{
					Reference:       "meterpoint1",
					Meter:           "meter1",
					Time:            getTime("2019-07-19T00:00:00Z"),
					Unit:            "kWh",
					DirectionOfFlow: "out",
					Value: Triplet{
						Value:   84.831,
						Quality: "0",
					},
					Description: "",
				},
				{
					Reference: "meterpoint2",
					Meter:     "meter2",
					Time:      getTime("2019-07-19T01:00:00Z"),
					Unit:      "kWh",
					Value: Triplet{
						Value:   85.078,
						Quality: "",
					},
					Description: "someDescription",
				},
				{
					Reference: "meterpoint3",
					Meter:     "meter3",
					Time:      getTime("2019-07-19T02:00:00Z"),
					Unit:      "kVArh",
					Value: Triplet{
						Value:   85.325,
						Quality: "x",
					},
					Description: "",
				},
			},
			TimeSeries: nil,
			EndMessage: EndMessage{
				ID:              "0",
				NumberOfObjects: 5,
			},
		},
	},
	{
		"testdata/timeseries.gs2",
		GS2{
			StartMessage: StartMessage{
				ID:           "0",
				MessageType:  "Settlement-data",
				Version:      "1.2",
				Time:         getTime("2019-09-25T06:00:08Z"),
				To:           "MDM",
				From:         "Sender",
				GMTReference: 1,
			},
			MeterReadings: nil,
			TimeSeries: []TimeSeries{
				{
					Reference:       "meterpoint1",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "in",
					Value: []Triplet{
						{Value: 0, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 2, Quality: "0"},
					},
					NoOfValues: 24,
					Sum:        27,
				},
				{
					Reference:       "meterpoint2",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "in",
					Value: []Triplet{
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
					},
					NoOfValues: 24,
					Sum:        12,
				},
				{
					Unit:            "kWh",
					DirectionOfFlow: "out",
					Meter:           "meter3",
					Start:           getTime("2020-03-26T22:00:00Z"),
					Stop:            getTime("2020-03-27T22:00:00Z"),
					Step:            time.Hour,
					TypeOfValue:     "interval",
					Value: []Triplet{
						{Value: .02},
						{Value: .02},
						{Value: .07},
						{Value: .13},
						{Value: .12},
						{Value: .11},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .01},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .01},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .01},
						{Value: .02},
						{Value: .02},
					},
					NoOfValues:  24,
					Sum:         0.8,
					Description: "somedescription",
				},
				{
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "out",
					Reference:       "meterpoint4",
					Installation:    "meterpoint4",
					Plant:           "0",
					MeterLocation:   "meterpoint4",
					Meter:           "meter4",
					Channel:         "1",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Value: []Triplet{
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
					},
					NoOfValues: 24,
					Sum:        0.0,
				},
				{
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "in",
					Reference:       "meterpoint5",
					Installation:    "0",
					Plant:           "0",
					MeterLocation:   "location5",
					Meter:           "meter5",
					Channel:         "1",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Value: []Triplet{
						{Value: 70.1},
						{Value: 72},
						{Value: 55.7},
						{Value: 57.6},
						{Value: 56.2},
						{Value: 59.5},
						{Value: 68.9},
						{Value: 79.8},
						{Value: 101.6},
						{Value: 84.9},
						{Value: 84},
						{Value: 85.7},
						{Value: 85.9},
						{Value: 83.9},
						{Value: 84.1},
						{Value: 68.9},
						{Value: 60.5},
						{Value: 58},
						{Value: 56.9},
						{Value: 59.8},
						{Value: 60.1},
						{Value: 62},
						{Value: 57.1},
						{Value: 66.8},
					},
					NoOfValues: 24,
					Sum:        1680.0,
				},
			},
			EndMessage: EndMessage{
				ID:              "0",
				NumberOfObjects: 7,
			},
		},
	},
	{
		"testdata/timeseries_noNewlines.gs2",
		GS2{
			StartMessage: StartMessage{
				ID:           "0",
				MessageType:  "Settlement-data",
				Version:      "1.2",
				Time:         getTime("2019-09-25T06:00:08Z"),
				To:           "MDM",
				From:         "Sender",
				GMTReference: 1,
			},
			MeterReadings: nil,
			TimeSeries: []TimeSeries{
				{
					Reference:       "meterpoint1",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "in",
					Value: []Triplet{
						{Value: 0, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 2, Quality: "0"},
					},
					NoOfValues: 24,
					Sum:        27,
				},
				{
					Reference:       "meterpoint2",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "in",
					Value: []Triplet{
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 2, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 1, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 0, Quality: "0"},
						{Value: 1, Quality: "0"},
					},
					NoOfValues: 24,
					Sum:        12,
				},
				{
					Unit:            "kWh",
					DirectionOfFlow: "out",
					Meter:           "meter3",
					Start:           getTime("2020-03-26T22:00:00Z"),
					Stop:            getTime("2020-03-27T22:00:00Z"),
					Step:            time.Hour,
					TypeOfValue:     "interval",
					Value: []Triplet{
						{Value: .02},
						{Value: .02},
						{Value: .07},
						{Value: .13},
						{Value: .12},
						{Value: .11},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .01},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .01},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .02},
						{Value: .01},
						{Value: .02},
						{Value: .02},
					},
					NoOfValues:  24,
					Sum:         0.8,
					Description: "somedescription",
				},
				{
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "out",
					Reference:       "meterpoint4",
					Installation:    "meterpoint4",
					Plant:           "0",
					MeterLocation:   "meterpoint4",
					Meter:           "meter4",
					Channel:         "1",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Value: []Triplet{
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
						{Value: 0, Quality: "x"},
					},
					NoOfValues: 24,
					Sum:        0.0,
				},
				{
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "in",
					Reference:       "meterpoint5",
					Installation:    "0",
					Plant:           "0",
					MeterLocation:   "location5",
					Meter:           "meter5",
					Channel:         "1",
					Start:           getTime("2020-03-26T23:00:00Z"),
					Stop:            getTime("2020-03-27T23:00:00Z"),
					Step:            time.Hour,
					Value: []Triplet{
						{Value: 70.1},
						{Value: 72},
						{Value: 55.7},
						{Value: 57.6},
						{Value: 56.2},
						{Value: 59.5},
						{Value: 68.9},
						{Value: 79.8},
						{Value: 101.6},
						{Value: 84.9},
						{Value: 84},
						{Value: 85.7},
						{Value: 85.9},
						{Value: 83.9},
						{Value: 84.1},
						{Value: 68.9},
						{Value: 60.5},
						{Value: 58},
						{Value: 56.9},
						{Value: 59.8},
						{Value: 60.1},
						{Value: 62},
						{Value: 57.1},
						{Value: 66.8},
					},
					NoOfValues: 24,
					Sum:        1680.0,
				},
			},
			EndMessage: EndMessage{
				ID:              "0",
				NumberOfObjects: 7,
			},
		},
	},
	{
		"testdata/durations.gs2",
		GS2{
			StartMessage: StartMessage{
				ID:           "0",
				MessageType:  "Settlement-data",
				Version:      "1.2",
				Time:         getTime("2020-03-27T01:00:00Z"),
				To:           "MDM",
				From:         "Sender",
				GMTReference: 2,
			},
			MeterReadings: nil,
			TimeSeries: []TimeSeries{
				{
					Reference: "meterpoint1",
					Start:     getTime("2020-03-26T22:00:00Z"),
					Stop:      getTime("2020-03-27T08:00:00Z"),
					Step:      time.Hour,
					Value: []Triplet{
						{Value: 1},
						{Value: 2},
						{Value: 3},
						{Value: 4},
						{Value: 5},
						{Value: 6},
						{Value: 7},
						{Value: 8},
						{Value: 9},
						{Value: 10},
					},
					NoOfValues: 10,
					Sum:        55,
				},
				{
					Reference: "meterpoint2",
					Start:     getTime("2020-03-26T22:00:00Z"),
					Stop:      getTime("2020-03-27T08:00:00Z"),
					Step:      2 * time.Hour,
					Value: []Triplet{
						{Value: 0.001},
						{Value: 0.002},
						{Value: 0.003},
						{Value: 0.004},
						{Value: 0.005},
					},
					NoOfValues: 5,
					Sum:        0.015,
				},
				{
					Reference: "meterpoint3",
					Start:     getTime("2020-03-26T22:00:00Z"),
					Stop:      getTime("2020-03-26T23:00:00Z"),
					Step:      15 * time.Minute,
					Value: []Triplet{
						{Value: 10001.01},
						{Value: 10002.02},
						{Value: 1003.03},
						{Value: 1004.04},
					},
					NoOfValues: 4,
					Sum:        22010.1,
				},
				{
					Reference: "meterpoint4",
					Start:     getTime("2020-03-26T22:00:00Z"),
					Stop:      getTime("2020-03-26T22:00:05Z"),
					Step:      time.Second,
					Value: []Triplet{
						{Value: 1.1},
						{Value: 2.2},
						{Value: 3.3},
						{Value: 4.4},
						{Value: 5.5},
					},
					NoOfValues: 5,
					Sum:        16.5,
				},
			},
			EndMessage: EndMessage{
				ID:              "0",
				NumberOfObjects: 6,
			},
		},
	},
}

// Helper function for putting times in tests
func getTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	return t
}
