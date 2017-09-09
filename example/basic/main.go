package main

import (
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.New(nil)
	c.Setup(func(ctx *canvas.Context) {
		ctx.SetColor(colornames.Green)
	})
	c.Draw(func(ctx *canvas.Context) {
		if ctx.MouseDragged() {
			ctx.DrawCircle(ctx.MouseX(), ctx.MouseY(), 5)
			ctx.Fill()
		}
	})
}
