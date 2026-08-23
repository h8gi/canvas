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

		// Title
		ctx.SetColor(colornames.White)
		ctx.DrawStringAnchored("Hello, Canvas!", float64(ctx.Width())/2, 60, 0.5, 0.5)

		// Subtitle
		ctx.SetColor(colornames.Gray)
		ctx.DrawStringAnchored("A lightweight 2D animation library for Go", float64(ctx.Width())/2, 90, 0.5, 0.5)

		// Real-time animation stats
		info := fmt.Sprintf("Frame: %d   Time: %.1fs   FPS: %.0f", ctx.FrameCount, ctx.Time, 1.0/ctx.DeltaTime)
		ctx.SetColor(colornames.Cyan)
		ctx.DrawStringAnchored(info, float64(ctx.Width())/2, 140, 0.5, 0.5)

		// Interactive status
		mouseInfo := fmt.Sprintf("Mouse: (%.0f, %.0f)", ctx.Mouse.X, ctx.Mouse.Y)
		ctx.SetColor(colornames.Yellow)
		ctx.DrawStringAnchored(mouseInfo, float64(ctx.Width())/2, 170, 0.5, 0.5)

		// Interactive cursor
		if ctx.IsMouseDragged {
			ctx.SetColor(colornames.Crimson)
		} else {
			ctx.SetColor(colornames.Cornflowerblue)
		}
		ctx.DrawCircle(ctx.Mouse.X, ctx.Mouse.Y, 12)
		ctx.Fill()

		// Footer instruction
		ctx.SetColor(colornames.Darkgray)
		ctx.DrawStringAnchored("Click and drag to interact", float64(ctx.Width())/2, 360, 0.5, 0.5)
	})
}
