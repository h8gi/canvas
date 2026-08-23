// draw context
package canvas

import (
	"image"
	"image/color"
	"sync"

	"github.com/fogleman/gg"
)

type Context struct {
	gg.Context
	mu               sync.Mutex
	IsMouseDragged   bool
	Mouse            Vec
	PMouse           Vec
	FrameCount       int
	DeltaTime        float64
	Time             float64
	justPressed      func(Key) bool
	pressed          func(Key) bool
	justReleased     func(Key) bool
	flippedPixBuffer []uint8
}

func NewContext(width, height int) *Context {
	var mu sync.Mutex
	return &Context{
		Context:          *gg.NewContext(width, height),
		mu:               mu,
		IsMouseDragged:   false,
		Mouse:            Vec{X: 0, Y: 0},
		PMouse:           Vec{X: 0, Y: 0},
		justPressed:      func(Key) bool { return false },
		pressed:          func(Key) bool { return false },
		justReleased:     func(Key) bool { return false },
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

// Background clears the canvas with the specified color.
func (ctx *Context) Background(c color.Color) {
	ctx.Push()
	ctx.SetColor(c)
	ctx.Clear()
	ctx.Pop()
}

// BackgroundRGB clears the canvas with an RGB color (r, g, b in 0.0 to 1.0).
func (ctx *Context) BackgroundRGB(r, g, b float64) {
	ctx.Push()
	ctx.SetRGB(r, g, b)
	ctx.Clear()
	ctx.Pop()
}

// BackgroundRGBA clears the canvas with an RGBA color (r, g, b, a in 0.0 to 1.0).
func (ctx *Context) BackgroundRGBA(r, g, b, a float64) {
	ctx.Push()
	ctx.SetRGBA(r, g, b, a)
	ctx.Clear()
	ctx.Pop()
}

// BackgroundHex clears the canvas with a hexadecimal color (e.g., "#ffffff" or "fff").
func (ctx *Context) BackgroundHex(hex string) {
	ctx.Push()
	ctx.SetHexColor(hex)
	ctx.Clear()
	ctx.Pop()
}

// IsKeyPressed returns true if the key was just pressed in the current frame.
func (ctx *Context) IsKeyPressed(k Key) bool {
	return ctx.justPressed(k)
}

// IsKeyDown returns true while the key is being held down.
func (ctx *Context) IsKeyDown(k Key) bool {
	return ctx.pressed(k)
}

// IsKeyJustReleased returns true if the key was released in the current frame.
func (ctx *Context) IsKeyJustReleased(k Key) bool {
	return ctx.justReleased(k)
}

// IsMousePressed returns true while the mouse button is being held down.
func (ctx *Context) IsMousePressed(btn MouseButton) bool {
	return ctx.pressed(btn)
}

// IsMouseJustPressed returns true if the mouse button was clicked in the current frame.
func (ctx *Context) IsMouseJustPressed(btn MouseButton) bool {
	return ctx.justPressed(btn)
}

// IsMouseJustReleased returns true if the mouse button was released in the current frame.
func (ctx *Context) IsMouseJustReleased(btn MouseButton) bool {
	return ctx.justReleased(btn)
}


