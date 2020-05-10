# GS2 decoder/encoder
Tool for decoding and encoding data according to the GS2 format (https://sintef.brage.unit.no/sintef-xmlui/handle/11250/2391930).

Current version is mostly forgiving. Mainly because there is a a lot of variation in the applications of this format in pratice.

Inspired by go json decoding/encoding, but does not support custom types as of yet. 

## Usage
```go
package main

import (
	"github.com/rejlersembriq/gs2"
	"log"
	"os"
)

func main() {
	file, _ := os.Open("someGS2File.gs2")

	// Decode
	g, err := gs2.NewDecoder(file).Decode()
	if err != nil {
		log.Fatalf("error decoding: %v", err)
	}

	// Encode
	if err = gs2.NewEncoder(file).Encode(g); err != nil {
		log.Fatalf("error encoding: %v", err)
	}
}
```
### Encoder/Decoder Options
Current options supported:
- Decoder
    - DecodeValidators (slice of Validator to be run on GS2 object after decoding)
- Encoder
    - EncodeValidators (slice of Validator to be run on GS2 object before encoding)
    - EncodeFloatPrecision (sets float precision when encoding floats. Default -1 = auto)
    
# Validator
A validator is simply a function with a pointer to a GS2 object as an argument and an error as a return value.
```go
type Validator func(*GS2) error
```
If you need some kind of validation that is not already defined you can define your own Validators and add them to the
Encoder/Decoder before encoding/decoding. NB: When adding Validators manually remeber to also add the default validators if they
are needed. Validators can also be disabled by providing an empty slice. 

# Example
```go
package main

import (
	"fmt"
	"github.com/rejlersembriq/gs2"
	"log"
	"os"
)

func customValidator(g *gs2.GS2) error {
	if g.StartMessage.Description == "" {
		return fmt.Errorf("startmessage needs a description")
	}

	return nil
}

func main() {
	file, _ := os.Open("someGS2File.gs2")

	g := gs2.GS2{}

	options := []gs2.EncoderOption{
		gs2.EncodeFloatPrecision(4),
		gs2.EncodeValidators(
			gs2.ValidateNoOfObjects,
			gs2.ValidateTimeSeriesValues,
			customValidator,
		),
	}

	// Encode
	encoder := gs2.NewEncoder(file, options...)
	if err := encoder.Encode(&g); err != nil {
		log.Fatalf("error encoding: %v", err)
	}
}
```
