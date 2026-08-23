// draw context
package canvas

import (
	"image"
	"sync"

	"github.com/fogleman/gg"
	"github.com/gopxl/pixel/v2"
)

type Context struct {
	gg.Context
	mu               sync.Mutex
	IsMouseDragged   bool
	Mouse            pixel.Vec
	PMouse           pixel.Vec
	pressed          func(pixel.Button) bool
	flippedPixBuffer []uint8
}

func NewContext(width, height int) *Context {
	var mu sync.Mutex
	return &Context{
		Context:          *gg.NewContext(width, height),
		mu:               mu,
		IsMouseDragged:   false,
		Mouse:            pixel.Vec{X: 0, Y: 0},
		PMouse:           pixel.Vec{X: 0, Y: 0},
		pressed:          func(pixel.Button) bool { return true },
		flippedPixBuffer: make([]uint8, width*height*4),
	}
}

func (ctx *Context) pix() []uint8 {
	return ctx.Image().(*image.RGBA).Pix
}

func (ctx *Context) flippedPix() []uint8 {
	flipV(ctx.pix(), ctx.flippedPixBuffer, ctx.Width(), ctx.Height())
	return ctx.flippedPixBuffer
}

func (ctx *Context) IsKeyPressed(b pixel.Button) bool {
	return ctx.pressed(b)
}

