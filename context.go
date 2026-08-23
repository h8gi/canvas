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
	mu             sync.Mutex
	IsMouseDragged bool
	Mouse          pixel.Vec
	PMouse         pixel.Vec
	pressed        func(pixel.Button) bool
}

func NewContext(width, height int) *Context {
	var mu sync.Mutex
	return &Context{
		*gg.NewContext(width, height),
		mu,
		false,
		pixel.Vec{X: 0, Y: 0},
		pixel.Vec{X: 0, Y: 0},
		func(pixel.Button) bool { return true },
	}
}

func (ctx *Context) pix() []uint8 {
	return ctx.Image().(*image.RGBA).Pix
}

func (ctx *Context) IsKeyPressed(b pixel.Button) bool {
	return ctx.pressed(b)
}

