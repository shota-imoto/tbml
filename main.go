package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/tbml/tabsvg"
)

func main() {
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

	// タイトル
	// canvas.Text(10, 20, "My Original Music Score", "font-size:14px;fill:black")

	// 楽譜を描画
	drawLine(canvas)
	// c := parse.ParseConfig("tab.yaml")
	// s, _ := c.Build()

	// for _, l := range s.Lines {

	// 	for _, m := range l.Measures {
	// 		m.Draw(canvas)
	// 	}
	// }

}

var SPACE int = 10
var NOTE_WIDTH int = 20
var MEASURE_LINE_DEFINE string = "stroke:#bbb;stroke-width:1"

func drawLine(canvas *svg.SVG) {
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

	f1, _ := m1.AddFingering("2", 3, 1)
	f1.Draw(canvas)
	t := f1.AddLegatoTechnique("3", 1, "h")
	t.Draw(canvas)

	f3, _ := m1.AddFingering("0", 2, 1)
	f3.Draw(canvas)
	f4, _ := m1.AddFingering("0", 5, 1)
	f4.Draw(canvas)
	f5, _ := m1.AddFingering("0", 1, 1)
	f5.Draw(canvas)

	f6, _ := m2.AddFingering("3", 2, 2)
	f6.Draw(canvas)
	t2 := f6.AddLegatoTechnique("8", 2, "s")
	t2.Draw(canvas)
	m2.AddWhiteSpace(2)

	fs, _ := m2.AddMultiFingering(2, tabsvg.FingeringInput{Fret: "0", Strings: 1}, tabsvg.FingeringInput{Fret: "0", Strings: 5})
	for _, f := range fs {
		f.Draw(canvas)
	}
	f7, _ := m2.AddFingering("0", 3, 1)
	f7.Draw(canvas)

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
