package canvas

import (
	"image"
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// implement fyne/canvasObject
type FyneCanvas struct {
	size      fyne.Size
	position  fyne.Position
	hidden    bool
	frameRate int
	title     string
	initFunc  func()
	drawFunc  func()
	context   *Context
}

func (c *FyneCanvas) Size() fyne.Size {
	return c.size
}

func (c *FyneCanvas) Resize(size fyne.Size) {
	c.size = size
	widget.Renderer(c).Layout(size)
}

func (c *FyneCanvas) Position() fyne.Position {
	return c.position
}

func (c *FyneCanvas) Move(pos fyne.Position) {
	c.position = pos
	widget.Renderer(c).Layout(c.size)
}

func (c *FyneCanvas) MinSize() fyne.Size {
	return widget.Renderer(c).MinSize()
}

func (c *FyneCanvas) Visible() bool {
	return c.hidden
}

func (c *FyneCanvas) Show() {
	c.hidden = false
}

func (c *FyneCanvas) Hide() {
	c.hidden = true
}

type fyneCanvasRenderer struct {
	render   *canvas.Raster
	objects  []fyne.CanvasObject
	imgCache *image.RGBA

	canvas *FyneCanvas
}

func (c *fyneCanvasRenderer) MinSize() fyne.Size {
	return fyne.NewSize(c.canvas.context.Width(), c.canvas.context.Height())
}

func (c *fyneCanvasRenderer) Layout(size fyne.Size) {
	c.render.Resize(size)
}

func (c *fyneCanvasRenderer) ApplyTheme() {}

func (c *fyneCanvasRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (c *fyneCanvasRenderer) Refresh() {
	canvas.Refresh(c.render)
}

func (c *fyneCanvasRenderer) Objects() []fyne.CanvasObject {
	return c.objects
}

func (c *fyneCanvasRenderer) Destroy() {
}

func (c *fyneCanvasRenderer) draw(w, h int) image.Image {
	c.imgCache = c.canvas.context.Image().(*image.RGBA)

	return c.imgCache
}

func (c *FyneCanvas) CreateRenderer() fyne.WidgetRenderer {
	renderer := &fyneCanvasRenderer{canvas: c}

	render := canvas.NewRaster(renderer.draw)
	renderer.render = render
	renderer.objects = []fyne.CanvasObject{render}
	renderer.ApplyTheme()

	return renderer
}

func (c *FyneCanvas) animate() {
	go func() {
		tick := time.NewTicker(time.Second / 6)

		for {
			select {
			case <-tick.C:
				c.drawFunc()
				widget.Refresh(c)
			}
		}
	}()
}

func (c *FyneCanvas) Tapped(ev *fyne.PointEvent) {

}

func (c *FyneCanvas) Setup(initializer func(*Context)) {
	c.initFunc = func() {
		c.context.mu.Lock()
		initializer(c.context)
		c.context.mu.Unlock()
	}
}

func (c *FyneCanvas) Draw(drawer func(*Context)) {
	c.drawFunc = func() {
		c.context.mu.Lock()
		drawer(c.context)
		c.context.mu.Unlock()
	}
	c.initFunc()
	c.startLoop()
}

func FyneNew(opts *NewCanvasOptions) *FyneCanvas {
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

	c := &FyneCanvas{
		frameRate: frameRate,
		title:     title,
	}
	c.context = NewContext(width, height)
	return c
}

func (c *FyneCanvas) startLoop() {
	a := app.New()

	window := a.NewWindow(c.title)
	window.SetContent(c)

	c.animate()

	window.Show()
	a.Run()
}
