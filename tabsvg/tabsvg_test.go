package tabsvg_test

import (
	"os"
	"testing"

	svg "github.com/ajstarks/svgo"
	"github.com/tbml/tabsvg"
)

func TestTabsvg(t *testing.T) {
	width := 1000
	height := 1000

	file, err := os.Create("tab_test.svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	c := svg.New(file)

	c.Start(width, height)
	defer c.End()

	s := tabsvg.NewScore(tabsvg.Cordinate{X: 100, Y: 200}, 20)
	l1 := s.AddNewLine(5, false)
	m1 := l1.AddNewMeasure(8, "")
	m2 := l1.AddNewMeasure(8, "")
	m3 := l1.AddNewMeasure(8, "")

	b, err := m1.AddBorder(tabsvg.StartPosition{})

	if err != nil {
		t.Fatal(err)
	}

	b.Draw(c)
	b2, _ := m2.AddBorder(tabsvg.StartPosition{})
	b2.Draw(c)
	b3, _ := m3.AddBorder(tabsvg.StartPosition{})
	b3.Draw(c)
	b4, _ := m3.AddBorder(tabsvg.EndPosition{})
	b4.Draw(c)

	m1.Draw(c)
	m2.Draw(c)
	m3.Draw(c)

	l2 := s.AddNewLine(6, true)
	m4 := l2.AddNewMeasure(6, "break")
	m5 := l2.AddNewMeasure(6, "柔らかめの音で弾く")
	m6 := l2.AddNewMeasure(6, "")

	fs2, _ := m4.AddFingerings(1, tabsvg.FingeringInput{Fret: "0", Strings: "1", Techniques: []tabsvg.AddLegatoTechniqueInput{}}, tabsvg.FingeringInput{Fret: "2", Strings: "3", Techniques: []tabsvg.AddLegatoTechniqueInput{{Fret: "4", Length: 2, Text: "s"}}})

	for _, f := range fs2 {
		f.Draw(c)
		for _, t := range f.Technique {
			t.Draw(c)
		}
	}

	fs3, _ := m4.AddFingerings(1, tabsvg.FingeringInput{Fret: "0", Strings: "5", Techniques: []tabsvg.AddLegatoTechniqueInput{}})
	fs3[0].Draw(c)
	fs4, _ := m4.AddFingerings(1, tabsvg.FingeringInput{Fret: "0", Strings: "2", Techniques: []tabsvg.AddLegatoTechniqueInput{}})
	fs4[0].Draw(c)

	m4.Draw(c)
	m5.Draw(c)
	m6.Draw(c)
	b.Draw(c)

	b5, err := m4.AddBorder(tabsvg.StartPosition{})

	if err != nil {
		t.Fatal(err)
	}

	b5.Draw(c)
	b6, _ := m5.AddBorder(tabsvg.StartPosition{})
	b6.Draw(c)
	b7, _ := m6.AddBorder(tabsvg.StartPosition{})
	b7.Draw(c)
	b8, _ := m6.AddBorder(tabsvg.EndPosition{})
	b8.Draw(c)
}
