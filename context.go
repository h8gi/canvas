// draw context
package canvas

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"sync"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/basicfont"
)

// Context embeds gg.Context for 2D vector drawing and provides input states,
// animation time metrics, and convenience helpers.
type Context struct {
	gg.Context

	mu sync.Mutex

	// IsMouseDragged is true when the left mouse button is pressed and held.
	IsMouseDragged bool

	// Mouse holds the current mouse position in canvas coordinates (top-left is (0, 0)).
	Mouse Vec

	// PMouse holds the previous frame's mouse position in canvas coordinates.
	PMouse Vec

	// MouseScroll holds the mouse wheel / trackpad scroll offset for the current frame.
	MouseScroll Vec

	// FrameCount is the total number of frames rendered since canvas startup.
	FrameCount int

	// DeltaTime is the elapsed time in seconds since the previous frame.
	DeltaTime float64

	// Time is the total elapsed time in seconds since canvas startup.
	Time float64

	justPressed      func(Key) bool
	pressed          func(Key) bool
	justReleased     func(Key) bool
	flippedPixBuffer []uint8
	firstVertex      bool
	rectMode         RectMode
	ellipseMode      EllipseMode
}

// NewContext creates a new Context with the specified dimensions and sets the default basic font face.
func NewContext(width, height int) *Context {
	var mu sync.Mutex
	ctx := &Context{
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
	ctx.SetFontFace(basicfont.Face7x13)
	return ctx
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

// BackgroundHSB clears the canvas with an HSB color (h: 0-360, s: 0-1, b: 0-1).
func (ctx *Context) BackgroundHSB(h, s, b float64) {
	ctx.Background(HSB(h, s, b))
}

// BackgroundHSBA clears the canvas with an HSBA color (h: 0-360, s: 0-1, b: 0-1, a: 0-1).
func (ctx *Context) BackgroundHSBA(h, s, b, a float64) {
	ctx.Background(HSBA(h, s, b, a))
}

// FillColor sets the current drawing color.
func (ctx *Context) FillColor(c color.Color) {
	ctx.SetColor(c)
}

// StrokeColor sets the current drawing color.
func (ctx *Context) StrokeColor(c color.Color) {
	ctx.SetColor(c)
}

// FillRGB sets the drawing color with RGB values (0.0 to 1.0).
func (ctx *Context) FillRGB(r, g, b float64) {
	ctx.SetRGB(r, g, b)
}

// StrokeRGB sets the drawing color with RGB values (0.0 to 1.0).
func (ctx *Context) StrokeRGB(r, g, b float64) {
	ctx.SetRGB(r, g, b)
}

// FillRGBA sets the drawing color with RGBA values (0.0 to 1.0).
func (ctx *Context) FillRGBA(r, g, b, a float64) {
	ctx.SetRGBA(r, g, b, a)
}

// StrokeRGBA sets the drawing color with RGBA values (0.0 to 1.0).
func (ctx *Context) StrokeRGBA(r, g, b, a float64) {
	ctx.SetRGBA(r, g, b, a)
}

// FillHex sets the drawing color with a hexadecimal string (e.g., "#ff0000").
func (ctx *Context) FillHex(hex string) {
	ctx.SetHexColor(hex)
}

// StrokeHex sets the drawing color with a hexadecimal string (e.g., "#ff0000").
func (ctx *Context) StrokeHex(hex string) {
	ctx.SetHexColor(hex)
}

// FillHSB sets the drawing color using Hue (0-360), Saturation (0-1), and Brightness (0-1).
func (ctx *Context) FillHSB(h, s, b float64) {
	ctx.SetColor(HSB(h, s, b))
}

// FillHSBA sets the drawing color using Hue (0-360), Saturation (0-1), Brightness (0-1), and Alpha (0-1).
func (ctx *Context) FillHSBA(h, s, b, a float64) {
	ctx.SetColor(HSBA(h, s, b, a))
}

// StrokeHSB sets the drawing color using Hue (0-360), Saturation (0-1), and Brightness (0-1).
func (ctx *Context) StrokeHSB(h, s, b float64) {
	ctx.SetColor(HSB(h, s, b))
}

// StrokeHSBA sets the drawing color using Hue (0-360), Saturation (0-1), Brightness (0-1), and Alpha (0-1).
func (ctx *Context) StrokeHSBA(h, s, b, a float64) {
	ctx.SetColor(HSBA(h, s, b, a))
}

// NoFill sets the fill color to transparent.
func (ctx *Context) NoFill() {
	ctx.SetColor(color.Transparent)
}

// NoStroke sets the stroke color to transparent.
func (ctx *Context) NoStroke() {
	ctx.SetColor(color.Transparent)
}

// SaveFrame saves the current context buffer as a PNG image to the given path.
func (ctx *Context) SaveFrame(path string) error {
	return ctx.SavePNG(path)
}

// SaveFrameSeq saves the current context buffer as a formatted sequence frame.
// pattern should contain a printf integer format verb (e.g. "frames/frame_%04d.png").
// Parent directories are automatically created if they do not exist.
func (ctx *Context) SaveFrameSeq(pattern string) error {
	filePath := fmt.Sprintf(pattern, ctx.FrameCount)
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return ctx.SaveFrame(filePath)
}

// CloneImage returns a deep copy of the current context RGBA image.
func (ctx *Context) CloneImage() *image.RGBA {
	src := ctx.Image().(*image.RGBA)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	copy(dst.Pix, src.Pix)
	return dst
}

// GetPixel returns the color of the pixel at (x, y) in canvas coordinates.
// If (x, y) is out of bounds, a transparent color is returned.
func (ctx *Context) GetPixel(x, y int) color.RGBA {
	if x < 0 || x >= ctx.Width() || y < 0 || y >= ctx.Height() {
		return color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}
	pix := ctx.pix()
	idx := (y*ctx.Width() + x) * 4
	return color.RGBA{
		R: pix[idx],
		G: pix[idx+1],
		B: pix[idx+2],
		A: pix[idx+3],
	}
}

// SetPixelColor sets the color of the pixel at (x, y) in canvas coordinates.
func (ctx *Context) SetPixelColor(x, y int, c color.Color) {
	if x < 0 || x >= ctx.Width() || y < 0 || y >= ctx.Height() {
		return
	}
	r, g, b, a := c.RGBA()
	pix := ctx.pix()
	idx := (y*ctx.Width() + x) * 4
	pix[idx] = uint8(r >> 8)
	pix[idx+1] = uint8(g >> 8)
	pix[idx+2] = uint8(b >> 8)
	pix[idx+3] = uint8(a >> 8)
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



