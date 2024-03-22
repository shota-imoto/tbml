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
	tabsvg.MeasureBorder{Measure: m1}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: m2}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: m3}.DrawStart(canvas)

	f1, _ := m1.AddFingering(2, 3, 1)
	f1.Draw(canvas)
	f2, _ := m1.AddFingering(3, 3, 2)
	f2.Draw(canvas)

	tabsvg.Technique{Start: f1, End: f2, Text: "h"}.Draw(canvas)

	f3, _ := m1.AddFingering(0, 2, 2)
	f3.Draw(canvas)
	f4, _ := m1.AddFingering(0, 5, 3)
	f4.Draw(canvas)
	f5, _ := m1.AddFingering(0, 1, 4)
	f5.Draw(canvas)

	f6, _ := m2.AddFingering(3, 2, 1)
	f6.Draw(canvas)
	// f6.DrawCenter(canvas)
	f7, _ := m2.AddFingering(8, 2, 3)
	f7.Draw(canvas)
	tabsvg.Technique{Start: f6, End: f7, Text: "s"}.Draw(canvas)

	l2 := s.AddNewLine(6, true)
	m4 := l2.AddNewMeasure(6, "break")
	m5 := l2.AddNewMeasure(6, "柔らかめの音で弾く")
	m6 := l2.AddNewMeasure(6, "")

	m4.Draw(canvas)
	m5.Draw(canvas)
	m6.Draw(canvas)
	tabsvg.MeasureBorder{Measure: m4}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: m5}.DrawStart(canvas)
	tabsvg.MeasureBorder{Measure: m6}.DrawStart(canvas)

}