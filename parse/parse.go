package parse

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tbml/tabsvg"
	"gopkg.in/yaml.v3"
)

type Fingering string

type Config struct {
	Score struct {
		Cordinate string
		Gap       int
		Lines     []struct {
			Strings  int
			WithText bool `yaml:"with_text"`
			Measures []struct {
				Beat      int
				Text      string
				Fingering []string
				Borders   []string
			}
		}
	}
}

var PositionMap = map[string]tabsvg.BorderPositionI{
	"start": tabsvg.StartPosition{},
	"end":   tabsvg.EndPosition{},
}

func (c Config) Build() (*tabsvg.Score, error) {
	cordinate, err := parseCordinate(c.Score.Cordinate)
	if err != nil {
		return &tabsvg.Score{}, fmt.Errorf("Build is failed: %v", err)
	}
	s := tabsvg.NewScore(cordinate, c.Score.Gap)

	for _, l := range c.Score.Lines {

		new_line := s.AddNewLine(l.Strings, l.WithText)

		// TODO: ネストが深すぎる。
		// たぶんループは見やすいよう切り出して、パフォーマンスは非同期化で回避するのが良さそう
		for _, m := range l.Measures {

			new_measure := new_line.AddNewMeasure(m.Beat, m.Text)
			for _, f := range m.Fingering {
				inputs, length, _ := ParseFingering(f)

				_, _ = new_measure.AddFingerings(length, inputs...)
			}

			for _, b := range m.Borders {
				p := PositionMap[b]
				if p == nil {
					return &tabsvg.Score{}, fmt.Errorf("Build is failed: bordersの値が不正です。-> %s", b)
				}
				new_measure.AddBorder(p)
			}

		}
	}
	return &s, nil
}

// TODO: パッケージから参照できないようにする
func ParseFingering(f string) ([]tabsvg.FingeringInput, int, error) {
	split := strings.Split(f, " ")
	if len(split) > 1 {
		l, err := strconv.Atoi(split[1])

		if err != nil {
			return []tabsvg.FingeringInput{}, 0, fmt.Errorf("ParseFingering is failed: parsed error %v", err)
		}
		inputs, err := parseFingeringStr(split[0], l)

		if err != nil {
			return inputs, l, fmt.Errorf("ParseFingering is failed: parsed error %v", err)
		}

		return inputs, l, nil

	} else if len(split) == 1 {
		l := 1
		inputs, err := parseFingeringStr(split[0], l)

		if err != nil {
			return []tabsvg.FingeringInput{}, l, fmt.Errorf("ParseFingering is failed: parsed error %v", err)
		}

		return inputs, 1, nil
	} else {
		return []tabsvg.FingeringInput{}, 0, fmt.Errorf("parseFingering: split error %v", split)
	}
}

func parseFingeringStr(strs string, length int) ([]tabsvg.FingeringInput, error) {
	f_str_ary := strings.Split(strs, "/")
	fis := []tabsvg.FingeringInput{}

	for _, f_str := range f_str_ary {
		f_split := strings.Split(f_str, ".")
		fret := f_split[0]
		strings := f_split[1]

		fis = append(fis, buildFingering(fret, strings, length))
	}
	return fis, nil
}

var LEGATO_TECHNIQUE = []string{"s", "h", "p"}

// 2s4とか3p2とかをparseしてFingeringInputを組み立てる。sは弦数
func buildFingering(fret_str string, s string, length int) tabsvg.FingeringInput {
	before := ""
	text := ""
	after := ""
	for _, t := range LEGATO_TECHNIQUE {
		split := strings.Split(fret_str, t)
		if len(split) == 2 {
			before = split[0]
			after = split[1]
			text = t
			break
		}
	}
	if before != "" {
		return tabsvg.FingeringInput{
			Strings:    s,
			Fret:       before,
			Techniques: []tabsvg.AddLegatoTechniqueInput{{Fret: after, Length: length, Text: text}}}
	}

	return tabsvg.FingeringInput{
		Strings: s,
		Fret:    fret_str,
	}
}

func (f Fingering) AddToMeasure(tabsvg.Measure) {}

func parseCordinate(str string) (tabsvg.Cordinate, error) {
	separeted := strings.Split(str, ",")
	if len(separeted) != 2 {
		return tabsvg.Cordinate{}, fmt.Errorf("parseCordinate is failed: 座標のフォーマットの誤り： %v", separeted)
	}
	x, err := strconv.Atoi(separeted[0])
	if err != nil {
		return tabsvg.Cordinate{}, fmt.Errorf("parseCordinate is failed: separeted[0] is not integer = %s", separeted[0])
	}
	y, err := strconv.Atoi(separeted[1])
	if err != nil {
		return tabsvg.Cordinate{}, fmt.Errorf("parseCordinate is failed: separeted[1] is not integer = %s", separeted[1])
	}

	return tabsvg.Cordinate{X: x, Y: y}, nil
}

func ParseConfig(filename string) Config {
	bytes, err := os.ReadFile("tab.yaml")
	if err != nil {
		panic(err)
	}
	return ParseYaml(bytes)
}

func ParseYaml(yaml_byte []byte) Config {
	s := Config{}
	err := yaml.Unmarshal(yaml_byte, &s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return s
}
