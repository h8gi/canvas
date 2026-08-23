package canvas

import (
	"image/color"
	"math"
	"testing"
)

func TestContext_ShapePrimitives(t *testing.T) {
	ctx := NewContext(100, 100)

	t.Run("DrawSquare", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(1, 0, 0)
		ctx.DrawSquare(10, 10, 20)
		ctx.Fill()

		pix := ctx.pix()
		// (15, 15) should be red inside square
		idx := (15*100 + 15) * 4
		if pix[idx] != 255 || pix[idx+1] != 0 || pix[idx+2] != 0 {
			t.Errorf("DrawSquare: expected pixel inside square to be red, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}
	})

	t.Run("DrawCenteredRectangle and DrawCenteredSquare", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(1, 1, 0)
		ctx.DrawCenteredRectangle(50, 50, 20, 30)
		ctx.Fill()

		pix := ctx.pix()
		// (50, 50) should be yellow
		idx := (50*100 + 50) * 4
		if pix[idx] != 255 || pix[idx+1] != 255 || pix[idx+2] != 0 {
			t.Errorf("DrawCenteredRectangle: expected center to be yellow, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}

		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(0, 1, 1)
		ctx.DrawCenteredSquare(50, 50, 20)
		ctx.Fill()

		pix = ctx.pix()
		// (50, 50) should be cyan
		idx = (50*100 + 50) * 4
		if pix[idx] != 0 || pix[idx+1] != 255 || pix[idx+2] != 255 {
			t.Errorf("DrawCenteredSquare: expected center to be cyan, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}
	})

	t.Run("RectMode", func(t *testing.T) {
		// Test RectModeCenter
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(1, 0, 0)
		ctx.SetRectMode(RectModeCenter)
		ctx.DrawRect(50, 50, 20, 20)
		ctx.Fill()

		pix := ctx.pix()
		idx := (50*100 + 50) * 4
		if pix[idx] != 255 || pix[idx+1] != 0 || pix[idx+2] != 0 {
			t.Errorf("RectModeCenter: expected center to be red, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}

		// Test RectModeRadius
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(0, 1, 0)
		ctx.SetRectMode(RectModeRadius)
		ctx.DrawRect(50, 50, 10, 10)
		ctx.Fill()

		pix = ctx.pix()
		idx = (50*100 + 50) * 4
		if pix[idx] != 0 || pix[idx+1] != 255 || pix[idx+2] != 0 {
			t.Errorf("RectModeRadius: expected center to be green, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}

		// Reset to default RectModeCorner
		ctx.SetRectMode(RectModeCorner)
	})

	t.Run("EllipseMode", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(1, 0, 1)
		ctx.SetEllipseMode(EllipseModeCenter)
		ctx.DrawEllipseMode(50, 50, 20, 20)
		ctx.Fill()

		pix := ctx.pix()
		idx := (50*100 + 50) * 4
		if pix[idx] != 255 || pix[idx+1] != 0 || pix[idx+2] != 255 {
			t.Errorf("EllipseModeCenter: expected center to be magenta, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}

		ctx.BackgroundRGB(0, 0, 0)
		ctx.SetEllipseMode(EllipseModeCorner)
		ctx.DrawEllipseMode(40, 40, 20, 20)
		ctx.Fill()

		pix = ctx.pix()
		idx = (50*100 + 50) * 4
		if pix[idx] != 255 || pix[idx+1] != 0 || pix[idx+2] != 255 {
			t.Errorf("EllipseModeCorner: expected center (50,50) to be magenta, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}

		ctx.SetEllipseMode(EllipseModeCenter)
	})

	t.Run("DrawTriangle", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(0, 1, 0)
		ctx.DrawTriangle(0, 0, 40, 0, 20, 40)
		ctx.Fill()

		pix := ctx.pix()
		// Centroid around (20, 15) should be green
		idx := (15*100 + 20) * 4
		if pix[idx] != 0 || pix[idx+1] != 255 || pix[idx+2] != 0 {
			t.Errorf("DrawTriangle: expected centroid to be green, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}
	})

	t.Run("DrawQuad", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(0, 0, 1)
		ctx.DrawQuad(10, 10, 50, 10, 40, 40, 0, 40)
		ctx.Fill()

		pix := ctx.pix()
		// (25, 25) should be blue
		idx := (25*100 + 25) * 4
		if pix[idx] != 0 || pix[idx+1] != 0 || pix[idx+2] != 255 {
			t.Errorf("DrawQuad: expected center to be blue, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}
	})

	t.Run("DrawPolygon", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(1, 0, 0)
		ctx.DrawPolygon(50, 50, 20, 6) // Hexagon
		ctx.Fill()

		pix := ctx.pix()
		idx := (50*100 + 50) * 4
		if pix[idx] != 255 || pix[idx+1] != 0 || pix[idx+2] != 0 {
			t.Errorf("DrawPolygon: expected center to be red, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}

		// Invalid sides (< 3) should safely no-op
		ctx.DrawPolygon(50, 50, 20, 2)
	})

	t.Run("DrawStar", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(1, 1, 0)
		ctx.DrawStar(50, 50, 30, 15, 5) // 5-pointed star
		ctx.Fill()

		pix := ctx.pix()
		idx := (50*100 + 50) * 4
		if pix[idx] != 255 || pix[idx+1] != 255 || pix[idx+2] != 0 {
			t.Errorf("DrawStar: expected center to be yellow, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}

		// Invalid points (< 2) should safely no-op
		ctx.DrawStar(50, 50, 30, 15, 1)
	})

	t.Run("DrawSector", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillRGB(0, 0, 1)
		ctx.DrawSector(50, 50, 25, 0, math.Pi) // Half circle
		ctx.Fill()

		pix := ctx.pix()
		// (50, 60) in lower half should be blue
		idx := (60*100 + 50) * 4
		if pix[idx] != 0 || pix[idx+1] != 0 || pix[idx+2] != 255 {
			t.Errorf("DrawSector: expected interior of sector to be blue, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}
	})

	t.Run("BeginShape and EndShape", func(t *testing.T) {
		ctx.BackgroundRGB(0, 0, 0)
		ctx.FillColor(color.White)
		ctx.BeginShape()
		ctx.Vertex(20, 20)
		ctx.Vertex(60, 20)
		ctx.Vertex(60, 60)
		ctx.Vertex(20, 60)
		ctx.EndShape(true)
		ctx.Fill()

		pix := ctx.pix()
		// (40, 40) should be white
		idx := (40*100 + 40) * 4
		if pix[idx] != 255 || pix[idx+1] != 255 || pix[idx+2] != 255 {
			t.Errorf("BeginShape/EndShape: expected center to be white, got [%d %d %d]", pix[idx], pix[idx+1], pix[idx+2])
		}
	})
}
