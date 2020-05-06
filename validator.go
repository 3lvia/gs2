package gs2

import (
	"fmt"
	"math"
)

// Validator is a function taking in a refenrece to a GS" object and returns an error if its not valid.
type Validator func(*GS2) error

// ValidateNoOfObjects validates that the reported number of objects are equal to the actual number of objects in the decoded object.
func ValidateNoOfObjects(g *GS2) error {
	startNoOfObjects := g.StartMessage.NumberOfObjects
	endNoOfObjects := g.EndMessage.NumberOfObjects

	if (startNoOfObjects != 0 && endNoOfObjects != 0) && (startNoOfObjects != endNoOfObjects) {
		return fmt.Errorf("conflicting number of objects in StartMessage and EndMessage")
	}

	var noOfObjects int
	if startNoOfObjects != 0 {
		noOfObjects = startNoOfObjects
	} else {
		noOfObjects = endNoOfObjects
	}

	actualNoOfObjects := len(g.MeterReadings) + len(g.TimeSeries) + 2

	if actualNoOfObjects != noOfObjects {
		return fmt.Errorf("number of objects not matching. Found %d, but start/end says %d", actualNoOfObjects, noOfObjects)
	}

	return nil
}

const delta = 0.000001

// ValidateTimeSeriesValues validates the number ov values are consistent and that sum og values in a time series block are equal to the
// sum attribute.
func ValidateTimeSeriesValues(g *GS2) error {
	for _, timeSeries := range g.TimeSeries {
		if len(timeSeries.Value) != timeSeries.NoOfValues {
			return fmt.Errorf("the number of values does not equal the No-of-values attribute. Expected %d, but got %d", timeSeries.NoOfValues, len(timeSeries.Value))
		}

		var sum float64
		for _, value := range timeSeries.Value {
			sum += value.Value
		}

		if math.Abs(sum-timeSeries.Sum) > delta {
			return fmt.Errorf("calculated sum is different from sum attribute. Expected: %f, but calculated %f", timeSeries.Sum, sum)
		}
	}

	return nil
}
