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
			}
		}
	}
}

func (c Config) Build() (tabsvg.Score, error) {
	cordinate, err := parseCordinate(c.Score.Cordinate)
	if err != nil {
		return tabsvg.Score{}, fmt.Errorf("Build is failed: %v", err)
	}
	s := tabsvg.NewScore(cordinate, c.Score.Gap)

	for _, l := range c.Score.Lines {

		new_line := s.AddNewLine(l.Strings, l.WithText)

		// TODO: ネストが深すぎる。
		// たぶんループは見やすいよう切り出して、パフォーマンスは非同期化で回避するのが良さそう
		for _, m := range l.Measures {
			_ = new_line.AddNewMeasure(m.Beat, m.Text)

		}
	}
	return s, nil
}

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
