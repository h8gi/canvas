package main

import (
	"fmt"
	"math"

	"github.com/h8gi/canvas"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     400,
		Height:    400,
		FrameRate: 30,
		Title:     "GIF Recording Example",
	})

	totalFrames := 60 // 2 seconds at 30 fps

	c.Setup(func(ctx *canvas.Context) {
		fmt.Printf("Recording %d frames to animation.gif...\n", totalFrames)

		c.RecordGIF("animation.gif", totalFrames, &canvas.GIFOptions{
			LoopCount: 0, // Infinite loop
			OnComplete: func(path string, err error) {
				if err != nil {
					fmt.Printf("Failed to save GIF: %v\n", err)
				} else {
					fmt.Printf("Successfully recorded GIF to %s!\n", path)
				}
			},
		})
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.BackgroundRGB(0.08, 0.08, 0.12)

		t := float64(ctx.FrameCount) / float64(totalFrames) * math.Pi * 2
		cx, cy := float64(ctx.Width())/2, float64(ctx.Height())/2
		numPetals := 8

		ctx.Push()
		ctx.SetLineWidth(2.5)

		for i := 0; i < numPetals; i++ {
			angle := float64(i)*(2*math.Pi/float64(numPetals)) + t
			dist := 70.0 + math.Sin(t+float64(i))*25.0
			x := cx + math.Cos(angle)*dist
			y := cy + math.Sin(angle)*dist

			hue := math.Mod(float64(i)*(360.0/float64(numPetals))+float64(ctx.FrameCount)*4, 360)
			ctx.StrokeHSBA(hue, 0.8, 0.95, 0.9)
			ctx.FillHSBA(hue, 0.8, 0.95, 0.3)

			ctx.DrawCircle(x, y, 24+math.Sin(t*2+float64(i))*8)
			ctx.Fill()
			ctx.Stroke()
		}

		ctx.Pop()

		// Draw progress text
		if ctx.FrameCount <= totalFrames {
			ctx.FillRGB(0.8, 0.8, 0.8)
			ctx.DrawString(fmt.Sprintf("Recording: %d / %d", ctx.FrameCount, totalFrames), 15, 30)
		} else {
			ctx.FillRGB(0.3, 0.9, 0.4)
			ctx.DrawString("Recording complete!", 15, 30)
		}
	})
}
