package main

import (
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/tbml/args"
	"github.com/tbml/parse"
)

func main() {
	a, err := args.Load()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out_path := a.OutPath()

	file, err := os.Create(out_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	c := svg.New(file)

	cfg := parse.ParseConfig(a.YamlPath)
	p, err := cfg.Build()

	if err != nil {
		panic(err)
	}

	c.Start(p.PageSize())
	defer c.End()

	p.Header.Draw(c)
	for _, l := range p.Score.Lines {
		for _, m := range l.Measures {
			m.Draw(c)
			for _, f := range m.Fingerings {
				f.Draw(c)
				for _, t := range f.Technique {
					t.Draw(c)
				}
			}

			for _, b := range m.Borders {
				b.Draw(c)
			}
		}
	}
}

var SPACE int = 10
var MEASURE_LINE_DEFINE string = "stroke:#bbb;stroke-width:1"
