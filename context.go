// draw context
package canvas

import (
	"image"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/fogleman/gg"
)

type Context struct {
	gg.Context
	mu             sync.Mutex
	IsMouseDragged bool
	Mouse          pixel.Vec
	PMouse         pixel.Vec
	pressed        func(pixelgl.Button) bool
}

func NewContext(width, height int) *Context {
	var mu sync.Mutex
	return &Context{
		*gg.NewContext(width, height),
		mu,
		false,
		pixel.Vec{0, 0},
		pixel.Vec{0, 0},
		func(pixelgl.Button) bool { return true },
	}
}

func (ctx *Context) pix() []uint8 {
	return ctx.Image().(*image.RGBA).Pix
}

func (ctx *Context) IsKeyPressed(b pixelgl.Button) bool {
	return ctx.pressed(b)
}
