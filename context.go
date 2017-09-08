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
}

func NewContext(width, height int) *Context {
	var mu sync.Mutex
	var mouseEvents [2]mouse.Event
	return &Context{
		*gg.NewContext(width, height),
		mu,
		mouseEvents,
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
