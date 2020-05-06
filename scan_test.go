package gs2

import (
	"testing"
)

func TestScanner(t *testing.T) {
	for _, test := range scanTestTable {
		s := newScanner()

		var steps []int
		for _, b := range []byte(test.data) {
			steps = append(steps, s.step(s, b))
		}

		if len(steps) != len(test.expected) {
			t.Errorf("Number of steps doesn't match. Expected: %d, but got %d\n", len(test.expected), len(steps))
		}

		for i := 0; i < min(len(steps), len(test.expected)); i++ {
			if steps[i] != test.expected[i] {
				t.Errorf("expected %d, but got %d", test.expected[i], steps[i])
			}
		}
	}
}

var scanTestTable = []struct {
	data     string
	expected []int
}{
	{`a`,
		[]int{
			scanError,
		},
	},
	{`##Block1
#Att1=val1
#Att2=val2`,
		[]int{
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
		},
	},
	{`##Block1#Att1=val1#Att2=val2`,
		[]int{
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
		},
	},
	{`##Block1
#Att1=val1
#Att2=val2

##Block2
#Att1=val1
#Att2=val2`,
		[]int{
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanSkipSpace,
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
		},
	},
	{`##Block1#Att1=val1#Att2=val2##Block2#Att1=val1#Att2=val2`,
		[]int{
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
		},
	},
	{`##Block1
#Att1=< val1 val2 val3 >
#Att2=val2`,
		[]int{
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue,
			scanArrayStart, scanArraySeparator, scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator,
			scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator,
			scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator, scanArrayEnd, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
		},
	},
	{`##Block1#Att1=< val1 val2 val3 >#Att2=val2`,
		[]int{
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue,
			scanArrayStart, scanArraySeparator, scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator,
			scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator,
			scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator, scanArrayEnd,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
		},
	},
	{`##Block1
#Att1=<val1 val2 val3>
#Att2=val2`,
		[]int{
			scanHash, scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanContinue, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue,
			scanArrayStart, scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator,
			scanContinue, scanContinue, scanContinue, scanContinue, scanArraySeparator,
			scanContinue, scanContinue, scanContinue, scanContinue, scanArrayEnd, scanSkipSpace,
			scanHash, scanContinue, scanContinue, scanContinue, scanContinue, scanBeginValue, scanContinue, scanContinue, scanContinue, scanContinue,
		},
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
