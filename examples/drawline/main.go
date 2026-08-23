package main

import (
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    400,
		FrameRate: 30,
		Title:     "Hello Canvas!",
	})

	c.Setup(func(ctx *canvas.Context) {
		ctx.Background(colornames.White)
		ctx.SetColor(colornames.Green)
		ctx.SetLineWidth(5)
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.Push()
		if ctx.IsMouseDragged {
			ctx.SetColor(colornames.Red)
		}
		ctx.DrawLine(ctx.Mouse.X, ctx.Mouse.Y,
			ctx.PMouse.X, ctx.PMouse.Y)
		ctx.Stroke()
		ctx.Pop()

		if ctx.IsKeyPressed(canvas.KeyUp) {
			ctx.Background(colornames.White)
		}
	})
}
