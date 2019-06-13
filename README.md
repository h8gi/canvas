# Canvas

`canvas` is 2d animation library.

## Usage

Create canvas object.

```go
c := canvas.NewCanvas(&canvas.CanvasConfig{
	Width: 300,
	Height: 300,
	FrameRate: 60,
})
```

Set drawing function and start loop.

```go
c.Draw(func(ctx *canvas.Context) {
	if ctx.IsMouseDragged {
		ctx.DrawCircle(ctx.Mouse.X, ctx.Mouse.Y, 5)
		ctx.Fill()
	}
})
```

Struct `gg.Context` is embedded in `canvas.Context`.
See [https://github.com/fogleman/gg](https://github.com/fogleman/gg) about details.

## Example

See `example` directory.

```go
package main

import (
	"github.com/faiface/pixel/pixelgl"
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
		ctx.SetColor(colornames.White)
		ctx.Clear()
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

		if ctx.IsKeyPressed(pixelgl.KeyUp) {
			ctx.Push()
			ctx.SetColor(colornames.White)
			ctx.Clear()
			ctx.Pop()
		}
	})
}
``` 

## Built With

- [gg](https://github.com/fogleman/gg) - 2D graphics library.
- [pixel](https://github.com/faiface/pixel) - 2D game library.
