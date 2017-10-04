package main

import (
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.New(&canvas.NewCanvasOptions{
		Width:     600,
		Height:    400,
		FrameRate: 30,
		Title:     "hello canvas!",
	})

	c.Setup(func(ctx *canvas.Context) {
		ctx.SetColor(colornames.White)
		ctx.Clear()
		ctx.SetColor(colornames.Green)
	})

	c.Draw(func(ctx *canvas.Context) {
		if ctx.MouseDragged() {
			ctx.DrawCircle(ctx.MouseX(), ctx.MouseY(), 5)
			ctx.Fill()
		}
		if ctx.KeyPressed() {
			ctx.SetColor(colornames.White)
			ctx.Clear()
			ctx.SetColor(colornames.Green)
		}
	})
}
