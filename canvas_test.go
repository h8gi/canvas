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
	if ctx.IsKeyPressed(pixel.KeySpace) {
		t.Error("expected default IsKeyPressed to return false")
	}
	if ctx.IsKeyDown(pixel.KeySpace) {
		t.Error("expected default IsKeyDown to return false")
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

func TestFlipV(t *testing.T) {
	// 2x2 image (4 pixels = 16 bytes)
	// Row 0: Pixel 0: (1, 1, 1, 1), Pixel 1: (2, 2, 2, 2)
	// Row 1: Pixel 2: (3, 3, 3, 3), Pixel 3: (4, 4, 4, 4)
	src := []uint8{
		1, 1, 1, 1, 2, 2, 2, 2,
		3, 3, 3, 3, 4, 4, 4, 4,
	}
	dst := make([]uint8, len(src))
	flipV(src, dst, 2, 2)

	// After flip, Row 0 in dst should be Row 1 of src, and Row 1 in dst should be Row 0 of src
	expected := []uint8{
		3, 3, 3, 3, 4, 4, 4, 4,
		1, 1, 1, 1, 2, 2, 2, 2,
	}
	for i, b := range dst {
		if b != expected[i] {
			t.Fatalf("at index %d: expected %d, got %d", i, expected[i], b)
		}
	}
}

func TestContext_FlippedPix(t *testing.T) {
	// Draw a dot at top-left (0, 0)
	ctx := NewContext(10, 10)
	ctx.SetColor(color.Black)
	ctx.Clear()
	// Set (0,0) to Red
	ctx.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	ctx.SetPixel(0, 0)

	// ctx.pix() (top-left origin) should have Red at pixel (0,0) which is index 0
	pix := ctx.pix()
	if pix[0] != 255 || pix[1] != 0 || pix[2] != 0 {
		t.Fatalf("expected (0,0) in pix to be red, got [%d %d %d]", pix[0], pix[1], pix[2])
	}

	// In OpenGL texture space (bottom-left origin), the top-left pixel (0,0) in gg
	// should end up at the top row in OpenGL, which is at the end of the flipped buffer:
	// y = height - 1 = 9, x = 0 -> index = (9 * 10 + 0) * 4 = 360
	flipped := ctx.flippedPix()
	if flipped[360] != 255 || flipped[361] != 0 || flipped[362] != 0 {
		t.Fatalf("expected top row in flippedPix to contain red at index 360, got [%d %d %d]", flipped[360], flipped[361], flipped[362])
	}
	// And bottom row in flippedPix (index 0) should be Black
	if flipped[0] != 0 || flipped[1] != 0 || flipped[2] != 0 {
		t.Fatalf("expected index 0 in flippedPix to be black, got [%d %d %d]", flipped[0], flipped[1], flipped[2])
	}
}

func TestToCanvasMousePos(t *testing.T) {
	height := 400

	// OpenGL (0, 0) is bottom-left -> Processing canvas (0, height)
	bottomLeft := toCanvasMousePos(pixel.V(0, 0), height)
	if bottomLeft.X != 0 || bottomLeft.Y != 400 {
		t.Errorf("expected bottom-left to be (0, 400), got (%v, %v)", bottomLeft.X, bottomLeft.Y)
	}

	// OpenGL (0, 400) is top-left -> Processing canvas (0, 0)
	topLeft := toCanvasMousePos(pixel.V(0, 400), height)
	if topLeft.X != 0 || topLeft.Y != 0 {
		t.Errorf("expected top-left to be (0, 0), got (%v, %v)", topLeft.X, topLeft.Y)
	}

	// OpenGL center (300, 200) -> Processing canvas (300, 200)
	center := toCanvasMousePos(pixel.V(300, 200), height)
	if center.X != 300 || center.Y != 200 {
		t.Errorf("expected center to be (300, 200), got (%v, %v)", center.X, center.Y)
	}
}

func TestContext_Background(t *testing.T) {
	ctx := NewContext(10, 10)
	ctx.BackgroundRGB(1, 0, 0)

	pix := ctx.pix()
	if pix[0] != 255 || pix[1] != 0 || pix[2] != 0 || pix[3] != 255 {
		t.Errorf("expected BackgroundRGB to fill with red, got [%d %d %d %d]", pix[0], pix[1], pix[2], pix[3])
	}

	ctx.Background(color.RGBA{G: 255, A: 255})
	pix = ctx.pix()
	if pix[0] != 0 || pix[1] != 255 || pix[2] != 0 || pix[3] != 255 {
		t.Errorf("expected Background to fill with green, got [%d %d %d %d]", pix[0], pix[1], pix[2], pix[3])
	}

	ctx.BackgroundHex("#0000ff")
	pix = ctx.pix()
	if pix[0] != 0 || pix[1] != 0 || pix[2] != 255 || pix[3] != 255 {
		t.Errorf("expected BackgroundHex to fill with blue, got [%d %d %d %d]", pix[0], pix[1], pix[2], pix[3])
	}
}

func TestContext_InputAndFrameState(t *testing.T) {
	ctx := NewContext(10, 10)

	ctx.justPressed = func(b Key) bool { return b == KeySpace }
	ctx.pressed = func(b Key) bool { return b == KeyA || b == MouseLeft }
	ctx.justReleased = func(b Key) bool { return b == KeyEscape }

	if !ctx.IsKeyPressed(KeySpace) {
		t.Error("expected KeySpace to be just pressed")
	}
	if ctx.IsKeyPressed(KeyA) {
		t.Error("expected KeyA not to be just pressed")
	}
	if !ctx.IsKeyDown(KeyA) {
		t.Error("expected KeyA to be down")
	}
	if !ctx.IsKeyJustReleased(KeyEscape) {
		t.Error("expected KeyEscape to be just released")
	}

	if !ctx.IsMousePressed(MouseLeft) {
		t.Error("expected MouseLeft to be pressed")
	}
	if ctx.IsMouseJustPressed(MouseLeft) {
		t.Error("expected MouseLeft not to be just pressed")
	}
	if ctx.IsMouseJustReleased(MouseLeft) {
		t.Error("expected MouseLeft not to be just released")
	}
}

func TestContext_DefaultFont(t *testing.T) {
	ctx := NewContext(100, 100)
	// Should not panic when drawing text without explicitly setting a font face
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("drawing text with default font panicked: %v", r)
		}
	}()
	ctx.SetColor(color.White)
	ctx.DrawString("Hello", 10, 20)
}

func TestContext_StyleHelpers(t *testing.T) {
	ctx := NewContext(50, 50)

	// FillHex
	ctx.FillHex("#ff0000")
	ctx.DrawRectangle(0, 0, 50, 50)
	ctx.Fill()
	pix := ctx.pix()
	if pix[0] != 255 || pix[1] != 0 || pix[2] != 0 {
		t.Errorf("expected FillHex to set fill color to red, got [%d %d %d]", pix[0], pix[1], pix[2])
	}

	// StrokeRGB & FillRGB
	ctx.FillRGB(0, 1, 0)
	ctx.StrokeRGB(0, 0, 1)

	// FillRGBA & StrokeRGBA
	ctx.FillRGBA(1, 1, 0, 0.5)
	ctx.StrokeRGBA(0, 1, 1, 0.8)

	// NoFill & NoStroke (transparent)
	ctx.NoFill()
	ctx.NoStroke()
}

func TestContext_SaveFrame(t *testing.T) {
	ctx := NewContext(10, 10)
	ctx.BackgroundRGB(1, 0, 0)

	tmpFile := t.TempDir() + "/test_frame.png"
	err := ctx.SaveFrame(tmpFile)
	if err != nil {
		t.Fatalf("SaveFrame failed: %v", err)
	}
}

func TestCanvas_LoopControl(t *testing.T) {
	c := NewCanvas(nil)
	if !c.IsLooping() {
		t.Error("expected canvas to be looping by default")
	}

	c.NoLoop()
	if c.IsLooping() {
		t.Error("expected canvas not to be looping after NoLoop()")
	}

	c.Loop()
	if !c.IsLooping() {
		t.Error("expected canvas to be looping after Loop()")
	}

	c.NoLoop()
	c.Redraw()
	if !c.redrawReq {
		t.Error("expected redrawReq to be true after Redraw()")
	}
}

func TestCanvas_ConfigOptions(t *testing.T) {
	cfg := &CanvasConfig{
		Width:      800,
		Height:     600,
		FrameRate:  60,
		Title:      "Extended Config",
		Resizable:  true,
		Fullscreen: true,
	}
	c := NewCanvas(cfg)
	if !c.resizable {
		t.Error("expected resizable to be true")
	}
	if !c.fullscreen {
		t.Error("expected fullscreen to be true")
	}
}




