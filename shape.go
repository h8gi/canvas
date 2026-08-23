package canvas

// DrawSquare draws a square with top-left corner at (x, y) and side length s.
func (ctx *Context) DrawSquare(x, y, s float64) {
	ctx.DrawRectangle(x, y, s, s)
}

// DrawTriangle draws a triangle with vertices at (x1, y1), (x2, y2), and (x3, y3).
func (ctx *Context) DrawTriangle(x1, y1, x2, y2, x3, y3 float64) {
	ctx.MoveTo(x1, y1)
	ctx.LineTo(x2, y2)
	ctx.LineTo(x3, y3)
	ctx.ClosePath()
}

// DrawQuad draws a quadrilateral with vertices at (x1, y1), (x2, y2), (x3, y3), and (x4, y4).
func (ctx *Context) DrawQuad(x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	ctx.MoveTo(x1, y1)
	ctx.LineTo(x2, y2)
	ctx.LineTo(x3, y3)
	ctx.LineTo(x4, y4)
	ctx.ClosePath()
}

// BeginShape starts building a custom path using Vertex calls.
func (ctx *Context) BeginShape() {
	ctx.ClearPath()
	ctx.firstVertex = true
}

// Vertex adds a coordinate point to the current custom shape path.
func (ctx *Context) Vertex(x, y float64) {
	if ctx.firstVertex {
		ctx.MoveTo(x, y)
		ctx.firstVertex = false
	} else {
		ctx.LineTo(x, y)
	}
}

// EndShape completes the custom shape path. If close is true, the path is closed.
func (ctx *Context) EndShape(close bool) {
	if close {
		ctx.ClosePath()
	}
}
