package gs2

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Encoder encodes a GS2 object and writes to an io.Writer. NB: year, month and day is not supported in Step attribute. Only hour,
//minute and seconds are used when encoding duration.
type Encoder struct {
	options encoderOptions
	w       io.Writer
	buf     []byte
}

type encoderOptions struct {
	floatPrecision int
	validators     []Validator
}

var defaultEncoderOptions = encoderOptions{
	floatPrecision: -1,
	validators: []Validator{
		ValidateNoOfObjects,
		ValidateTimeSeriesValues,
	},
}

// EncoderOption sets configuration for a Encoder.
type EncoderOption func(*encoderOptions)

// EncodeFloatPrecision sets precision when encoding floats.
func EncodeFloatPrecision(i int) EncoderOption {
	return func(o *encoderOptions) {
		o.floatPrecision = i
	}
}

// EncodeValidators sets the validators to be run before encoding an object. Will overwrite the default ones. So remeber to add
// the defaults as well if needed.
func EncodeValidators(v ...Validator) DecoderOption {
	return func(o *decoderOptions) {
		o.validators = v
	}
}

// NewEncoder returna a new Encoder writing to w.
func NewEncoder(w io.Writer, opt ...EncoderOption) *Encoder {
	opts := defaultEncoderOptions

	for _, o := range opt {
		o(&opts)
	}

	return &Encoder{
		options: opts,
		w:       w,
	}
}

// Encode encodes and writes a GS2 object to an io.Writer.
func (e *Encoder) Encode(g *GS2) error {
	for _, validator := range e.options.validators {
		if err := validator(g); err != nil {
			return err
		}
	}

	return e.encode(reflect.ValueOf(g))
}

func (e *Encoder) encode(v reflect.Value) error {
	indirect := reflect.Indirect(v)

	for i := 0; i < indirect.NumField(); i++ {
		blockName, exists := indirect.Type().Field(i).Tag.Lookup("gs2")
		if !exists {
			return fmt.Errorf("type %s does not have a gs2 tag defined", indirect.Type().Field(i).Type)
		}

		field := indirect.Field(i)
		if field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				e.write([]byte("##" + blockName + "\n"))
				if err := e.block(field.Index(j)); err != nil {
					return err
				}
				e.write([]byte("\n"))
			}
		} else {
			e.write([]byte("##" + blockName + "\n"))
			if err := e.block(field); err != nil {
				return err
			}

			if field.Type() != reflect.TypeOf((*EndMessage)(nil)).Elem() {
				e.write([]byte("\n"))
			}
		}
	}

	if _, err := e.w.Write(e.buf); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) block(v reflect.Value) error {
	indirect := reflect.Indirect(v)

	for i := 0; i < indirect.NumField(); i++ {
		field := indirect.Field(i)

		attributeName, exists := indirect.Type().Field(i).Tag.Lookup("gs2")
		if !exists {
			return fmt.Errorf("type %s does not have a gs2 tag defined", indirect.Type().Field(i).Type)
		}

		// TODO: Add omitempty/required tags. This is a temporary hack to print sums that are 0.0.
		if field.IsZero() && attributeName != "Sum" {
			continue
		}

		e.write([]byte("#" + attributeName + "="))
		if err := e.attribute(indirect.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) attribute(v reflect.Value) error {
	indirect := reflect.Indirect(v)

	if indirect.Kind() == reflect.Slice {
		e.write([]byte("< "))
		for j := 0; j < indirect.Len(); j++ {
			if err := e.value(indirect.Index(j)); err != nil {
				return err
			}
			e.write([]byte(" "))
		}
		e.write([]byte(">\n"))
	} else {
		if err := e.value(v); err != nil {
			return err
		}
		e.write([]byte("\n"))
	}

	return nil
}

// TODO: GMT-Reference should be on the form +/-hh
func (e *Encoder) value(v reflect.Value) error {
	indirect := reflect.Indirect(v)

	switch indirect.Kind() {
	case reflect.String:
		e.write([]byte(indirect.String()))
	case reflect.Int:
		e.write([]byte(strconv.FormatInt(indirect.Int(), 10)))
	case reflect.Int64:
		switch indirect.Type() {
		case reflect.TypeOf((*time.Duration)(nil)).Elem():
			e.write([]byte(encodeDuration(indirect.Interface().(time.Duration))))
		default:
			return fmt.Errorf("type %s not supported", indirect.Type())
		}
	case reflect.Float64:
		e.write([]byte(strconv.FormatFloat(indirect.Float(), 'f', e.options.floatPrecision, 64)))
	case reflect.Struct:
		switch indirect.Type() {
		case reflect.TypeOf((*time.Time)(nil)).Elem():
			e.write([]byte(indirect.Interface().(time.Time).Format(gs2TimeLayout)))
		case reflect.TypeOf((*time.Duration)(nil)).Elem():
			e.write([]byte(encodeDuration(indirect.Interface().(time.Duration))))
		case reflect.TypeOf((*Triplet)(nil)).Elem():
			e.write([]byte(e.encodeTriplet(indirect.Interface().(Triplet))))
		default:
			return fmt.Errorf("type %s not supported", indirect.Type())
		}
	default:
		return fmt.Errorf("type %s not supported", indirect.Type())
	}

	return nil
}

func (e *Encoder) encodeTriplet(t Triplet) string {
	value := strconv.FormatFloat(t.Value, 'f', e.options.floatPrecision, 64)
	var timePart string
	if !reflect.ValueOf(t.Time).IsZero() {
		timePart = t.Time.Format(gs2TimeLayout)
	}

	return value + "/" + timePart + "/" + t.Quality
}

func (e *Encoder) write(val []byte) {
	e.buf = append(e.buf, val...)
}

func encodeDuration(d time.Duration) string {
	return strings.Replace(time.Time{}.Add(d).Format(gs2TimeLayout), "0001-01-01", "0000-00-00", 1)
}
