package main

import (
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.FyneNew(&canvas.NewCanvasOptions{
		Width:     600,
		Height:    400,
		FrameRate: 30,
		Title:     "hello canvas!",
	})

	c.Setup(func(ctx *canvas.Context) {
		ctx.SetColor(colornames.White)
		ctx.Clear()
		ctx.SetColor(colornames.Green)
		ctx.SetLineWidth(5)
	})

	c.Draw(func(ctx *canvas.Context) {
		if ctx.IsMouseDragged() {
			ctx.DrawLine(ctx.MouseX(), ctx.MouseY(),
				ctx.PreviousMouseX(), ctx.PreviousMouseY())
			ctx.Stroke()
		}

		if ctx.IsKeyPressed() {
			ctx.Push()
			ctx.SetColor(colornames.White)
			ctx.Clear()
			ctx.Pop()
		}
	})
}
