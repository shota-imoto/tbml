package main

import (
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
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

}

var SPACE int = 10
var NOTE_WIDTH int = 14
var MEASURE_LINE_DEFINE string = "stroke:black;stroke-width:1"

func drawLine(canvas *svg.SVG) {
	s := NewScore(Cordinate{200, 200}, 20)
	l1 := s.addNewLine(5)
	m1 := l1.addNewMeasure(8)
	m2 := l1.addNewMeasure(8)
	m3 := l1.addNewMeasure(8)
	m1.Draw(canvas)
	m2.Draw(canvas)
	m3.Draw(canvas)
	MeasureBorder{Measure: m1}.drawStart(canvas)
	MeasureBorder{Measure: m2}.drawStart(canvas)
	MeasureBorder{Measure: m3}.drawStart(canvas)

	f := Fingering{Measure: m1, Fret: 12, Strings: 3, Beat: 1}
	f.Draw(canvas)
	Fingering{Measure: m1, Fret: 4, Strings: 3, Beat: 2}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 3, Beat: 3}.Draw(canvas)
	Fingering{Measure: m1, Fret: 4, Strings: 3, Beat: 4}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 3, Beat: 5}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 3, Beat: 6}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 3, Beat: 7}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 3, Beat: 8}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 1, Beat: 8}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 2, Beat: 8}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 3, Beat: 8}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 4, Beat: 8}.Draw(canvas)
	Fingering{Measure: m1, Fret: 12, Strings: 5, Beat: 8}.Draw(canvas)

	l2 := s.addNewLine(6)
	m4 := l2.addNewMeasure(6)
	m5 := l2.addNewMeasure(6)
	m6 := l2.addNewMeasure(6)

	m4.Draw(canvas)
	m5.Draw(canvas)
	m6.Draw(canvas)
	MeasureBorder{Measure: m4}.drawStart(canvas)
	MeasureBorder{Measure: m5}.drawStart(canvas)
	MeasureBorder{Measure: m6}.drawStart(canvas)
	Fingering{Measure: m4, Fret: 4, Strings: 3, Beat: 2}.Draw(canvas)

	// drawMeasure(canvas, Cordinate{100, 100}, 5, 8)
}

type Cordinate struct {
	X int
	Y int
}

// 縦方向にLineを並べたものがScore
type Score struct {
	Base  Cordinate
	EndY  int
	Lines []Line
	Gap   int
}

func NewScore(b Cordinate, g int) Score {
	return Score{Base: b, EndY: b.Y, Lines: []Line{}, Gap: g}
}

func (s *Score) addNewLine(strings int) Line {
	l := NewLine(Cordinate{s.Base.X, s.EndY}, strings)
	s.Lines = append(s.Lines, l)
	s.EndY = s.EndY + l.Height + s.Gap

	return l
}

// 横方向にMeasureを並べたもの
type Line struct {
	Base     Cordinate
	EndX     int
	Measures []Measure
	Height   int
	Strings  int
}

func NewLine(b Cordinate, s int) Line {
	h := (s - 1) * SPACE
	return Line{Base: b, EndX: b.X, Measures: []Measure{}, Strings: s, Height: h}
}

func (l *Line) addNewMeasure(beat int) Measure {
	m := Measure{Base: Cordinate{l.EndX, l.Base.Y}, Strings: l.Strings, Beat: beat}
	l.Measures = append(l.Measures, m)
	l.EndX = l.EndX + m.Width()

	return m
}

// 小節
type Measure struct {
	Base    Cordinate // 小節の左上を0点とする
	Strings int       // 弦の数
	Beat    int       // 拍数
}

func (m Measure) Width() int {
	return NOTE_WIDTH * m.Beat
}

func (m Measure) Draw(c *svg.SVG) error {
	x1 := m.Base.X
	x2 := m.Base.X + m.Width()

	for i := 0; i < m.Strings; i++ {
		y, err := m.XthStringY(i + 1)
		if err != nil {
			return err
		}
		c.Line(x1, y, x2, y, MEASURE_LINE_DEFINE)
	}
	return nil
}

func (m Measure) XthStringY(xth int) (int, error) {
	if xth > m.Strings {
		return 0, fmt.Errorf("xthが弦の数より多い")
	}
	return m.Base.Y + (xth-1)*SPACE, nil
}

type MeasureBorder struct {
	Measure Measure
}

func (b MeasureBorder) drawStart(c *svg.SVG) error {
	x := b.Measure.Base.X
	y1 := b.Measure.Base.Y
	y2, err := b.Measure.XthStringY(b.Measure.Strings)
	if err != nil {
		return err
	}

	c.Line(x, y1, x, y2, MEASURE_LINE_DEFINE)
	return nil
}

// 運指。小節に基づいて描画位置が決まる
type Fingering struct {
	Measure Measure
	Fret    int // フレット数
	Strings int // 何弦
	Beat    int // 何拍目
}

func (f Fingering) Draw(canvas *svg.SVG) error {
	x := f.Measure.Base.X + (f.Beat-1)*NOTE_WIDTH
	y, err := f.Measure.XthStringY(f.Strings)
	if err != nil {
		return err
	}
	corrected_x := x + (NOTE_WIDTH / 2) // text-anchorによって指定した座標を中心に文字列が表示されるため、文字列幅/2を加える
	corrected_y := y + (SPACE / 2)
	canvas.Text(corrected_x, corrected_y, fmt.Sprint(f.Fret), "text-anchor:middle")
	return nil
}
