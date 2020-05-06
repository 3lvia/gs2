package gs2

import (
	"fmt"
	"unicode"
)

const (
	scanError = iota

	scanContinue
	scanHash
	scanBeginValue
	scanArrayStart
	scanArraySeparator
	scanArrayEnd
	scanSkipSpace
)

type scanner struct {
	step func(*scanner, byte) int
	err  error
}

func newScanner() *scanner {
	return &scanner{step: stateBeginScan}
}

func stateBeginScan(s *scanner, b byte) int {
	if unicode.IsSpace(rune(b)) {
		return scanSkipSpace
	}

	switch b {
	case '#':
		s.step = stateHash
		return scanHash
	}

	return s.error("got invalid character %q looking for start", b)
}

func stateHash(s *scanner, b byte) int {
	switch b {
	case '#':
		s.step = stateBlock
		return scanHash
	case '=':
		s.step = stateValue
		return scanBeginValue
	}

	return scanContinue
}

func stateBlock(s *scanner, b byte) int {
	if unicode.IsSpace(rune(b)) {
		return scanSkipSpace
	}

	if b == '#' {
		s.step = stateHash
		return scanHash
	}

	return scanContinue
}

func stateValue(s *scanner, b byte) int {
	if b != ' ' && unicode.IsSpace(rune(b)) {
		return scanSkipSpace
	}

	switch b {
	case '#':
		s.step = stateHash
		return scanHash
	case '<':
		s.step = stateArray
		return scanArrayStart
	}

	return scanContinue
}

func stateArray(s *scanner, b byte) int {
	switch b {
	case ' ':
		return scanArraySeparator
	case '>':
		s.step = stateBeginScan
		return scanArrayEnd
	}

	return scanContinue
}

func stateError(s *scanner, b byte) int {
	return scanError
}

func (s *scanner) error(format string, args ...interface{}) int {
	s.step = stateError
	s.err = fmt.Errorf(format, args...)
	return scanError
}

func (s *scanner) reset() {
	s.step = stateBeginScan
	s.err = nil
}
