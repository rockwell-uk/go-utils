package timeutils

import (
	"fmt"
	"time"
)

var divs = []time.Duration{
	time.Duration(1),
	time.Duration(10),
	time.Duration(100),
	time.Duration(1000),
}

func Round(d time.Duration, digits int) time.Duration {

	switch {
	case d > time.Second:
		d = d.Round(time.Second / divs[digits])
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / divs[digits])
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / divs[digits])
	}

	return d
}

func FormatDuration(d time.Duration, digits int) string {

	s := float64(d.Milliseconds()) / 1000

	ph := "%."
	formatStr := fmt.Sprintf("%v%vfs", ph, digits)

	return fmt.Sprintf(formatStr, s)
}

func FormatTime(d time.Time) string {
	return d.Format("2006-01-02 15:04:05")
}

func Took(start time.Time) time.Duration {

	var duration time.Duration = time.Since(start)

	return Round(duration, 2)
}
