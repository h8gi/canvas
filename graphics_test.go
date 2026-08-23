package canvas

import (
	"testing"
)

func TestOffscreenGraphics(t *testing.T) {
	t.Run("CreateGraphics and DrawGraphics", func(t *testing.T) {
		mainCtx := NewContext(100, 100)
		mainCtx.BackgroundRGB(0, 0, 0)

		// Create offscreen layer
		layer := CreateGraphics(40, 40)
		layer.BackgroundRGB(1, 0, 0) // Red square

		// Draw layer onto main context at (20, 20)
		mainCtx.DrawGraphics(layer, 20, 20)

		// (10, 10) should be black background
		pOutside := mainCtx.GetPixel(10, 10)
		if pOutside.R != 0 || pOutside.G != 0 || pOutside.B != 0 {
			t.Errorf("expected background to be black, got %+v", pOutside)
		}

		// (30, 30) should be red from layer
		pInside := mainCtx.GetPixel(30, 30)
		if pInside.R != 255 || pInside.G != 0 || pInside.B != 0 {
			t.Errorf("expected layer area to be red, got %+v", pInside)
		}

		// Nil graphic should safely no-op
		mainCtx.DrawGraphics(nil, 0, 0)
	})

	t.Run("DrawGraphicsAnchored", func(t *testing.T) {
		mainCtx := NewContext(100, 100)
		mainCtx.BackgroundRGB(0, 0, 0)

		layer := NewGraphics(20, 20)
		layer.BackgroundRGB(0, 1, 0) // Green square

		// Center the 20x20 layer at (50, 50) -> bounds (40..60, 40..60)
		mainCtx.DrawGraphicsAnchored(layer, 50, 50, 0.5, 0.5)

		pCenter := mainCtx.GetPixel(50, 50)
		if pCenter.R != 0 || pCenter.G != 255 || pCenter.B != 0 {
			t.Errorf("expected anchored center to be green, got %+v", pCenter)
		}

		// Nil graphic should safely no-op
		mainCtx.DrawGraphicsAnchored(nil, 50, 50, 0.5, 0.5)
	})

	t.Run("Canvas CreateGraphics", func(t *testing.T) {
		c := NewCanvas(nil)
		g := c.CreateGraphics(50, 50)
		if g == nil || g.Width() != 50 || g.Height() != 50 {
			t.Errorf("Canvas.CreateGraphics failed")
		}
	})
}
