package canvas

import "math"

// RectMode represents how rectangle coordinates are interpreted.
type RectMode int

const (
	// RectModeCorner interprets (x, y) as the top-left corner (default).
	RectModeCorner RectMode = iota
	// RectModeCenter interprets (x, y) as the center point.
	RectModeCenter
	// RectModeRadius interprets (x, y) as the center point, and (w, h) as half-width and half-height.
	RectModeRadius
)

// EllipseMode represents how ellipse / circle coordinates are interpreted.
type EllipseMode int

const (
	// EllipseModeCenter interprets (x, y) as the center point (default).
	EllipseModeCenter EllipseMode = iota
	// EllipseModeCorner interprets (x, y) as the top-left corner of the bounding box.
	EllipseModeCorner
	// EllipseModeRadius interprets (x, y) as the center point, and (w, h) as horizontal and vertical radii.
	EllipseModeRadius
)

// SetRectMode sets the current rectangle drawing mode (RectModeCorner, RectModeCenter, or RectModeRadius).
func (ctx *Context) SetRectMode(mode RectMode) {
	ctx.rectMode = mode
}

// SetEllipseMode sets the current ellipse drawing mode (EllipseModeCenter, EllipseModeCorner, or EllipseModeRadius).
func (ctx *Context) SetEllipseMode(mode EllipseMode) {
	ctx.ellipseMode = mode
}

// DrawSquare draws a square with top-left corner at (x, y) and side length s.
func (ctx *Context) DrawSquare(x, y, s float64) {
	ctx.DrawRectangle(x, y, s, s)
}

// DrawCenteredRectangle draws a rectangle centered at (x, y) with width w and height h.
func (ctx *Context) DrawCenteredRectangle(x, y, w, h float64) {
	ctx.DrawRectangle(x-w/2.0, y-h/2.0, w, h)
}

// DrawCenteredSquare draws a square centered at (x, y) with side length s.
func (ctx *Context) DrawCenteredSquare(x, y, s float64) {
	ctx.DrawRectangle(x-s/2.0, y-s/2.0, s, s)
}

// DrawRect draws a rectangle using the current RectMode.
func (ctx *Context) DrawRect(x, y, w, h float64) {
	switch ctx.rectMode {
	case RectModeCenter:
		ctx.DrawRectangle(x-w/2.0, y-h/2.0, w, h)
	case RectModeRadius:
		ctx.DrawRectangle(x-w, y-h, w*2.0, h*2.0)
	default: // RectModeCorner
		ctx.DrawRectangle(x, y, w, h)
	}
}

// DrawEllipseMode draws an ellipse using the current EllipseMode.
// In default EllipseModeCenter, (w, h) are the full width and height (diameters).
func (ctx *Context) DrawEllipseMode(x, y, w, h float64) {
	switch ctx.ellipseMode {
	case EllipseModeCorner:
		ctx.DrawEllipse(x+w/2.0, y+h/2.0, w/2.0, h/2.0)
	case EllipseModeRadius:
		ctx.DrawEllipse(x, y, w, h)
	default: // EllipseModeCenter
		ctx.DrawEllipse(x, y, w/2.0, h/2.0)
	}
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

// DrawPolygon draws a regular polygon centered at (x, y) with the specified outer radius and number of sides.
func (ctx *Context) DrawPolygon(x, y, radius float64, sides int) {
	if sides < 3 || radius <= 0 {
		return
	}

	angleStep := (math.Pi * 2.0) / float64(sides)
	startAngle := -math.Pi / 2.0 // Top vertex first

	for i := 0; i < sides; i++ {
		angle := startAngle + float64(i)*angleStep
		px := x + radius*math.Cos(angle)
		py := y + radius*math.Sin(angle)
		if i == 0 {
			ctx.MoveTo(px, py)
		} else {
			ctx.LineTo(px, py)
		}
	}
	ctx.ClosePath()
}

// DrawStar draws an n-pointed star centered at (x, y) with outer radius radius1 and inner radius radius2.
func (ctx *Context) DrawStar(x, y, radius1, radius2 float64, points int) {
	if points < 2 || radius1 <= 0 || radius2 <= 0 {
		return
	}

	totalVertices := points * 2
	angleStep := math.Pi / float64(points)
	startAngle := -math.Pi / 2.0 // Top vertex first

	for i := 0; i < totalVertices; i++ {
		r := radius1
		if i%2 != 0 {
			r = radius2
		}
		angle := startAngle + float64(i)*angleStep
		px := x + r*math.Cos(angle)
		py := y + r*math.Sin(angle)
		if i == 0 {
			ctx.MoveTo(px, py)
		} else {
			ctx.LineTo(px, py)
		}
	}
	ctx.ClosePath()
}

// DrawSector draws a pie/wedge sector centered at (x, y) with radius r between startAngle and endAngle (in radians).
func (ctx *Context) DrawSector(x, y, r, startAngle, endAngle float64) {
	if r <= 0 {
		return
	}
	ctx.MoveTo(x, y)
	ctx.DrawEllipticalArc(x, y, r, r, startAngle, endAngle)
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
