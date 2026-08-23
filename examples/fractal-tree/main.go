package main

import (
	"math"

	"github.com/h8gi/canvas"
)

func branch(ctx *canvas.Context, length, thickness float64, depth int, time float64) {
	if depth <= 0 || length < 3 {
		// Draw a glowing flower / leaf at the tip
		hue := math.Mod(330+float64(depth)*15, 360)
		ctx.FillHSBA(hue, 0.7, 0.95, 0.7)
		ctx.DrawCircle(0, 0, 4)
		ctx.Fill()
		return
	}

	// Calculate branch color based on depth
	treeHue := canvas.Map(float64(depth), 0, 10, 80, 25)
	ctx.StrokeHSBA(treeHue, 0.6, 0.8, 0.9)
	ctx.SetLineWidth(thickness)

	ctx.DrawLine(0, 0, 0, -length)
	ctx.Stroke()

	ctx.Translate(0, -length)

	// Wind sway calculation with noise
	wind := canvas.Noise1D(time*0.5+float64(depth)*0.3)*0.4 - 0.2
	branchAngle := 0.45 + math.Sin(time*0.8)*0.05

	// Right branch
	ctx.Push()
	ctx.Rotate(branchAngle + wind)
	branch(ctx, length*0.72, thickness*0.7, depth-1, time)
	ctx.Pop()

	// Left branch
	ctx.Push()
	ctx.Rotate(-branchAngle + wind)
	branch(ctx, length*0.72, thickness*0.7, depth-1, time)
	ctx.Pop()
}

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    540,
		FrameRate: 60,
		Title:     "Fractal Tree with Wind Animation",
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.BackgroundRGB(0.05, 0.05, 0.09)

		ctx.Push()
		// Move origin to bottom center
		ctx.Translate(float64(ctx.Width())/2, float64(ctx.Height())-30)

		branch(ctx, 110.0, 9.0, 9, ctx.Time)
		ctx.Pop()

		ctx.FillRGB(0.7, 0.7, 0.7)
		ctx.DrawString("Organic Fractal Tree (Perlin noise sway)", 15, 25)
	})
}
