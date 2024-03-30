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
	Borders    []*MeasureBorder
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

	// 重音なので一回だけしかlengthを加えない。forループの外で加える
	m.sumLength += length

	return fingerings, nil
}

type FingeringInput struct {
	Fret       string
	Strings    string
	Techniques []AddLegatoTechniqueInput
}

type FingeringBuilderI interface {
	build() (*Fingering, error)
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

// 空白スペース用のbuilder
type WhiteSpaceBuilder struct{}

func (b WhiteSpaceBuilder) build() (*Fingering, error) {
	return &Fingering{}, nil
}

// 休符用のbuilder
type RestBuilder struct {
	measure Measure
	input   FingeringInput
	length  int
}

func (b RestBuilder) build() (*Fingering, error) {
	x := b.measure.Base.X + b.measure.SumLength()*NOTE_WIDTH

	s, err := strconv.Atoi(b.input.Strings)
	if err != nil {
		return &Fingering{}, fmt.Errorf("RestBuilder.build is failed: %v", err)
	}

	y, err := b.measure.XthStringY(s)
	if err != nil {
		return &Fingering{}, fmt.Errorf("RestBuilder.build is failed: %v", err)
	}
	center_x := x + (NOTE_WIDTH / 2)

	return &Fingering{Center: Cordinate{center_x, y}, CorrectionY: FINGERING_CORRECTION_Y, Length: b.length, Fret: b.input.Fret, Strings: s}, nil
}

// 一般的な音符のbuilder
type FingeringBuilder struct {
	measure Measure
	input   FingeringInput
	length  int
}

func (b FingeringBuilder) build() (*Fingering, error) {
	s, err := strconv.Atoi(b.input.Strings)

	if err != nil {
		return &Fingering{}, fmt.Errorf("FingeringBuilder.build is failed: %v", err)
	}

	y, err := b.measure.XthStringY(s)
	if err != nil {
		return &Fingering{}, fmt.Errorf("FingeringBuilder.build: %v", err)
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

type BorderPositionI interface {
	xPosition(Measure) int
}

type StartPosition struct{}

func (p StartPosition) xPosition(m Measure) int {
	return m.Base.X
}

type EndPosition struct{}

func (p EndPosition) xPosition(m Measure) int {
	return m.Base.X + m.Width()
}

func (m *Measure) AddBorder(position BorderPositionI) (*MeasureBorder, error) {
	x := position.xPosition(*m)
	y1 := m.Base.Y
	y2, err := m.XthStringY(m.Strings)
	if err != nil {
		return &MeasureBorder{}, fmt.Errorf("AddBorder is failed: %v", err)
	}

	b := &MeasureBorder{Measure: *m, Top: Cordinate{X: x, Y: y1}, Bottom: Cordinate{X: x, Y: y2}}
	m.Borders = append(m.Borders, b)
	return b, nil
}
