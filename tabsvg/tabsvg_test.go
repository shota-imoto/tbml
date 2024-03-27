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

	file, err := os.Create("tab.svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	canvas := svg.New(file)

	canvas.Start(width, height)
	defer canvas.End()

	s := tabsvg.NewScore(tabsvg.Cordinate{X: 200, Y: 200}, 20)
	l1 := s.AddNewLine(5, false)
	m1 := l1.AddNewMeasure(8, "")
	m2 := l1.AddNewMeasure(8, "")
	m3 := l1.AddNewMeasure(8, "")

	m1.Draw(canvas)
	m2.Draw(canvas)
	m3.Draw(canvas)
	tabsvg.MeasureBorder{Measure: *m1}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: *m2}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: *m3}.DrawStart(canvas)

	l2 := s.AddNewLine(6, true)
	m4 := l2.AddNewMeasure(6, "break")
	m5 := l2.AddNewMeasure(6, "柔らかめの音で弾く")
	m6 := l2.AddNewMeasure(6, "")

	fs2, _ := m4.AddFingerings(1, tabsvg.FingeringInput{Fret: "0", Strings: "1", Techniques: []tabsvg.AddLegatoTechniqueInput{}}, tabsvg.FingeringInput{Fret: "2", Strings: "3", Techniques: []tabsvg.AddLegatoTechniqueInput{{Fret: "4", Length: 2, Text: "s"}}})

	for _, f := range fs2 {
		f.Draw(canvas)
		for _, t := range f.Technique {
			t.Draw(canvas)
		}
	}

	fs3, _ := m4.AddFingerings(1, tabsvg.FingeringInput{Fret: "0", Strings: "5", Techniques: []tabsvg.AddLegatoTechniqueInput{}})
	fs3[0].Draw(canvas)
	fs4, _ := m4.AddFingerings(1, tabsvg.FingeringInput{Fret: "0", Strings: "2", Techniques: []tabsvg.AddLegatoTechniqueInput{}})
	fs4[0].Draw(canvas)

	m4.Draw(canvas)
	m5.Draw(canvas)
	m6.Draw(canvas)
	tabsvg.MeasureBorder{Measure: *m4}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: *m5}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: *m6}.DrawStart(canvas)
}
