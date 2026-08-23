# Canvas

[![CI](https://github.com/h8gi/canvas/actions/workflows/ci.yml/badge.svg)](https://github.com/h8gi/canvas/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/h8gi/canvas.svg)](https://pkg.go.dev/github.com/h8gi/canvas)

`canvas` is a lightweight, Processing-like 2D animation and creative coding library for Go.

It combines the power and simplicity of [`fogleman/gg`](https://github.com/fogleman/gg) (2D vector drawing) with [`gopxl/pixel`](https://github.com/gopxl/pixel) for window management and real-time graphics rendering.

## Features

- **Processing / p5.js-like workflow**: Familiar `Setup` and `Draw` loop structure.
- **Top-left origin coordinate system**: Natural 2D coordinate system matching standard Canvas & Processing conventions.
- **Embedded `gg.Context`**: Full 2D vector drawing capabilities (shapes, curves, paths, colors, text, transformations).
- **Zero-dependency user API**: Built-in keyboard/mouse keys and helpers (`canvas.KeyUp`, `canvas.MouseLeft`, `canvas.V(x, y)`).
- **Time & Frame helpers**: `ctx.FrameCount`, `ctx.DeltaTime`, `ctx.Time`.
- **Creative coding math utilities**: `canvas.Map`, `canvas.Lerp`, `canvas.Constrain`, `canvas.Dist`, `canvas.Random`, `canvas.Radians`, `canvas.Degrees`.

## Installation

```sh
go get github.com/h8gi/canvas
```

## Quick Start

```go
package main

import (
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    400,
		FrameRate: 60,
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
		ctx.DrawLine(ctx.Mouse.X, ctx.Mouse.Y, ctx.PMouse.X, ctx.PMouse.Y)
		ctx.Stroke()
		ctx.Pop()

		if ctx.IsKeyPressed(canvas.KeyUp) {
			ctx.Background(colornames.White)
		}
	})
}
```

## Examples

Check the [examples](examples) directory for complete programs:

- [`drawline`](examples/drawline): Interactive mouse drawing and key interactions.
- [`lifegame`](examples/lifegame): Conway's Game of Life simulation.
- [`rotate-objects`](examples/rotate-objects): Coordinate system rotation and object transformations.
- [`text`](examples/text): Text rendering, time tracking, and frame rate display.

## Built With

- [gg](https://github.com/fogleman/gg) - 2D graphics library
- [pixel](https://github.com/gopxl/pixel) - 2D OpenGL game engine

