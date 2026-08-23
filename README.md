# Canvas

[![CI](https://github.com/h8gi/canvas/actions/workflows/ci.yml/badge.svg)](https://github.com/h8gi/canvas/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/h8gi/canvas.svg)](https://pkg.go.dev/github.com/h8gi/canvas)

`canvas` is a lightweight, Processing-like 2D animation and creative coding library for Go.

It combines the power and simplicity of [`fogleman/gg`](https://github.com/fogleman/gg) (2D vector drawing) with [`gopxl/pixel`](https://github.com/gopxl/pixel) for window management and real-time graphics rendering.

## Features

- **Simple & Intuitive**: Start animating with familiar `Setup` and `Draw` callbacks—no window or loop boilerplate.
- **Rich 2D Vector Drawing**: Full 2D vector graphics (shapes, paths, text, images, colors, and transformations) powered by [`gg`](https://github.com/fogleman/gg).
- **Interactive**: Built-in mouse and keyboard handling for interactive sketches, simulations, and creative coding.
- **Easy Recording & Export**: Built-in automated animated GIF recording (`RecordGIF`) and sequential PNG frame capture (`SaveFrameSeq`).

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

## Showcase

| [**Perlin Flow Field**](examples/flow-field) | [**Boids Flocking**](examples/boids) | [**Particle Fountain**](examples/particle-system) |
| :---: | :---: | :---: |
| [![Flow Field](assets/preview-flowfield.png)](examples/flow-field) | [![Boids](assets/preview-boids.png)](examples/boids) | [![Particles](assets/preview-particles.png)](examples/particle-system) |
| [**Fractal Tree**](examples/fractal-tree) | [**Noise Landscape**](examples/noise-landscape) | [**Generative Grid**](examples/static-art) |
| [![Fractal Tree](assets/preview-fractal-tree.png)](examples/fractal-tree) | [![Noise Landscape](assets/preview-noise-landscape.png)](examples/noise-landscape) | [![Static Art](assets/preview-static-art.png)](examples/static-art) |

Explore all complete, runnable programs in the [examples](examples) directory.

## Built With

- [gg](https://github.com/fogleman/gg) - 2D graphics library
- [pixel](https://github.com/gopxl/pixel) - 2D OpenGL game engine

## License

[MIT License](LICENSE) © 2019-2026 Hiroki Yagi (h8gi)


