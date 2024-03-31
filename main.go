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
	width := 1000
	height := 1000

	file, err := os.Create(out_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	c := svg.New(file)

	c.Start(width, height)
	defer c.End()

	// タイトル
	// c.Text(10, 20, "My Original Music Score", "font-size:14px;fill:black")

	// 楽譜を描画
	// drawLine(c)
	cfg := parse.ParseConfig("tab.yaml")
	p, err := cfg.Build()

	if err != nil {
		panic(err)
	}
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
