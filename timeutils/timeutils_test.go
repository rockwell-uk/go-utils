package timeutils

import (
	"fmt"
	"testing"
	"time"
)

func TestRound(t *testing.T) {
	tests := map[string]struct {
		Duration    time.Duration
		Expected    map[int]string
		Unformatted string
	}{
		"1h0m1.123s": {
			time.Hour + time.Second + 123*time.Millisecond,
			map[int]string{
				0: "1h0m1s",
				1: "1h0m1.1s",
				2: "1h0m1.12s",
				3: "1h0m1.123s",
			},
			"1h0m1.123s",
		},
		"123.456789ms": {
			123456789 * time.Nanosecond,
			map[int]string{
				0: "123ms",
				1: "123.5ms",
				2: "123.46ms",
				3: "123.457ms",
			},
			"123.456789ms",
		},
		"123.456µs": {
			123456 * time.Nanosecond,
			map[int]string{
				0: "123µs",
				1: "123.5µs",
				2: "123.46µs",
				3: "123.456µs",
			},
			"123.456µs",
		},
		"123ns": {
			123 * time.Nanosecond,
			map[int]string{
				0: "123ns",
				1: "123ns",
				2: "123ns",
				3: "123ns",
			},
			"123ns",
		},
	}

	for name, test := range tests {
		for digits := 0; digits <= 3; digits++ {
			actual := fmt.Sprintf("%v", Round(test.Duration, digits))

			if test.Expected[digits] != actual {
				t.Errorf("%s [%v]: expected %v, got %v", name, digits, test.Expected[digits], actual)
			}
		}

		unformatted := fmt.Sprintf("%v", test.Duration)

		if test.Unformatted != unformatted {
			t.Errorf("%s: expected %v, got %v", name, test.Unformatted, unformatted)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := map[string]struct {
		Duration    time.Duration
		Digits      int
		Expected    string
		Unformatted string
	}{
		"seconds": {
			11000000011,
			2,
			"11.00s",
			"11.000000011s",
		},
		"mins and secs": {
			71000000011,
			2,
			"71.00s",
			"1m11.000000011s",
		},
		"hours, mins and secs": {
			7010000000111,
			2,
			"7010.00s",
			"1h56m50.000000111s",
		},
	}

	for name, test := range tests {
		actual := FormatDuration(test.Duration, 2)

		if test.Expected != actual {
			t.Errorf("%s: expected %v, got %v", name, test.Expected, actual)
		}

		unformatted := fmt.Sprintf("%v", test.Duration)

		if test.Unformatted != unformatted {
			t.Errorf("%s: expected %v, got %v", name, test.Unformatted, unformatted)
		}
	}
}

func TestFormatTime(t *testing.T) {
	tests := map[string]struct {
		Time        string
		Digits      int
		Expected    string
		Unformatted string
	}{
		"seconds": {
			"2021-12-23T11:45:26.371Z",
			2,
			"2021-12-23 11:45:26",
			"2021-12-23 11:45:26.371 +0000 UTC",
		},
	}

	for name, test := range tests {
		pt, err := time.Parse(time.RFC3339, test.Time)
		if err != nil {
			t.Fatal(err)
		}

		actual := FormatTime(pt)

		if test.Expected != actual {
			t.Errorf("%s: expected %v, got %v", name, test.Expected, actual)
		}

		unformatted := fmt.Sprintf("%v", pt)

		if test.Unformatted != unformatted {
			t.Errorf("%s: expected %v, got %v", name, test.Unformatted, unformatted)
		}
	}
}
