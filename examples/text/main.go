package main

import (
	"fmt"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    400,
		FrameRate: 60,
		Title:     "Text Rendering Example",
	})

	c.Setup(func(ctx *canvas.Context) {
		ctx.SetFontFace(basicfont.Face7x13)
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.BackgroundHex("#1e1e2e")

		// Draw header text
		ctx.SetColor(colornames.White)
		ctx.DrawStringAnchored("Hello Canvas! (Processing-like 2D Library)", float64(ctx.Width())/2, 50, 0.5, 0.5)

		// Draw frame and time info
		info := fmt.Sprintf("Frame: %d | Time: %.2fs | FPS: %.1f", ctx.FrameCount, ctx.Time, 1.0/ctx.DeltaTime)
		ctx.SetColor(colornames.Cyan)
		ctx.DrawStringAnchored(info, float64(ctx.Width())/2, 90, 0.5, 0.5)

		// Draw mouse position and instructions
		mouseInfo := fmt.Sprintf("Mouse: (%.0f, %.0f) [Origin: Top-Left]", ctx.Mouse.X, ctx.Mouse.Y)
		ctx.SetColor(colornames.Yellow)
		ctx.DrawStringAnchored(mouseInfo, float64(ctx.Width())/2, 130, 0.5, 0.5)

		// Interactive box following mouse
		ctx.SetColor(colornames.Crimson)
		ctx.DrawCircle(ctx.Mouse.X, ctx.Mouse.Y, 15)
		ctx.Fill()

		// Key instruction
		ctx.SetColor(colornames.Lightgray)
		ctx.DrawStringAnchored("Move mouse around. Text renders correctly from top to bottom.", float64(ctx.Width())/2, 350, 0.5, 0.5)
	})
}
