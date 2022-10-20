package gs2

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const gs2TimeLayout = "2006-01-02.15:04:05"

const scanEnd = -1

// Decoder reads and decodes GS2 input. NB: year, month and day is not supported in Step attribute. Only hour, minute and seconds
// are used when decoding duration.
type Decoder struct {
	options       decoderOptions
	rdr           io.Reader
	scan          *scanner
	buf           []byte
	bytesRead     int
	lastByteRead  byte
	lastScanState int
	typeCache     map[reflect.Type]map[string]int
}

type decoderOptions struct {
	validators []Validator
}

var defaultDecoderOptions = decoderOptions{
	validators: []Validator{
		ValidateNoOfObjects,
		ValidateTimeSeriesValues,
	},
}

// DecoderOption sets configuration for a Decoder.
type DecoderOption func(*decoderOptions)

// DecodeValidators sets the validators to be run after decoding an object. Will overwrite the default ones. So remeber to add the
// defaults as well if needed.
func DecodeValidators(v ...Validator) DecoderOption {
	return func(o *decoderOptions) {
		o.validators = v
	}
}

// NewDecoder returna a new Decoder reading from r.
func NewDecoder(r io.Reader, opt ...DecoderOption) *Decoder {
	opts := defaultDecoderOptions

	for _, o := range opt {
		o(&opts)
	}

	return &Decoder{
		options:   opts,
		rdr:       r,
		scan:      newScanner(),
		typeCache: make(map[reflect.Type]map[string]int),
	}
}

// Decode reads the input and puts it in a GS2 object.
func (d *Decoder) Decode() (*GS2, error) {
	result := &GS2{}

	err := d.decode(reflect.ValueOf(result))
	if err != nil {
		return nil, err
	}

	for _, validator := range d.options.validators {
		if err := validator(result); err != nil {
			return nil, err
		}
	}

	var gmtOffset = gmtReferenceToOffset(result.StartMessage.GMTReference)

	result.StartMessage.Time = addGmtOffset(result.StartMessage.Time, gmtOffset)
	result.EndMessage.Time = addGmtOffset(result.EndMessage.Time, gmtOffset)

	for i := range result.MeterReadings {
		result.MeterReadings[i].Time = addGmtOffset(result.MeterReadings[i].Time, gmtOffset)
	}

	for i := range result.TimeSeries {
		result.TimeSeries[i].Start = addGmtOffset(result.TimeSeries[i].Start, gmtOffset)
		result.TimeSeries[i].Stop = addGmtOffset(result.TimeSeries[i].Stop, gmtOffset)
	}

	return result, nil
}

func addGmtOffset(incomingTime time.Time, gmtOffset time.Duration) time.Time {
	if (incomingTime == time.Time{}) {
		return incomingTime
	}

	return incomingTime.Add(gmtOffset)
}

func gmtReferenceToOffset(gmtReference int) time.Duration {
	return time.Hour * time.Duration(-gmtReference)
}

func (d *Decoder) decode(v reflect.Value) error {
	if err := d.fillBuffer(); err != nil {
		return err
	}

	// Scan to the first # in the file, which should be the first block. The following block will be identified by one # since
	// the first # of the block will be the delimiter of the previous blocks last value. As per the specification spaces are not
	// to be used as delimiters.
	d.scanWhile(scanHash)

	for d.bytesRead < len(d.buf) {
		d.scanNext()
		switch d.lastScanState {
		// Scan for two ## which is the start of a block.
		case scanHash:
			if err := d.block(v); err != nil {
				return err
			}
		case scanSkipSpace:
			continue
		default:
			return fmt.Errorf("unable to find start of block. Got character %q", d.buf[d.bytesRead-1])
		}
	}

	return nil
}

func (d *Decoder) block(v reflect.Value) error {
	dataStart := d.bytesRead

	var blockName []byte
loop:
	for {
		d.scanNext()
		switch d.lastScanState {
		case scanContinue:
			blockName = append(blockName, d.lastByteRead)
		case scanSkipSpace:
		case scanHash:
			break loop
		default:
			return fmt.Errorf("unexpected state while decoding block: %s", string(d.buf[dataStart:d.bytesRead]))
		}
	}

	field, exists := d.getField(string(blockName), reflect.Indirect(v).Type())
	if !exists {
		d.skipBlock()
		return nil
	}

	vf := reflect.Indirect(v).Field(field)

	var block reflect.Value
	if vf.Kind() == reflect.Slice {
		block = reflect.New(vf.Type().Elem())
	} else {
		block = reflect.New(vf.Type())
	}

	for d.lastScanState == scanHash && d.peek(0) != '#' {
		if err := d.attribute(block); err != nil {
			return err
		}
	}

	if vf.Kind() == reflect.Slice {
		vf.Set(reflect.Append(vf, reflect.Indirect(block)))
	} else {
		vf.Set(reflect.Indirect(block))
	}

	return nil
}

func (d *Decoder) attribute(v reflect.Value) error {
	dataStart := d.bytesRead

	var attributeName []byte
loop:
	for {
		d.scanNext()
		switch d.lastScanState {
		case scanContinue:
			attributeName = append(attributeName, d.lastByteRead)
		case scanBeginValue:
			break loop
		default:
			return fmt.Errorf("unexpected state while decoding attribute: %s", string(d.buf[dataStart:d.bytesRead]))
		}
	}

	field, exists := d.getField(string(attributeName), reflect.Indirect(v).Type())
	if !exists {
		d.skipAttribute()
		return nil
	}

	vf := reflect.Indirect(v).Field(field)
	attribute := reflect.New(vf.Type())

	if reflect.Indirect(attribute).Kind() == reflect.Slice {
		if err := d.array(attribute); err != nil {
			return err
		}
	} else {
		if err := d.value(attribute); err != nil {
			return err
		}
	}

	vf.Set(reflect.Indirect(attribute))

	return nil
}

func (d *Decoder) value(v reflect.Value) error {
	dataStart := d.bytesRead

	var value []byte
loop:
	for {
		d.scanNext()

		switch d.lastScanState {
		case scanContinue:
			value = append(value, d.lastByteRead)

		case scanSkipSpace:
		case scanHash:
			break loop
		case scanEnd:
			break loop
		default:
			return fmt.Errorf("unexpected state while decoding value: %s", string(d.buf[dataStart:d.bytesRead]))
		}
	}

	indirect := reflect.Indirect(v)
	switch indirect.Kind() {
	case reflect.String:
		indirect.SetString(string(value))
	case reflect.Int:
		pi, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return err
		}
		indirect.SetInt(pi)
	case reflect.Int64:
		switch indirect.Type() {
		case reflect.TypeOf((*time.Duration)(nil)).Elem():
			d, err := parseDuration(string(value))
			if err != nil {
				return err
			}
			indirect.Set(reflect.ValueOf(d))
		default:
			return fmt.Errorf("unsupported type %q", indirect.Type().Name())
		}
	case reflect.Float64:
		pf, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return err
		}
		indirect.SetFloat(pf)
	case reflect.Struct:
		switch indirect.Type() {
		case reflect.TypeOf((*time.Time)(nil)).Elem():
			t, err := parseTime(string(value))
			if err != nil {
				return err
			}
			indirect.Set(reflect.ValueOf(t))
		case reflect.TypeOf((*Triplet)(nil)).Elem():
			t, err := parseTriplet(string(value))
			if err != nil {
				return err
			}
			indirect.Set(reflect.ValueOf(t))
		default:
			return fmt.Errorf("unsupported type %q", indirect.Type().Name())
		}
	default:
		return fmt.Errorf("unsupported type %q", indirect.Type().Name())
	}

	return nil
}

func (d *Decoder) array(v reflect.Value) error {
	dataStart := d.bytesRead

	// Scan for the start of an array
	d.scanWhile(scanArrayStart)

	var value []byte
loop:
	for {
		d.scanNext()

		switch d.lastScanState {
		case scanContinue:
			value = append(value, d.lastByteRead)
		case scanArrayEnd:
			fallthrough
		case scanArraySeparator:
			if len(value) > 0 {
				indirect := reflect.Indirect(v)
				switch indirect.Type().Elem().Kind() {
				case reflect.Struct:
					switch indirect.Type().Elem() {
					case reflect.TypeOf((*Triplet)(nil)).Elem():
						trip, err := parseTriplet(string(value))
						if err != nil {
							return err
						}
						indirect.Set(reflect.Append(indirect, reflect.ValueOf(trip)))
					default:
						return fmt.Errorf("unsupported type %q for arrays", indirect.Type().Elem())
					}
				default:
					return fmt.Errorf("unsupported type %q for arrays", indirect.Type().Elem())
				}
				value = value[:0]
			}
		case scanSkipSpace:
		case scanHash:
			break loop
		default:
			return fmt.Errorf("unexpected state while decoding value: %s", string(d.buf[dataStart:d.bytesRead]))
		}
	}

	return nil
}

func (d *Decoder) fillBuffer() (err error) {
	d.buf, err = ioutil.ReadAll(d.rdr)
	return
}

func (d *Decoder) scanNext() {
	if d.bytesRead == len(d.buf) {
		d.lastScanState = scanEnd
		return
	}
	d.lastScanState = d.scan.step(d.scan, d.buf[d.bytesRead])
	d.lastByteRead = d.buf[d.bytesRead]
	d.bytesRead++
}

func (d *Decoder) scanWhile(state int) {
	for {
		if d.lastScanState == state || d.lastScanState == scanEnd {
			return
		}
		d.scanNext()
	}
}

func (d *Decoder) peek(n int) byte {
	return d.buf[d.bytesRead+n]
}

func (d *Decoder) skipBlock() {
	for {
		d.scanNext()
		switch d.lastScanState {
		case scanHash:
			if d.peek(0) == '#' {
				return
			}
		case scanEnd:
			return
		}
	}
}

func (d *Decoder) skipAttribute() {
	d.scanWhile(scanHash)
}

func (d *Decoder) getField(key string, typ reflect.Type) (int, bool) {
	cachedTyp, isCached := d.typeCache[typ]
	if isCached {
		index, exists := cachedTyp[key]
		if exists {
			return index, true
		}
	}

	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("gs2")
		if strings.EqualFold(key, strings.Split(tag, ",")[0]) {
			if isCached {
				cachedTyp[key] = i
			} else {
				m := make(map[string]int)
				m[key] = i
				d.typeCache[typ] = m
			}

			return i, true
		}
	}
	return 0, false
}

func parseTriplet(val string) (Triplet, error) {
	split := strings.Split(val, "/")

	var v float64
	var t time.Time
	var q string
	var err error

	if len(split) > 0 && split[0] != "" {
		v, err = strconv.ParseFloat(split[0], 64)
		if err != nil {
			return Triplet{}, err
		}
	}

	if len(split) > 1 && split[1] != "" {
		t, err = parseTime(split[1])
		if err != nil {
			return Triplet{}, err
		}
	}

	if len(split) > 2 && split[2] != "" {
		q = split[2]
	}

	return Triplet{
		Value:   v,
		Time:    t,
		Quality: q,
	}, nil
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	var modifier time.Duration
	if strings.Contains(s, "24:00:00") {
		s = strings.Replace(s, "24:00:00", "00:00:00", 1)
		modifier = 24 * time.Hour
	}

	t, err := time.ParseInLocation(gs2TimeLayout, s, time.UTC)
	return t.Add(modifier), err
}

func parseDuration(s string) (time.Duration, error) {
	split := strings.Split(s, ".")
	tp := strings.Split(split[1], ":")

	var duration time.Duration = 0
	var durationUnits = [...]time.Duration{time.Hour, time.Minute, time.Second}

	for i, timePart := range tp {
		var durationFactor, err = strconv.Atoi(timePart)
		if err != nil {
			return 0, err
		}

		duration += durationUnits[i] * time.Duration(durationFactor)
	}

	return duration, nil
}
