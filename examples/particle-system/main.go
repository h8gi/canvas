package main

import (
	"math"

	"github.com/h8gi/canvas"
)

type Particle struct {
	pos      canvas.Vec
	vel      canvas.Vec
	acc      canvas.Vec
	lifespan float64
	maxLife  float64
	radius   float64
	hue      float64
}

func newParticle(x, y, hue float64) *Particle {
	angle := canvas.Random(0, math.Pi*2)
	speed := canvas.Random(1.5, 6.0)
	life := canvas.Random(40, 80)
	return &Particle{
		pos:      canvas.V(x, y),
		vel:      canvas.V(math.Cos(angle)*speed, math.Sin(angle)*speed-2.0),
		acc:      canvas.V(0, 0.08), // gravity
		lifespan: life,
		maxLife:  life,
		radius:   canvas.Random(3, 8),
		hue:      hue,
	}
}

func (p *Particle) update() {
	p.vel = canvas.VecAdd(p.vel, p.acc)
	p.vel = canvas.VecMult(p.vel, 0.98) // air drag
	p.pos = canvas.VecAdd(p.pos, p.vel)
	p.lifespan--
}

func (p *Particle) isDead() bool {
	return p.lifespan <= 0
}

func (p *Particle) draw(ctx *canvas.Context) {
	normLife := p.lifespan / p.maxLife
	alpha := canvas.EaseOutQuad(normLife)
	r := p.radius * canvas.EaseOutCubic(normLife)

	ctx.FillHSBA(p.hue, 0.85, 1.0, alpha*0.8)
	ctx.DrawCircle(p.pos.X, p.pos.Y, r)
	ctx.Fill()
}

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     700,
		Height:    500,
		FrameRate: 60,
		Title:     "Interactive Particle Fountain",
	})

	var particles []*Particle
	baseHue := 0.0
	var trailLayer *canvas.Context

	c.Setup(func(ctx *canvas.Context) {
		trailLayer = canvas.CreateGraphics(ctx.Width(), ctx.Height())
		trailLayer.BackgroundRGB(0.04, 0.04, 0.08)
	})

	c.Draw(func(ctx *canvas.Context) {
		baseHue = math.Mod(baseHue+0.5, 360)

		// Fade previous frame on trail layer
		trailLayer.BackgroundRGBA(0.04, 0.04, 0.08, 0.15)

		// Spawn new particles from mouse or center
		spawnPos := ctx.Mouse
		if spawnPos.X == 0 && spawnPos.Y == 0 {
			spawnPos = canvas.V(float64(ctx.Width())/2, float64(ctx.Height())/2)
		}

		spawnCount := 5
		if ctx.IsMouseDragged {
			spawnCount = 15
		}

		for i := 0; i < spawnCount; i++ {
			hue := math.Mod(baseHue+canvas.Random(-30, 30), 360)
			particles = append(particles, newParticle(spawnPos.X, spawnPos.Y, hue))
		}

		// Update and draw particles on trail layer
		active := particles[:0]
		for _, p := range particles {
			p.update()
			if !p.isDead() {
				p.draw(trailLayer)
				active = append(active, p)
			}
		}
		particles = active

		// Blit trail layer onto main window context
		ctx.DrawGraphics(trailLayer, 0, 0)

		ctx.FillRGB(0.8, 0.8, 0.8)
		ctx.DrawString("Move mouse to emit particles, drag to erupt", 15, 25)
	})
}
