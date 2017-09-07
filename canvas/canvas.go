// simple boids viewer
package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"sync"
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
	drawFunc  func()
	shared    struct {
		mu              sync.Mutex
		uploadEventSent bool
		img             *image.RGBA
	}
}

// self-referential functions pattern
// https://commandcenter.blogspot.jp/2014/01/self-referential-functions-and-design.html
type option func(*Canvas)

func (c *Canvas) Option(opts ...option) {
	for _, opt := range opts {
		opt(c)
	}
}

func FrameRate(r int) option {
	return func(c *Canvas) {
		c.frameRate = r
	}
}

// set window size
func Size(width, height int) option {
	return func(c *Canvas) {
		c.width = width
		c.height = height
	}
}

func New() *Canvas {
	c := &Canvas{
		width:     600,
		height:    400,
		frameRate: 60,
	}
	return c
}

// start main loop
func (c *Canvas) Main(f func(*image.RGBA)) {
	c.drawFunc = func() {
		c.shared.mu.Lock()
		f(c.shared.img)
		c.shared.mu.Unlock()
	}
	c.start()
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

func (c *Canvas) start() {
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

		c.shared.img = image.NewRGBA(b.Bounds())
		copy(c.shared.img.Pix, b.RGBA().Pix)

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
				c.shared.mu.Lock()
				copy(b.RGBA().Pix, c.shared.img.Pix)
				c.shared.mu.Unlock()
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
