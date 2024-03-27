package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/tbml/parse"
)

func main() {
	width := 1000
	height := 1000

	file, err := os.Create("tab.svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	canvas := svg.New(file)

	canvas.Start(width, height)
	defer canvas.End()

	// タイトル
	// canvas.Text(10, 20, "My Original Music Score", "font-size:14px;fill:black")

	// 楽譜を描画
	// drawLine(canvas)
	c := parse.ParseConfig("tab.yaml")
	s, _ := c.Build()

	for _, l := range s.Lines {
		for _, m := range l.Measures {
			m.Draw(canvas)
			for _, f := range m.Fingerings {
				f.Draw(canvas)
				for _, t := range f.Technique {
					t.Draw(canvas)
				}
			}
		}
	}

}

var SPACE int = 10
var NOTE_WIDTH int = 20
var MEASURE_LINE_DEFINE string = "stroke:#bbb;stroke-width:1"
