package tabsvg

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
)

var SPACE int = 10
var NOTE_WIDTH int = 20
var MEASURE_LINE_DEFINE string = "stroke:#bbb;stroke-width:1"

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

func NewScore(b Cordinate, gap int) Score {
	fmt.Println("tabsvg package")
	return Score{Base: b, EndY: b.Y, Lines: []Line{}, Gap: gap}
}

func (s *Score) AddNewLine(strings int, with_text bool) Line {
	l := NewLine(Cordinate{s.Base.X, s.EndY}, strings, with_text)
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
	WithText bool
}

var MEASURE_TEXT_HEIGHT int = 10

func NewLine(b Cordinate, s int, with_text bool) Line {
	h := (s - 1) * SPACE
	if with_text {

		h = h + MEASURE_TEXT_HEIGHT
	}
	return Line{Base: b, EndX: b.X, Measures: []Measure{}, Strings: s, Height: h, WithText: with_text}
}

func (l *Line) AddNewMeasure(beat int, text string) Measure {
	y := l.Base.Y

	if l.WithText {
		y = y + MEASURE_TEXT_HEIGHT
	}

	m := Measure{Base: Cordinate{l.EndX, y}, Strings: l.Strings, Beat: beat, Text: text, withText: l.WithText}
	l.Measures = append(l.Measures, m)
	l.EndX = l.EndX + m.Width()

	return m
}

// 小節
type Measure struct {
	Base     Cordinate // 小節の左上を0点とする
	Strings  int       // 弦の数
	Beat     int       // 拍数
	Text     string    // 小節ごとのメモ
	withText bool
}

func (m Measure) Width() int {
	return NOTE_WIDTH * m.Beat
}

func (m Measure) Draw(c *svg.SVG) error {
	x1 := m.Base.X
	x2 := m.Base.X + m.Width()

	// テキストの描画
	if m.withText {
		c.Text(x1, m.Base.Y-MEASURE_TEXT_HEIGHT, m.Text)
	}

	// 譜面の描画
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

var FINGERING_CORRECTION_Y int = 5

func (m *Measure) AddFingering(fret, strings, beat int) (Fingering, error) {

	x := m.Base.X + (beat-1)*NOTE_WIDTH

	y, err := m.XthStringY(strings)
	if err != nil {
		return Fingering{}, err
	}
	center_x := x + (NOTE_WIDTH / 2)
	center_y := y
	return Fingering{Center: Cordinate{center_x, center_y}, CorrectionY: FINGERING_CORRECTION_Y, Fret: fret, Strings: strings}, nil
}

type MeasureBorder struct {
	Measure Measure
}

func (b MeasureBorder) DrawStart(c *svg.SVG) error {
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
	Center Cordinate
	Fret   int // フレット数

	// 縦方向の補正。
	// svgoのText関数はCordinate.Yを底辺として描画するため補正がないと弦の上にフレット番号が乗って表示されてしまう
	// Fingering.Centerを参照して位置が決まる要素が存在するため、あらかじめCenter.Yに加えるのではなくDraw時に補正する
	CorrectionY int

	// 以下は、初期化の際にCenterが格納される想定なので正直いらない
	Strings int // 何弦。
	Beat    int // 何拍目。
}

var FINGERING_TEXT_DEFINE string = "text-anchor:middle"

func (f Fingering) Draw(c *svg.SVG) {
	c.Text(f.Center.X, f.Center.Y+f.CorrectionY, fmt.Sprint(f.Fret), FINGERING_TEXT_DEFINE)
}

// Fingering動作確認用。文字と座標の中心がずれるためSPACEを変更したときに縦方向の補正値を調整する必要がある
func (f Fingering) DrawCenter(c *svg.SVG) {
	c.Circle(f.Center.X, f.Center.Y, 2)
}

type Technique struct {
	Start Fingering
	End   Fingering
	Text  string
}

var TECHNIQUE_LINE_DEFINE string = "stroke:#444;stroke-width:1.2"

func (t Technique) Draw(c *svg.SVG) error {
	x := (t.Start.Center.X + t.End.Center.X) / 2
	if t.Start.Center.Y != t.End.Center.Y {
		// Center計算時の小数点以下の扱いによってはバグるかも
		return fmt.Errorf("Technique.Draw is failed: 別の弦なのでおかしい")
	}
	y := t.Start.Center.Y + SPACE

	c.Text(x, y, t.Text, FINGERING_TEXT_DEFINE)

	// start_xとend_xの÷3は補正値。NOTE_WIDTHが極端な値になると文字列と被ったり、文字列から離れすぎたりする
	start_x := t.Start.Center.X + NOTE_WIDTH/3
	end_x := t.End.Center.X - NOTE_WIDTH/3

	c.Line(start_x, t.Start.Center.Y, end_x, t.End.Center.Y, TECHNIQUE_LINE_DEFINE)
	return nil
}
