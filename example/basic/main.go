//  use https://github.com/fogleman/gg/blob/master/examples/sine.go as refrence
package main

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/h8gi/canvas"
)

func main() {
	c := canvas.New()
	const W = 1200
	const H = 60
	c.Option(
		canvas.Size(W, H),
		canvas.FrameRate(30),
	)
	c.Main(func(dc *gg.Context) {
		dc.SetHexColor("#fff")
		dc.Clear()
		dc.ScaleAbout(0.95, 0.75, W/2, H/2)
		for i := 0; i < W; i++ {
			a := float64(i) * 2 * math.Pi / W * 8
			x := float64(i)
			y := (math.Sin(a) + 1) / 2 * H
			dc.LineTo(x, y)
		}
		dc.ClosePath()
		dc.SetHexColor("#3E606F")
		dc.FillPreserve()
		dc.SetHexColor("#19344180")
		dc.SetLineWidth(8)
		dc.Stroke()
	})
}
