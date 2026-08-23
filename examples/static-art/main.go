package main

import (
	"math/rand"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    400,
		FrameRate: 30,
		Title:     "Static Generative Art (NoLoop Example)",
		Resizable: true,
	})

	c.Setup(func(ctx *canvas.Context) {
		// Stop continuous frame rendering after the first frame
		c.NoLoop()
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.BackgroundHex("#181825")

		// Draw generative geometric circles
		for i := 0; i < 40; i++ {
			x := canvas.Random(50, float64(ctx.Width())-50)
			y := canvas.Random(50, float64(ctx.Height())-50)
			radius := canvas.Random(10, 60)

			ctx.FillRGBA(rand.Float64(), rand.Float64(), rand.Float64(), 0.6)
			ctx.StrokeHex("#cdd6f4")
			ctx.SetLineWidth(1.5)
			ctx.DrawCircle(x, y, radius)
			ctx.Fill()
			ctx.Stroke()
		}

		// Instructions
		ctx.FillColor(colornames.White)
		ctx.DrawStringAnchored("Generative Static Art (Rendered once with NoLoop)", float64(ctx.Width())/2, 370, 0.5, 0.5)
	})
}
