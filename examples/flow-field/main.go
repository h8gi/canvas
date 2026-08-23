package main

import (
	"math"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     800,
		Height:    500,
		FrameRate: 60,
		Title:     "Perlin Noise Flow Field",
	})

	type Particle struct {
		x, y  float64
		hue   float64
		speed float64
	}

	var particles []Particle
	initParticles := func() {
		particles = make([]Particle, 800)
		for i := range particles {
			particles[i] = Particle{
				x:     canvas.Random(0, 800),
				y:     canvas.Random(0, 500),
				hue:   canvas.Random(0, 360),
				speed: canvas.Random(1.5, 3.5),
			}
		}
	}

	c.Setup(func(ctx *canvas.Context) {
		ctx.BackgroundHex("#0b0f19")
		initParticles()
	})

	c.KeyPressed(func(ctx *canvas.Context, k canvas.Key) {
		if k == canvas.KeySpace {
			ctx.BackgroundHex("#0b0f19")
			canvas.NoiseSeed(int64(canvas.Random(0, 10000)))
			initParticles()
		}
		if k == canvas.KeyS {
			_ = ctx.SaveFrame("flowfield.png")
		}
	})

	c.Draw(func(ctx *canvas.Context) {
		// Semi-transparent background overlay for smooth light trail effect
		ctx.Push()
		ctx.SetRGBA(0.04, 0.06, 0.1, 0.05)
		ctx.DrawRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()))
		ctx.Fill()
		ctx.Pop()

		scale := 0.005
		for i := range particles {
			p := &particles[i]
			angle := canvas.Noise2D(p.x*scale, p.y*scale) * math.Pi * 4.0

			nx := p.x + math.Cos(angle)*p.speed
			ny := p.y + math.Sin(angle)*p.speed

			ctx.Push()
			ctx.StrokeHSBA(p.hue, 0.85, 0.95, 0.6)
			ctx.SetLineWidth(1.5)
			ctx.DrawLine(p.x, p.y, nx, ny)
			ctx.Stroke()
			ctx.Pop()

			p.x = nx
			p.y = ny
			p.hue = math.Mod(p.hue+0.2, 360.0)

			// Wrap edges
			if p.x < 0 {
				p.x = float64(ctx.Width())
			} else if p.x > float64(ctx.Width()) {
				p.x = 0
			}
			if p.y < 0 {
				p.y = float64(ctx.Height())
			} else if p.y > float64(ctx.Height()) {
				p.y = 0
			}
		}

		// Help text on top
		if ctx.FrameCount < 180 {
			ctx.FillColor(colornames.White)
			ctx.DrawStringAnchored("Press SPACE to re-seed, S to save frame", float64(ctx.Width())/2, 25, 0.5, 0.5)
		}
	})
}
