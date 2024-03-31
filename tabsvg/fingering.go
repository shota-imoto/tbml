package tabsvg

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
)

// 運指。小節に基づいて描画位置が決まる
type Fingering struct {
	Center    Cordinate
	Fret      string // フレット数
	Length    int    // 音長
	Technique []TechniqueInterface

	// 縦方向の補正。
	// svgoのText関数はCordinate.Yを底辺として描画するため補正がないと弦の上にフレット番号が乗って表示されてしまう
	// Fingering.Centerを参照して位置が決まる要素が存在するため、あらかじめCenter.Yに加えるのではなくDraw時に補正する
	CorrectionY int

	// 初期化の際にCenterが格納される想定なのでStringsは正直いらない
	Strings int // 何弦。
}

var FINGERING_TEXT_DEFINE string = "text-anchor:middle;font-size:10px"

func (f Fingering) Draw(c *svg.SVG) {
	c.Text(f.Center.X, f.Center.Y+f.CorrectionY, fmt.Sprint(f.Fret), FINGERING_TEXT_DEFINE)
}

// Fingering動作確認用。文字と座標の中心がずれるためSPACEを変更したときに縦方向の補正値を調整する必要がある
func (f Fingering) DrawCenter(c *svg.SVG) {
	c.Circle(f.Center.X, f.Center.Y, 2)
}

var FINGERING_CORRECTION_Y int = 5

type AddLegatoTechniqueInput struct {
	Fret   string
	Length int
	Text   string
}

func (f *Fingering) AddLegatoTechnique(input AddLegatoTechniqueInput) *LegatoTechnique {
	// Legatoの元の音の中央座標に幅を加えるためNOTE_WIDTH/2を減ずる必要はない
	after_x := f.Center.X + f.Length*NOTE_WIDTH

	after := Fingering{Center: Cordinate{X: after_x, Y: f.Center.Y}, Fret: input.Fret, Strings: f.Strings, Length: input.Length, CorrectionY: FINGERING_CORRECTION_Y, Technique: []TechniqueInterface{}}

	x := (f.Center.X + after_x) / 2
	d := after_x - f.Center.X
	t := LegatoTechnique{Center: Cordinate{X: x, Y: f.Center.Y}, Distance: d, AfterNote: after, Text: input.Text}
	f.Technique = append(f.Technique, &t)
	return &t
}
