package stringutils

import (
	"testing"
)

func TestUcFirst(t *testing.T) {

	tests := map[string]string{
		"":    "",
		" ":   " ",
		"foo": "Foo",
		"bar": "Bar",
		"found 55 table(s) in the source database [glasshouse]": "Found 55 table(s) in the source database [glasshouse]",
	}

	for subject, expected := range tests {

		actual := UcFirst(subject)

		if expected != actual {
			t.Errorf("Expected [%s]\nGot [%s]", expected, actual)
		}
	}
}

func TestToSnakeCase(t *testing.T) {

	tests := map[string]string{
		"":        "",
		" ":       " ",
		"foo bar": "foo bar",
		"fooBar":  "foo_bar",
	}

	for subject, expected := range tests {

		actual := ToSnakeCase(subject)

		if expected != actual {
			t.Errorf("Expected [%s]\nGot [%s]", expected, actual)
		}
	}
}

func TestSpacePad(t *testing.T) {
	tests := []struct {
		strToPad    string
		length      int
		paddedLeft  string
		paddedRight string
	}{
		{
			strToPad:    "apple",
			length:      8,
			paddedLeft:  "   apple",
			paddedRight: "apple   ",
		},
		{
			strToPad:    "banana",
			length:      5,
			paddedLeft:  "banana",
			paddedRight: "banana",
		},
	}
	for _, tt := range tests {
		pl := SpacePadLeft(tt.strToPad, tt.length)
		if tt.paddedLeft != pl {
			t.Fatalf("expected '%v', actual '%v'", tt.paddedLeft, pl)
		}

		pr := SpacePadRight(tt.strToPad, tt.length)
		if tt.paddedRight != pr {
			t.Fatalf("expected '%v', actual '%v'", tt.paddedRight, pr)
		}
	}
}
