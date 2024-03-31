package tabsvg

import (
	"fmt"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

var SPACE int = 10
var NOTE_WIDTH int = 20
var MEASURE_LINE_DEFINE string = "stroke:#bbb;stroke-width:1"

type Cordinate struct {
	X int
	Y int
}

type NoteWidth int

type Page struct {
	Base      Cordinate
	Gap       int
	Header    Header
	Score     *Score
	NoteWidth NoteWidth
}

type PageOption struct {
	PageGap  int
	ScoreGap int
}

func NewPage(b Cordinate, gap int, note_width NoteWidth) Page {
	return Page{Base: b, Gap: gap, NoteWidth: note_width}
}

func (p *Page) SetHeader(title string, key string, bpm int) *Header {
	h := Header{Base: p.Base, Title: title, Key: key, BPM: bpm}
	p.Header = h
	return &h
}

var HEADER_HEIGHT = 40

func (p *Page) SetScore(gap int) *Score {
	s := NewScore(Cordinate{X: p.Base.X, Y: p.Base.Y + HEADER_HEIGHT}, gap, p.NoteWidth)
	p.Score = &s
	return &s
}

type Header struct {
	Base  Cordinate
	Title string
	Key   string
	BPM   int
}

var HEADER_FONT_DEFINE = "font-size:25;"
var KEY_FONT_DEFINE = "font-size:12;"

func (h Header) Draw(c *svg.SVG) {
	c.Text(h.Base.X, h.Base.Y, h.Title, HEADER_FONT_DEFINE)

	texts := make([]string, 0, 2)
	if h.Key != "" {
		texts = append(texts, fmt.Sprintf("key of %s", h.Key))
	}

	if h.BPM != 0 {
		bpm := strconv.Itoa(h.BPM)
		texts = append(texts, fmt.Sprintf("♪ = %s", bpm))
	}

	c.Text(h.Base.X, h.Base.Y+20, strings.Join(texts, ", "), KEY_FONT_DEFINE)
}

// 縦方向にLineを並べたものがScore
type Score struct {
	Base      Cordinate
	EndY      int
	Lines     []*Line
	Gap       int
	NoteWidth NoteWidth
}

func NewScore(b Cordinate, gap int, note_width NoteWidth) Score {
	return Score{Base: b, EndY: b.Y, Lines: []*Line{}, Gap: gap, NoteWidth: note_width}
}

func (s *Score) AddNewLine(strings int, with_text bool) *Line {
	l := NewLine(Cordinate{s.Base.X, s.EndY}, strings, with_text, s.NoteWidth)
	s.Lines = append(s.Lines, &l)
	s.EndY = s.EndY + l.Height + s.Gap

	return &l
}

// 横方向にMeasureを並べたもの
type Line struct {
	Base      Cordinate
	EndX      int
	Measures  []*Measure
	Height    int
	Strings   int
	WithText  bool
	NoteWidth NoteWidth
}

var MEASURE_TEXT_HEIGHT int = 10

func NewLine(b Cordinate, s int, with_text bool, note_width NoteWidth) Line {
	h := (s - 1) * SPACE
	if with_text {

		h = h + MEASURE_TEXT_HEIGHT
	}
	return Line{Base: b, EndX: b.X, Measures: []*Measure{}, Strings: s, Height: h, WithText: with_text, NoteWidth: note_width}
}

func (l *Line) AddNewMeasure(beat int, text string) *Measure {
	y := l.Base.Y
	if l.WithText {
		y = y + MEASURE_TEXT_HEIGHT
	}

	m := Measure{Base: Cordinate{l.EndX, y}, Strings: l.Strings, Beat: beat, Text: text, withText: l.WithText, NoteWidth: l.NoteWidth}
	l.Measures = append(l.Measures, &m)
	l.EndX = l.EndX + m.Width()

	return &m
}

type MeasureBorder struct {
	Measure Measure
	Top     Cordinate
	Bottom  Cordinate
}

func (b MeasureBorder) Draw(c *svg.SVG) {
	c.Line(b.Top.X, b.Top.Y, b.Bottom.X, b.Bottom.Y, MEASURE_LINE_DEFINE)
}

type TechniqueInterface interface {
	Draw(*svg.SVG)
}

type LegatoTechnique struct {
	// Legato先とLegato元の中間地点
	Center    Cordinate
	Distance  int
	AfterNote Fingering
	Text      string
	NoteWidth
}

var TECHNIQUE_LINE_DEFINE string = "stroke:#444;stroke-width:1.2"

func (t LegatoTechnique) Draw(c *svg.SVG) {
	text_y := t.Center.Y + SPACE
	c.Text(t.Center.X, text_y, t.Text, FINGERING_TEXT_DEFINE)

	// start_xとend_xの÷3は補正値。t.NoteWidthが極端な値になると文字列と被ったり、文字列から離れすぎたりする
	start_x := t.Center.X - t.Distance/2 + int(t.NoteWidth)/3
	end_x := t.Center.X + t.Distance/2 - int(t.NoteWidth)/3
	c.Line(start_x, t.Center.Y, end_x, t.Center.Y, TECHNIQUE_LINE_DEFINE)

	t.AfterNote.Draw(c)
}
