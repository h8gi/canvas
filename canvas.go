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

	title             string
	resizable         bool
	fullscreen        bool
	looping           bool
	redrawReq         bool
	initFunc          func()
	drawFunc          func()
	keyPressedFunc    func(*Context, Key)
	keyReleasedFunc   func(*Context, Key)
	mousePressedFunc  func(*Context, MouseButton)
	mouseReleasedFunc func(*Context, MouseButton)
	mouseMovedFunc    func(*Context, Vec)
	context           *Context
	recorder          *gifRecorder
	mu                sync.Mutex
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
		recorder:   newGIFRecorder(frameRate),
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

// RecordGIF starts recording the specified number of rendered frames into an animated GIF file at path.
// If opts is nil, default recording options are used.
func (c *Canvas) RecordGIF(path string, frames int, opts ...*GIFOptions) {
	var opt *GIFOptions
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	c.recorder.start(path, frames, opt)
}

// IsRecordingGIF returns true if a GIF recording is currently in progress.
func (c *Canvas) IsRecordingGIF() bool {
	return c.recorder.isRecording()
}

// KeyPressed registers a callback invoked when a keyboard key is pressed.
func (c *Canvas) KeyPressed(fn func(*Context, Key)) {
	c.keyPressedFunc = fn
}

// OnKeyPressed is an alias for KeyPressed.
func (c *Canvas) OnKeyPressed(fn func(*Context, Key)) {
	c.KeyPressed(fn)
}

// KeyReleased registers a callback invoked when a keyboard key is released.
func (c *Canvas) KeyReleased(fn func(*Context, Key)) {
	c.keyReleasedFunc = fn
}

// OnKeyReleased is an alias for KeyReleased.
func (c *Canvas) OnKeyReleased(fn func(*Context, Key)) {
	c.KeyReleased(fn)
}

// MousePressed registers a callback invoked when a mouse button is clicked.
func (c *Canvas) MousePressed(fn func(*Context, MouseButton)) {
	c.mousePressedFunc = fn
}

// OnMousePressed is an alias for MousePressed.
func (c *Canvas) OnMousePressed(fn func(*Context, MouseButton)) {
	c.MousePressed(fn)
}

// MouseReleased registers a callback invoked when a mouse button is released.
func (c *Canvas) MouseReleased(fn func(*Context, MouseButton)) {
	c.mouseReleasedFunc = fn
}

// OnMouseReleased is an alias for MouseReleased.
func (c *Canvas) OnMouseReleased(fn func(*Context, MouseButton)) {
	c.MouseReleased(fn)
}

// MouseMoved registers a callback invoked when the mouse cursor position changes.
func (c *Canvas) MouseMoved(fn func(*Context, Vec)) {
	c.mouseMovedFunc = fn
}

// OnMouseMoved is an alias for MouseMoved.
func (c *Canvas) OnMouseMoved(fn func(*Context, Vec)) {
	c.MouseMoved(fn)
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

		// Dispatch key events
		if c.keyPressedFunc != nil {
			for _, k := range allKeys {
				if win.JustPressed(k) {
					c.keyPressedFunc(c.context, k)
				}
			}
		}
		if c.keyReleasedFunc != nil {
			for _, k := range allKeys {
				if win.JustReleased(k) {
					c.keyReleasedFunc(c.context, k)
				}
			}
		}

		// Dispatch mouse events
		if c.mousePressedFunc != nil {
			for _, btn := range allMouseButtons {
				if win.JustPressed(btn) {
					c.mousePressedFunc(c.context, btn)
				}
			}
		}
		if c.mouseReleasedFunc != nil {
			for _, btn := range allMouseButtons {
				if win.JustReleased(btn) {
					c.mouseReleasedFunc(c.context, btn)
				}
			}
		}
		if c.mouseMovedFunc != nil && c.context.Mouse != c.context.PMouse {
			c.mouseMovedFunc(c.context, c.context.Mouse)
		}

		c.mu.Lock()
		shouldDraw := c.looping || c.redrawReq
		c.redrawReq = false
		c.mu.Unlock()

		if shouldDraw {
			c.context.FrameCount++
			c.drawFunc()
			wincan.SetPixels(c.context.flippedPix())
			c.recorder.capture(c.context)
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

