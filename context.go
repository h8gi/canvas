// draw context
package canvas

import (
	"image"
	"sync"

	"golang.org/x/mobile/event/mouse"

	"github.com/fogleman/gg"
)

type Context struct {
	gg.Context
	mu sync.Mutex
	// mouse coordinates are rescaled by window size.
	mouseEvents [2]mouse.Event
	pressed     bool
}

func NewContext(width, height int) *Context {
	var mu sync.Mutex
	var mouseEvents [2]mouse.Event
	return &Context{
		*gg.NewContext(width, height),
		mu,
		mouseEvents,
		false,
	}
}

func (ctx *Context) pix() []uint8 {
	return ctx.Image().(*image.RGBA).Pix
}

func (ctx *Context) pushMouseEvent(m mouse.Event) {
	ctx.mouseEvents = [2]mouse.Event{ctx.mouseEvents[1], m}
}

func (ctx *Context) PreviousMouseEvent() mouse.Event {
	return ctx.mouseEvents[0]
}

func (ctx *Context) MouseEvent() mouse.Event {
	return ctx.mouseEvents[1]
}

func (ctx *Context) MouseX() float64 {
	return float64(ctx.MouseEvent().X)
}

func (ctx *Context) MouseY() float64 {
	return float64(ctx.MouseEvent().Y)
}

func (ctx *Context) PreviousMouseX() float64 {
	return float64(ctx.PreviousMouseEvent().X)
}

func (ctx *Context) PreviousMouseY() float64 {
	return float64(ctx.PreviousMouseEvent().Y)
}

func (ctx *Context) MousePressed() bool {
	return ctx.pressed
}
