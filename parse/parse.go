package parse

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Score struct {
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

func ParseConfig(filename string) Score {
	bytes, err := os.ReadFile("tab.yaml")
	if err != nil {
		panic(err)
	}
	return ParseYaml(bytes)
}

func ParseYaml(yaml_byte []byte) Score {
	s := Score{}
	err := yaml.Unmarshal(yaml_byte, &s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("--- t:\n%v\n\n", s)
	fmt.Printf("--- t:\n%v\n\n", s.Score.Lines[1].Measures[1].Fingering)
	return s
}
