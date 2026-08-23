package main

import (
	"fmt"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    400,
		FrameRate: 60,
		Title:     "Perlin Noise Wave Animation",
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.BackgroundHex("#0f172a")

		// Draw multiple undulating noise waves
		numWaves := 6
		for w := 0; w < numWaves; w++ {
			baseY := float64(ctx.Height())*0.4 + float64(w)*35
			colorAlpha := 0.25 + float64(w)*0.12

			ctx.FillRGBA(0.2, 0.6+float64(w)*0.06, 0.9, colorAlpha)
			ctx.StrokeHex("#38bdf8")
			ctx.SetLineWidth(1.5)

			ctx.MoveTo(0, float64(ctx.Height()))
			ctx.LineTo(0, baseY)

			step := 8.0
			for x := 0.0; x <= float64(ctx.Width()); x += step {
				// Perlin noise calculation with time offset
				n := canvas.Noise(x*0.005+float64(w)*0.8, ctx.Time*0.4+float64(w)*0.2, float64(w))
				yOffset := (n - 0.5) * 90.0
				ctx.LineTo(x, baseY+yOffset)
			}

			ctx.LineTo(float64(ctx.Width()), float64(ctx.Height()))
			ctx.ClosePath()
			ctx.Fill()
		}

		// Text info using built-in default font
		ctx.FillColor(colornames.White)
		ctx.DrawStringAnchored("Perlin Noise Wave Animation", 20, 30, 0, 0.5)

		info := fmt.Sprintf("Frame: %d | Time: %.2fs", ctx.FrameCount, ctx.Time)
		ctx.FillHex("#94a3b8")
		ctx.DrawStringAnchored(info, 20, 50, 0, 0.5)
	})
}
