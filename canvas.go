// simple animation library
package canvas

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// drawing area
type Canvas struct {
	Width     int
	Height    int
	FrameRate int
	title     string
	initFunc  func()
	drawFunc  func()
	context   *Context
}

type CanvasConfig struct {
	Width, Height, FrameRate int
	Title                    string
}

func NewCanvas(opts *CanvasConfig) *Canvas {
	width, height, frameRate := 600, 400, 30
	title := "canvas"

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
	}

	c := &Canvas{
		Width:     width,
		Height:    height,
		FrameRate: frameRate,
		title:     title,
	}
	c.context = NewContext(width, height)
	// set init drawer
	c.Setup(func(*Context) {})
	return c
}

// initialize drawer
func (c *Canvas) Setup(initializer func(*Context)) {
	c.initFunc = func() {
		c.context.mu.Lock()
		initializer(c.context)
		c.context.mu.Unlock()
	}
}

// start main loop
func (c *Canvas) Draw(drawer func(*Context)) {
	c.drawFunc = func() {
		c.context.mu.Lock()
		drawer(c.context)
		c.context.mu.Unlock()
	}
	c.initFunc()
	pixelgl.Run(c.startLoop)
}

func (c *Canvas) startLoop() {
	cfg := pixelgl.WindowConfig{
		Title:  c.title,
		Bounds: pixel.R(0, 0, float64(c.Width), float64(c.Height)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	c.context.pressed = win.Pressed
	wincan := win.Canvas()
	wincan.SetPixels(c.context.pix())
	win.Update()

	for !win.Closed() {
		c.context.IsMouseDragged = win.Pressed(pixelgl.MouseButtonLeft)
		c.context.PMouse = c.context.Mouse
		c.context.Mouse = win.MousePosition()
		c.drawFunc()
		wincan.SetPixels(c.context.pix())
		win.Update()
	}
}
