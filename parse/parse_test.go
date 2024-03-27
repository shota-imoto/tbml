package parse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tbml/tabsvg"
)

func TestParseFingering(t *testing.T) {
	tests := []struct {
		str             string
		fingering_input []tabsvg.FingeringInput
		length          int
	}{
		{
			str: "2s4.3 2",
			fingering_input: []tabsvg.FingeringInput{
				{Strings: "3", Fret: "2", Techniques: []tabsvg.AddLegatoTechniqueInput{{Fret: "4", Length: 2, Text: "s"}}},
			},
			length: 2,
		},
		{
			str: "5.1/0.5 2",
			fingering_input: []tabsvg.FingeringInput{
				{Strings: "1", Fret: "5"},
				{Strings: "5", Fret: "0"},
			},
			length: 2,
		},
		{
			str: "0.1/2s4.3 2",
			fingering_input: []tabsvg.FingeringInput{
				{Strings: "1", Fret: "0"},
				{Strings: "3", Fret: "2", Techniques: []tabsvg.AddLegatoTechniqueInput{{Fret: "4", Length: 2, Text: "s"}}},
			},
			length: 2,
		},
		{
			str: "0.3",
			fingering_input: []tabsvg.FingeringInput{
				{Strings: "3", Fret: "0"},
			},
			length: 1,
		},
	}

	for _, test := range tests {
		fs, l, err := ParseFingering(test.str)

		if err != nil {
			t.Errorf("error is raised: %v", err)
		}

		if diff := cmp.Diff(fs, test.fingering_input); diff != "" {
			t.Errorf("fingering_inputs (-fs +fingering_input):%s\n", diff)
		}

		if l != test.length {
			t.Errorf("length is different: expect: %d, want: %d", test.length, l)
		}
	}

}
