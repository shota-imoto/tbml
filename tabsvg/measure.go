package tabsvg

import (
	"fmt"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

// 小節
type Measure struct {
	Base       Cordinate // 小節の左上を0点とする
	Strings    int       // 弦の数
	Beat       int       // 拍数
	Text       string    // 小節ごとのメモ
	withText   bool
	sumLength  int
	Fingerings []*Fingering
}

func (m Measure) SumLength() int {
	return m.sumLength
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

func (m *Measure) AddFingering(fret string, strings, length int) (*Fingering, error) {

	x := m.Base.X + m.SumLength()*NOTE_WIDTH

	y, err := m.XthStringY(strings)
	if err != nil {
		return &Fingering{}, err
	}
	center_x := x + (NOTE_WIDTH / 2)
	center_y := y

	f := Fingering{Center: Cordinate{center_x, center_y}, CorrectionY: FINGERING_CORRECTION_Y, Length: length, Fret: fret, Strings: strings}
	m.Fingerings = append(m.Fingerings, &f)
	m.sumLength += length
	return &f, nil
}

// 重音を表現する場合に２要素以上のinputsを受け取る
func (m *Measure) AddFingerings(length int, inputs ...FingeringInput) ([]*Fingering, error) {
	fingerings := []*Fingering{}
	for _, input := range inputs {
		b := builder(*m, input, length)
		f, err := b.build()
		if err != nil {
			return []*Fingering{}, fmt.Errorf("AddFingerings is failed: %v", err)
		}
		m.Fingerings = (append(m.Fingerings, f))
		fingerings = append(fingerings, f)
	}
	m.sumLength += length

	return fingerings, nil
}

type FingeringInput struct {
	Fret       string
	Strings    string
	Techniques []AddLegatoTechniqueInput
}

func builder(m Measure, input FingeringInput, length int) FingeringBuilderI {
	if input.Fret != "x" {
		return FingeringBuilder{input: input, length: length, measure: m}

	}

	if input.Strings != "x" {
		return RestBuilder{input: input, length: length, measure: m}
	} else {
		return WhiteSpaceBuilder{}
	}
}

type FingeringBuilderI interface {
	build() (*Fingering, error)
}

type WhiteSpaceBuilder struct{}

func (b WhiteSpaceBuilder) build() (*Fingering, error) {
	return &Fingering{}, nil
}

type RestBuilder struct {
	measure Measure
	input   FingeringInput
	length  int
}

func (b RestBuilder) build() (*Fingering, error) {
	x := b.measure.Base.X + b.measure.SumLength()*NOTE_WIDTH

	s, err := strconv.Atoi(b.input.Strings)
	if err != nil {
		return &Fingering{}, fmt.Errorf("AddFingerings is failed: %v", err)
	}

	y, err := b.measure.XthStringY(s)
	if err != nil {
		return &Fingering{}, fmt.Errorf("AddFingerings is failed: %v", err)
	}
	center_x := x + (NOTE_WIDTH / 2)

	return &Fingering{Center: Cordinate{center_x, y}, CorrectionY: FINGERING_CORRECTION_Y, Length: b.length, Fret: b.input.Fret, Strings: s}, nil
}

type FingeringBuilder struct {
	measure Measure
	input   FingeringInput
	length  int
}

func (b FingeringBuilder) build() (*Fingering, error) {
	s, err := strconv.Atoi(b.input.Strings)

	if err != nil {
		return &Fingering{}, fmt.Errorf("AddFingerings is failed: %v", err)
	}

	y, err := b.measure.XthStringY(s)
	if err != nil {
		return &Fingering{}, fmt.Errorf("AddFingerings is failed: %v", err)
	}
	x := b.measure.Base.X + b.measure.SumLength()*NOTE_WIDTH
	center_x := x + (NOTE_WIDTH / 2)
	f := Fingering{Center: Cordinate{center_x, y}, CorrectionY: FINGERING_CORRECTION_Y, Length: b.length, Fret: b.input.Fret, Strings: s}
	for _, tech_input := range b.input.Techniques {
		tech := f.AddLegatoTechnique(AddLegatoTechniqueInput{Fret: tech_input.Fret, Length: b.length, Text: tech_input.Text})
		f.Technique = append(f.Technique, tech)
	}
	return &f, nil
}
