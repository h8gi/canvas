// draw context
package canvas

import (
	"image"
	"sync"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"

	"github.com/fogleman/gg"
)

type Context struct {
	gg.Context
	mu sync.Mutex
	// mouse coordinates are rescaled by window size.
	mouseEvents [2]mouse.Event
	keyEvent    key.Event
	dragged     bool
}

func NewContext(width, height int) *Context {
	var mu sync.Mutex
	var mouseEvents [2]mouse.Event
	var k key.Event
	return &Context{
		*gg.NewContext(width, height),
		mu,
		mouseEvents,
		k,
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

func (ctx *Context) IsMouseDragged() bool {
	return ctx.dragged
}

func (ctx *Context) IsKeyPressed() bool {
	return ctx.keyEvent.Direction == key.DirPress
}

func (ctx *Context) IsKeyDown() bool {
	return ctx.keyEvent.Direction == key.DirPress || ctx.keyEvent.Direction == key.DirNone
}

func (ctx *Context) KeyEvent() key.Event {
	return ctx.keyEvent
}

func (ctx *Context) KeyCode() key.Code {
	return ctx.keyEvent.Code
}
