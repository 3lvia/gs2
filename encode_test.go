package gs2

import (
	"bytes"
	"testing"
	"time"
)

func TestEncoder_Encode(t *testing.T) {
	for _, test := range encodeTestTable {
		var buf bytes.Buffer
		if err := NewEncoder(&buf).Encode(&test.g); err != nil {
			t.Fatalf("unexpected error when encoding: %v", err)
		}

		if buf.String() != test.expected {
			t.Errorf("Expected:\n%s got:\n%s", test.expected, buf.String())
		}
	}
}

var encodeTestTable = []struct {
	g        GS2
	expected string
}{
	{
		GS2{
			StartMessage: StartMessage{
				ID:          "0",
				MessageType: "Settlement-data",
				Version:     "1.2",
				Time:        getTime("2020-04-04T20:00:00Z"),
			},
			MeterReadings: []MeterReading{
				{
					Reference: "meterpoint1",
					Meter:     "meter1",
					Time:      getTime("2020-04-04T00:00:00Z"),
					Unit:      "kWh",
					Value: Triplet{
						Value:   1.1,
						Quality: "",
					},
					Description: "someDescription",
				},
				{
					Reference: "meterpoint2",
					Meter:     "meter2",
					Time:      getTime("2020-04-03T00:00:00Z"),
					Unit:      "kWh",
					Value: Triplet{
						Value:   0,
						Quality: "x",
					},
				},
				{
					Reference: "meterpoint3",
					Meter:     "meter3",
					Time:      getTime("2020-04-02T00:00:00Z"),
					Unit:      "kWh",
					Value: Triplet{
						Value:   2.2,
						Quality: "",
					},
				},
			},
			EndMessage: EndMessage{
				ID:              "0",
				NumberOfObjects: 5,
			},
		},
		`##Start-message
#Id=0
#Message-type=Settlement-data
#Version=1.2
#Time=2020-04-04.20:00:00

##Meter-reading
#Reference=meterpoint1
#Time=2020-04-04.00:00:00
#Unit=kWh
#Value=1.1//
#Meter=meter1
#Description=someDescription

##Meter-reading
#Reference=meterpoint2
#Time=2020-04-03.00:00:00
#Unit=kWh
#Value=0//x
#Meter=meter2

##Meter-reading
#Reference=meterpoint3
#Time=2020-04-02.00:00:00
#Unit=kWh
#Value=2.2//
#Meter=meter3

##End-message
#Id=0
#Number-of-objects=5
`,
	},
	{
		GS2{
			StartMessage: StartMessage{
				ID:          "0",
				MessageType: "Settlement-data",
				Version:     "1.2",
				Time:        getTime("2020-04-04T20:00:00Z"),
			},
			TimeSeries: []TimeSeries{
				{
					Reference:       "meterpoint1",
					Start:           getTime("2020-04-03T00:00:00Z"),
					Stop:            getTime("2020-04-04T00:00:00Z"),
					Step:            time.Hour,
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "out",
					Value: []Triplet{
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
						{Value: 1},
					},
					NoOfValues:  24,
					Sum:         24.0,
					Description: "hourly values",
				},
				{
					Reference:       "meterpoint2",
					Start:           getTime("2020-04-03T00:00:00Z"),
					Stop:            getTime("2020-04-04T00:00:00Z"),
					Step:            time.Hour,
					Unit:            "kWh",
					TypeOfValue:     "interval",
					DirectionOfFlow: "in",
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
					NoOfValues:  24,
					Sum:         0.0,
					Description: "hourly values",
				},
			},
			EndMessage: EndMessage{
				ID:              "0",
				NumberOfObjects: 4,
			},
		},
		`##Start-message
#Id=0
#Message-type=Settlement-data
#Version=1.2
#Time=2020-04-04.20:00:00

##Time-series
#Reference=meterpoint1
#Start=2020-04-03.00:00:00
#Stop=2020-04-04.00:00:00
#Step=0000-00-00.01:00:00
#Unit=kWh
#Type-of-value=interval
#Direction-of-flow=out
#Value=< 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// 1// >
#No-of-values=24
#Sum=24
#Description=hourly values

##Time-series
#Reference=meterpoint2
#Start=2020-04-03.00:00:00
#Stop=2020-04-04.00:00:00
#Step=0000-00-00.01:00:00
#Unit=kWh
#Type-of-value=interval
#Direction-of-flow=in
#Value=< 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x 0//x >
#No-of-values=24
#Sum=0
#Description=hourly values

##End-message
#Id=0
#Number-of-objects=4
`,
	},
}
