package main

import (
	"os"
	"testing"

	svg "github.com/ajstarks/svgo"
)

func TestDrawLine(t *testing.T) {
	width := 1000
	height := 1000

	file, err := os.Create("tab_test.svg")
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
	drawLine(canvas)

}
