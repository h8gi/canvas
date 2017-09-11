// simple animation library
package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

// timer event
type tickEvent struct{}

// drawing area
type Canvas struct {
	width     int
	height    int
	frameRate int
	title     string
	initFunc  func()
	drawFunc  func()
	context   *Context
}

type NewCanvasOptions struct {
	Width, Height, FrameRate int
	Title                    string
}

func New(opts *NewCanvasOptions) *Canvas {
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
		width:     width,
		height:    height,
		frameRate: frameRate,
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
	c.startLoop()
}

// start inner simulation loop
func (c *Canvas) simulate(q screen.EventDeque) {
	duration := time.Second / time.Duration(c.frameRate)
	for {
		// memory lock
		c.drawFunc()
		q.Send(tickEvent{})
		time.Sleep(duration)
	}
}

func (c *Canvas) startLoop() {
	rand.Seed(time.Now().UnixNano())

	driver.Main(func(s screen.Screen) {
		// create window
		bufSize := image.Point{c.width, c.height}
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  c.width,
			Height: c.height,
			Title:  c.title,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		// create image buffer
		b, err := s.NewBuffer(bufSize)
		if err != nil {
			log.Fatal(err)
		}
		defer b.Release()

		tex, err := s.NewTexture(bufSize)
		if err != nil {
			log.Fatal(err)
		}
		defer tex.Release()
		tex.Fill(tex.Bounds(), color.White, draw.Src)

		// invoke timer event
		go c.simulate(w)

		var sz size.Event // window size
		var m mouse.Event // latest mouse event
		var k key.Event   // latest key event
		// event loop
		for {
			publish := false

			e := w.NextEvent()
			// handle event
			switch e := e.(type) {
			case lifecycle.Event: // close button. BUG: doesn't exit from program.
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				k = e
			case mouse.Event:
				m = e
				// rescaling mouse coord
				m.Y = float32(c.height) * (m.Y / float32(sz.HeightPx))
				m.X = float32(c.width) * (m.X / float32(sz.WidthPx))

				if e.Direction == mouse.DirPress {
					c.context.mu.Lock()
					c.context.dragged = true
					c.context.mu.Unlock()
				}
				if e.Direction == mouse.DirRelease {
					c.context.mu.Lock()
					c.context.dragged = false
					c.context.mu.Unlock()
				}
			case paint.Event:
				publish = true

			case tickEvent:

				c.context.mu.Lock()
				// push latest mouse event to context
				c.context.pushMouseEvent(m)
				// copy image from shared memory
				copy(b.RGBA().Pix, c.context.pix())
				// set mouse event
				c.context.keyEvent = k
				c.context.mu.Unlock()

				// upload buffer to texture
				tex.Upload(image.Point{}, b, b.Bounds())
				publish = true
			case size.Event:
				sz = e
			case error:
				log.Print(e)
			}

			if publish {
				w.Scale(sz.Bounds(), tex, tex.Bounds(), draw.Src, nil)
				w.Publish()
			}
		}
	})
}
