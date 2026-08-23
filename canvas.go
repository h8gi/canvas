// Package canvas provides a lightweight, Processing-like 2D animation and creative coding library for Go.
//
// It integrates fogleman/gg for 2D vector drawing with gopxl/pixel for windowing and hardware rendering.
package canvas

import (
	"sync"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

// Canvas represents the drawing area and manages the animation loop and window lifecycle.
type Canvas struct {
	// Width is the width of the canvas window in pixels.
	Width int
	// Height is the height of the canvas window in pixels.
	Height int
	// FrameRate is the target frames per second.
	FrameRate int

	title      string
	resizable  bool
	fullscreen bool
	looping    bool
	redrawReq  bool
	initFunc   func()
	drawFunc   func()
	context    *Context
	mu         sync.Mutex
}

// CanvasConfig holds configuration parameters for creating a new Canvas.
type CanvasConfig struct {
	// Width is the window width in pixels (default: 600).
	Width int
	// Height is the window height in pixels (default: 400).
	Height int
	// FrameRate is the target FPS (default: 60).
	FrameRate int
	// Title is the window title bar text (default: "canvas").
	Title string
	// Resizable allows the window to be resized by the user.
	Resizable bool
	// Fullscreen opens the window in fullscreen on the primary monitor.
	Fullscreen bool
}

// NewCanvas creates and initializes a new Canvas instance with the given configuration options.
// If opts is nil, default options are applied.
func NewCanvas(opts *CanvasConfig) *Canvas {
	width, height, frameRate := 600, 400, 60
	title := "canvas"
	resizable := false
	fullscreen := false

	if opts != nil {
		if opts.Width > 0 {
			width = opts.Width
		}
		if opts.Height > 0 {
			height = opts.Height
		}
		if opts.FrameRate > 0 {
			frameRate = opts.FrameRate
		}
		title = opts.Title
		resizable = opts.Resizable
		fullscreen = opts.Fullscreen
	}

	c := &Canvas{
		Width:      width,
		Height:     height,
		FrameRate:  frameRate,
		title:      title,
		resizable:  resizable,
		fullscreen: fullscreen,
		looping:    true,
		redrawReq:  true,
	}
	c.context = NewContext(width, height)
	// set init drawer
	c.Setup(func(*Context) {})
	return c
}

// NoLoop stops the continuous animation loop.
func (c *Canvas) NoLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.looping = false
}

// Loop resumes the continuous animation loop.
func (c *Canvas) Loop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.looping = true
}

// Redraw requests a single frame draw when the loop is paused.
func (c *Canvas) Redraw() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.redrawReq = true
}

// IsLooping returns true if the animation loop is currently active.
func (c *Canvas) IsLooping() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.looping
}

// Setup registers an initialization function that runs once before the draw loop starts.
func (c *Canvas) Setup(initializer func(*Context)) {
	c.initFunc = func() {
		c.context.mu.Lock()
		initializer(c.context)
		c.context.mu.Unlock()
	}
}

// Draw registers the main rendering callback and starts the window event loop.
// Note: Draw must be called from the main thread / main function.
func (c *Canvas) Draw(drawer func(*Context)) {
	c.drawFunc = func() {
		c.context.mu.Lock()
		drawer(c.context)
		c.context.mu.Unlock()
	}
	c.initFunc()
	opengl.Run(c.startLoop)
}

func (c *Canvas) startLoop() {
	cfg := opengl.WindowConfig{
		Title:     c.title,
		Bounds:    pixel.R(0, 0, float64(c.Width), float64(c.Height)),
		VSync:     true,
		Resizable: c.resizable,
	}
	if c.fullscreen {
		cfg.Monitor = opengl.PrimaryMonitor()
	}

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	c.context.justPressed = win.JustPressed
	c.context.pressed = win.Pressed
	c.context.justReleased = win.JustReleased

	wincan := win.Canvas()
	wincan.SetPixels(c.context.flippedPix())
	win.Update()

	ticker := time.NewTicker(time.Second / time.Duration(c.FrameRate))
	defer ticker.Stop()

	startTime := time.Now()
	lastTime := startTime

	for !win.Closed() {
		now := time.Now()
		c.context.DeltaTime = now.Sub(lastTime).Seconds()
		c.context.Time = now.Sub(startTime).Seconds()
		lastTime = now

		c.context.IsMouseDragged = win.Pressed(pixel.MouseButtonLeft)
		c.context.PMouse = c.context.Mouse
		c.context.Mouse = toCanvasMousePos(win.MousePosition(), c.Height)

		c.mu.Lock()
		shouldDraw := c.looping || c.redrawReq
		c.redrawReq = false
		c.mu.Unlock()

		if shouldDraw {
			c.context.FrameCount++
			c.drawFunc()
			wincan.SetPixels(c.context.flippedPix())
		}

		win.Update()
		<-ticker.C
	}
}

func flipV(src, dst []uint8, width, height int) {
	stride := width * 4
	for y := 0; y < height; y++ {
		srcRow := src[y*stride : (y+1)*stride]
		dstRow := dst[(height-1-y)*stride : (height-y)*stride]
		copy(dstRow, srcRow)
	}
}

func toCanvasMousePos(pixelPos pixel.Vec, height int) pixel.Vec {
	return pixel.V(pixelPos.X, float64(height)-pixelPos.Y)
}

