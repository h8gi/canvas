package canvas

import (
	"image/color"
	"testing"

	"github.com/gopxl/pixel/v2"
)

func TestNewCanvas(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		c := NewCanvas(nil)
		if c.Width != 600 {
			t.Errorf("expected Width 600, got %d", c.Width)
		}
		if c.Height != 400 {
			t.Errorf("expected Height 400, got %d", c.Height)
		}
		if c.FrameRate != 60 {
			t.Errorf("expected FrameRate 60, got %d", c.FrameRate)
		}
		if c.title != "canvas" {
			t.Errorf("expected title 'canvas', got '%s'", c.title)
		}
		if c.context == nil {
			t.Fatal("expected context to be initialized")
		}
	})

	t.Run("custom options", func(t *testing.T) {
		cfg := &CanvasConfig{
			Width:     800,
			Height:    600,
			FrameRate: 30,
			Title:     "Custom Title",
		}
		c := NewCanvas(cfg)
		if c.Width != 800 {
			t.Errorf("expected Width 800, got %d", c.Width)
		}
		if c.Height != 600 {
			t.Errorf("expected Height 600, got %d", c.Height)
		}
		if c.FrameRate != 30 {
			t.Errorf("expected FrameRate 30, got %d", c.FrameRate)
		}
		if c.title != "Custom Title" {
			t.Errorf("expected title 'Custom Title', got '%s'", c.title)
		}
	})
}

func TestCanvas_Setup(t *testing.T) {
	c := NewCanvas(nil)
	called := false
	c.Setup(func(ctx *Context) {
		called = true
		ctx.SetColor(color.White)
	})

	if c.initFunc == nil {
		t.Fatal("expected initFunc to be non-nil")
	}
	c.initFunc()
	if !called {
		t.Error("expected initializer function to be called")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(100, 200)
	if ctx == nil {
		t.Fatal("expected ctx to be non-nil")
	}
	if ctx.Width() != 100 {
		t.Errorf("expected width 100, got %d", ctx.Width())
	}
	if ctx.Height() != 200 {
		t.Errorf("expected height 200, got %d", ctx.Height())
	}
	if ctx.IsMouseDragged {
		t.Error("expected IsMouseDragged to be false initially")
	}
	if ctx.Mouse.X != 0 || ctx.Mouse.Y != 0 {
		t.Errorf("expected Mouse to be (0,0), got (%v,%v)", ctx.Mouse.X, ctx.Mouse.Y)
	}
	if ctx.PMouse.X != 0 || ctx.PMouse.Y != 0 {
		t.Errorf("expected PMouse to be (0,0), got (%v,%v)", ctx.PMouse.X, ctx.PMouse.Y)
	}
	if !ctx.IsKeyPressed(pixel.KeySpace) {
		t.Error("expected default pressed func to return true")
	}

	// Test pix()
	pix := ctx.pix()
	expectedLen := 100 * 200 * 4
	if len(pix) != expectedLen {
		t.Errorf("expected pix len %d, got %d", expectedLen, len(pix))
	}
}

func TestContext_Drawing(t *testing.T) {
	ctx := NewContext(10, 10)
	ctx.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	ctx.Clear()

	pix := ctx.pix()
	// First pixel should be Red (255, 0, 0, 255)
	if pix[0] != 255 || pix[1] != 0 || pix[2] != 0 || pix[3] != 255 {
		t.Errorf("expected pixel to be red [255 0 0 255], got [%d %d %d %d]", pix[0], pix[1], pix[2], pix[3])
	}
}
