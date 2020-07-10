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
		// first push
		ctx.Push()
		ctx.SetColor(colornames.Green)
		ctx.RotateAbout(50, 0, 0)
		ctx.DrawRectangle(6, 6, 30, 30)
		ctx.DrawRectangle(40, 6, 30, 30)
		ctx.Fill()
		ctx.Pop() // first pop

		// second push
		ctx.Push()
		ctx.RotateAbout(120, 100, 100)
		ctx.SetColor(colornames.Red)
		ctx.DrawRectangle(100, 100, 50, 50)
		ctx.Fill()

		ctx.SetColor(colornames.Blue)
		ctx.DrawRectangle(120, 80, 50, 50)
		ctx.Fill()
		ctx.Pop() // second pop

		ctx.DrawRectangle(120, 80, 50, 50)
		ctx.Fill()
	})

	c.Draw(func(ctx *canvas.Context) {
	})
}
