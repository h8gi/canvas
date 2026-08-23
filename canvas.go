// simple animation library
package canvas

import (
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
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
	width, height, frameRate := 600, 400, 60
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
	opengl.Run(c.startLoop)
}

func (c *Canvas) startLoop() {
	cfg := opengl.WindowConfig{
		Title:  c.title,
		Bounds: pixel.R(0, 0, float64(c.Width), float64(c.Height)),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	c.context.pressed = win.JustPressed
	wincan := win.Canvas()
	wincan.SetPixels(c.context.flippedPix())
	win.Update()

	ticker := time.NewTicker(time.Second / time.Duration(c.FrameRate))
	defer ticker.Stop()

	for !win.Closed() {
		c.context.IsMouseDragged = win.Pressed(pixel.MouseButtonLeft)
		c.context.PMouse = c.context.Mouse
		c.context.Mouse = toCanvasMousePos(win.MousePosition(), c.Height)
		c.drawFunc()
		wincan.SetPixels(c.context.flippedPix())
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

