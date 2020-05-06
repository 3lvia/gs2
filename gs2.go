package gs2

import "time"

// GS2 - versjon 1.2
//http://hdl.handle.net/11250/2391930

// GS2 represents the data of a GS2 file.
type GS2 struct {
	StartMessage  StartMessage   `gs2:"Start-message"`
	MeterReadings []MeterReading `gs2:"Meter-reading"`
	TimeSeries    []TimeSeries   `gs2:"Time-series"`
	EndMessage    EndMessage     `gs2:"End-message"`
}

// StartMessage should always be the first object in any GS2-file-
type StartMessage struct {
	ID              string    `gs2:"Id"`
	MessageType     string    `gs2:"Message-type"`
	Version         string    `gs2:"Version"`
	Time            time.Time `gs2:"Time"`
	To              string    `gs2:"To"`
	From            string    `gs2:"From"`
	ReferenceTable  string    `gs2:"Reference-table"`
	GMTReference    int       `gs2:"GMT-reference"`
	NumberOfObjects int       `gs2:"Number-of-objects"`
	TypeOfObjects   string    `gs2:"Type-of-objects"`
	ContainsObjects string    `gs2:"Contains-objects"`
	RequestedAction string    `gs2:"Requested-action"`
	Description     string    `gs2:"Description"`
}

// EndMessage should always be the last object in any GS2-file-
type EndMessage struct {
	ID              string    `gs2:"Id"`
	MessageType     string    `gs2:"Message-type"`
	Version         string    `gs2:"Version"`
	Time            time.Time `gs2:"Time"`
	To              string    `gs2:"To"`
	From            string    `gs2:"From"`
	ReferenceTable  string    `gs2:"Reference-table"`
	GMTReference    int       `gs2:"GMT-reference"`
	NumberOfObjects int       `gs2:"Number-of-objects"`
	TypeOfObjects   string    `gs2:"Type-of-objects"`
	ContainsObjects string    `gs2:"Contains-objects"`
	RequestedAction string    `gs2:"Requested-action"`
	Description     string    `gs2:"Description"`
}

// MeterReading contains a single value that is a channel reading at a given point in time.
type MeterReading struct {
	Reference     string    `gs2:"Reference"`
	Time          time.Time `gs2:"Time"`
	Unit          string    `gs2:"Unit"`
	Value         Triplet   `gs2:"Value"`
	Installation  string    `gs2:"Installation"`
	Plant         string    `gs2:"Plant"`
	MeterLocation string    `gs2:"Meter-location"`
	NetOwner      string    `gs2:"Net-owner"`
	Supplier      string    `gs2:"Supplier"`
	Customer      string    `gs2:"Customer"`
	Meter         string    `gs2:"Meter"`
	Channel       string    `gs2:"Channel"`
	Description   string    `gs2:"Description"`
}

// TimeSeries contains time series of metered values within the interval given by start and stop.
type TimeSeries struct {
	Reference       string        `gs2:"Reference"`
	Start           time.Time     `gs2:"Start"`
	Stop            time.Time     `gs2:"Stop"`
	Step            time.Duration `gs2:"Step"`
	Unit            string        `gs2:"Unit"`
	TypeOfValue     string        `gs2:"Type-of-value"`
	DirectionOfFlow string        `gs2:"Direction-of-flow"`
	Value           []Triplet     `gs2:"Value"`
	NoOfValues      int           `gs2:"No-of-values"`
	Sum             float64       `gs2:"Sum"`
	Installation    string        `gs2:"Installation"`
	Plant           string        `gs2:"Plant"`
	MeterLocation   string        `gs2:"Meter-location"`
	NetOwner        string        `gs2:"Net-owner"`
	Supplier        string        `gs2:"Supplier"`
	Customer        string        `gs2:"Customer"`
	Meter           string        `gs2:"Meter"`
	Channel         string        `gs2:"Channel"`
	Description     string        `gs2:"Description"`
}

// Triplet represents a value triplet with value, time and quality.
type Triplet struct {
	Value   float64
	Time    time.Time
	Quality string
}
