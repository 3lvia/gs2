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
	ID              string    `gs2:"Id,omitempty"`
	MessageType     string    `gs2:"Message-type,omitempty"`
	Version         string    `gs2:"Version,omitempty"`
	Time            time.Time `gs2:"Time,omitempty"`
	To              string    `gs2:"To,omitempty"`
	From            string    `gs2:"From,omitempty"`
	ReferenceTable  string    `gs2:"Reference-table,omitempty"`
	GMTReference    int       `gs2:"GMT-reference,omitempty"`
	NumberOfObjects int       `gs2:"Number-of-objects,omitempty"`
	TypeOfObjects   string    `gs2:"Type-of-objects,omitempty"`
	ContainsObjects string    `gs2:"Contains-objects,omitempty"`
	RequestedAction string    `gs2:"Requested-action,omitempty"`
	Description     string    `gs2:"Description,omitempty"`
}

// EndMessage should always be the last object in any GS2-file-
type EndMessage struct {
	ID              string    `gs2:"Id,omitempty"`
	MessageType     string    `gs2:"Message-type,omitempty"`
	Version         string    `gs2:"Version,omitempty"`
	Time            time.Time `gs2:"Time,omitempty"`
	To              string    `gs2:"To,omitempty"`
	From            string    `gs2:"From,omitempty"`
	ReferenceTable  string    `gs2:"Reference-table,omitempty"`
	GMTReference    int       `gs2:"GMT-reference,omitempty"`
	NumberOfObjects int       `gs2:"Number-of-objects"`
	TypeOfObjects   string    `gs2:"Type-of-objects,omitempty"`
	ContainsObjects string    `gs2:"Contains-objects,omitempty"`
	RequestedAction string    `gs2:"Requested-action,omitempty"`
	Description     string    `gs2:"Description,omitempty"`
}

// MeterReading contains a single value that is a channel reading at a given point in time.
type MeterReading struct {
	Reference     string    `gs2:"Reference,omitempty"`
	Time          time.Time `gs2:"Time,omitempty"`
	Unit          string    `gs2:"Unit,omitempty"`
	Value         Triplet   `gs2:"Value"`
	Installation  string    `gs2:"Installation,omitempty"`
	Plant         string    `gs2:"Plant,omitempty"`
	MeterLocation string    `gs2:"Meter-location,omitempty"`
	NetOwner      string    `gs2:"Net-owner,omitempty"`
	Supplier      string    `gs2:"Supplier,omitempty"`
	Customer      string    `gs2:"Customer,omitempty"`
	Meter         string    `gs2:"Meter,omitempty"`
	Channel       string    `gs2:"Channel,omitempty"`
	Description   string    `gs2:"Description,omitempty"`
}

// TimeSeries contains time series of metered values within the interval given by start and stop.
type TimeSeries struct {
	Reference       string        `gs2:"Reference,omitempty"`
	Start           time.Time     `gs2:"Start,omitempty"`
	Stop            time.Time     `gs2:"Stop,omitempty"`
	Step            time.Duration `gs2:"Step,omitempty"`
	Unit            string        `gs2:"Unit,omitempty"`
	TypeOfValue     string        `gs2:"Type-of-value,omitempty"`
	DirectionOfFlow string        `gs2:"Direction-of-flow,omitempty"`
	Value           []Triplet     `gs2:"Value,omitempty"`
	NoOfValues      int           `gs2:"No-of-values"`
	Sum             float64       `gs2:"Sum"`
	Installation    string        `gs2:"Installation,omitempty"`
	Plant           string        `gs2:"Plant,omitempty"`
	MeterLocation   string        `gs2:"Meter-location,omitempty"`
	NetOwner        string        `gs2:"Net-owner,omitempty"`
	Supplier        string        `gs2:"Supplier,omitempty"`
	Customer        string        `gs2:"Customer,omitempty"`
	Meter           string        `gs2:"Meter,omitempty"`
	Channel         string        `gs2:"Channel,omitempty"`
	Description     string        `gs2:"Description,omitempty"`
}

// Triplet represents a value triplet with value, time and quality.
type Triplet struct {
	Value   float64
	Time    time.Time
	Quality string
}
