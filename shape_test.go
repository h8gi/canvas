package canvas

import (
	"image/color"
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
