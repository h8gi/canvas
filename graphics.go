package canvas

import "math"

// CreateGraphics creates an offscreen Context with the specified width and height.
// It can be drawn onto another Context using DrawGraphics or DrawGraphicsAnchored.
func CreateGraphics(width, height int) *Context {
	return NewContext(width, height)
}

// NewGraphics is an alias for CreateGraphics.
func NewGraphics(width, height int) *Context {
	return CreateGraphics(width, height)
}

// CreateGraphics creates an offscreen Context with the specified dimensions.
func (c *Canvas) CreateGraphics(width, height int) *Context {
	return CreateGraphics(width, height)
}

// DrawGraphics renders the offscreen context g onto ctx at (x, y).
func (ctx *Context) DrawGraphics(g *Context, x, y float64) {
	if g == nil {
		return
	}
	ctx.DrawImage(g.Image(), int(math.Round(x)), int(math.Round(y)))
}

// DrawGraphicsAnchored renders the offscreen context g onto ctx with anchor alignments (ax, ay).
// For example, ax = 0.5, ay = 0.5 centers the graphic at (x, y).
func (ctx *Context) DrawGraphicsAnchored(g *Context, x, y, ax, ay float64) {
	if g == nil {
		return
	}
	ctx.DrawImageAnchored(g.Image(), int(math.Round(x)), int(math.Round(y)), ax, ay)
}
