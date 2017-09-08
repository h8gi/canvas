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
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

// timer event
type TickEvent struct{}

// drawing area
type Canvas struct {
	width     int
	height    int
	frameRate int
	initFunc  func()
	drawFunc  func()
	context   *Context
}

// self-referential functions pattern
// https://commandcenter.blogspot.jp/2014/01/self-referential-functions-and-design.html
type option func(*Canvas)

func (c *Canvas) Option(opts ...option) {
	for _, opt := range opts {
		opt(c)
	}
}

func FrameRate(fps int) option {
	return func(c *Canvas) {
		c.frameRate = fps
	}
}

// set window size
func Size(width, height int) option {
	return func(c *Canvas) {
		c.width = width
		c.height = height
		c.context = NewContext(width, height)
	}
}

func New() *Canvas {
	c := &Canvas{}
	c.Option(
		Size(600, 400),
		FrameRate(60),
	)
	return c
}

// initialize
func (c *Canvas) Setup(f func(*Context)) {
	c.context.mu.Lock()
	f(c.context)
	c.context.mu.Unlock()
}

// start main loop
func (c *Canvas) Draw(f func(*Context)) {
	c.drawFunc = func() {
		c.context.mu.Lock()
		f(c.context)
		c.context.mu.Unlock()
	}
	c.startLoop()
}

// start inner simulation loop
func (c *Canvas) simulate(q screen.EventDeque) {
	duration := time.Second / time.Duration(c.frameRate)
	for {
		// memory lock
		c.drawFunc()
		q.Send(TickEvent{})
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
			Title:  "Basic Shiny Example",
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
				if e.Code == key.CodeEscape {
					return
				}
			case paint.Event:
				publish = true

			case TickEvent:
				c.context.mu.Lock()
				// copy image from shared memory
				copy(b.RGBA().Pix, c.context.Image().(*image.RGBA).Pix)
				c.context.mu.Unlock()
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
