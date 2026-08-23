package main

import (
	"math"

	"github.com/h8gi/canvas"
)

type Boid struct {
	pos   canvas.Vec
	vel   canvas.Vec
	acc   canvas.Vec
	maxS  float64
	maxF  float64
	hue   float64
}

func newBoid(x, y float64) *Boid {
	angle := canvas.Random(0, math.Pi*2)
	return &Boid{
		pos:  canvas.V(x, y),
		vel:  canvas.V(math.Cos(angle)*2, math.Sin(angle)*2),
		acc:  canvas.V(0, 0),
		maxS: 3.5,
		maxF: 0.08,
		hue:  canvas.Random(180, 260),
	}
}

func (b *Boid) update(w, h float64) {
	b.vel = canvas.VecAdd(b.vel, b.acc)
	b.vel = canvas.VecLimit(b.vel, b.maxS)
	b.pos = canvas.VecAdd(b.pos, b.vel)
	b.acc = canvas.V(0, 0)

	// Wrap around edges
	if b.pos.X < 0 {
		b.pos.X += w
	}
	if b.pos.X >= w {
		b.pos.X -= w
	}
	if b.pos.Y < 0 {
		b.pos.Y += h
	}
	if b.pos.Y >= h {
		b.pos.Y -= h
	}
}

func (b *Boid) applyForce(force canvas.Vec) {
	b.acc = canvas.VecAdd(b.acc, force)
}

func (b *Boid) flock(boids []*Boid, mouse canvas.Vec, mousePressed bool) {
	sep := b.separate(boids)
	ali := b.align(boids)
	coh := b.cohere(boids)

	sep = canvas.VecMult(sep, 1.5)
	ali = canvas.VecMult(ali, 1.0)
	coh = canvas.VecMult(coh, 1.0)

	b.applyForce(sep)
	b.applyForce(ali)
	b.applyForce(coh)

	if mousePressed {
		seekMouse := b.seek(mouse)
		b.applyForce(canvas.VecMult(seekMouse, 2.0))
	}
}

func (b *Boid) seek(target canvas.Vec) canvas.Vec {
	desired := canvas.VecSub(target, b.pos)
	desired = canvas.VecSetMag(desired, b.maxS)
	steer := canvas.VecSub(desired, b.vel)
	return canvas.VecLimit(steer, b.maxF)
}

func (b *Boid) separate(boids []*Boid) canvas.Vec {
	desiredSeparation := 25.0
	steer := canvas.V(0, 0)
	count := 0

	for _, other := range boids {
		d := canvas.VecDist(b.pos, other.pos)
		if d > 0 && d < desiredSeparation {
			diff := canvas.VecSub(b.pos, other.pos)
			diff = canvas.VecNormalize(diff)
			diff = canvas.VecDiv(diff, d) // Closer boids have greater influence
			steer = canvas.VecAdd(steer, diff)
			count++
		}
	}

	if count > 0 {
		steer = canvas.VecDiv(steer, float64(count))
	}
	if canvas.VecMag(steer) > 0 {
		steer = canvas.VecSetMag(steer, b.maxS)
		steer = canvas.VecSub(steer, b.vel)
		steer = canvas.VecLimit(steer, b.maxF)
	}
	return steer
}

func (b *Boid) align(boids []*Boid) canvas.Vec {
	neighborDist := 50.0
	sum := canvas.V(0, 0)
	count := 0

	for _, other := range boids {
		d := canvas.VecDist(b.pos, other.pos)
		if d > 0 && d < neighborDist {
			sum = canvas.VecAdd(sum, other.vel)
			count++
		}
	}

	if count > 0 {
		sum = canvas.VecDiv(sum, float64(count))
		sum = canvas.VecSetMag(sum, b.maxS)
		steer := canvas.VecSub(sum, b.vel)
		return canvas.VecLimit(steer, b.maxF)
	}
	return canvas.V(0, 0)
}

func (b *Boid) cohere(boids []*Boid) canvas.Vec {
	neighborDist := 50.0
	sum := canvas.V(0, 0)
	count := 0

	for _, other := range boids {
		d := canvas.VecDist(b.pos, other.pos)
		if d > 0 && d < neighborDist {
			sum = canvas.VecAdd(sum, other.pos)
			count++
		}
	}

	if count > 0 {
		sum = canvas.VecDiv(sum, float64(count))
		return b.seek(sum)
	}
	return canvas.V(0, 0)
}

func (b *Boid) draw(ctx *canvas.Context) {
	angle := canvas.Heading(b.vel) + math.Pi/2

	ctx.Push()
	ctx.Translate(b.pos.X, b.pos.Y)
	ctx.Rotate(angle)

	ctx.FillHSBA(b.hue, 0.7, 0.95, 0.85)
	ctx.StrokeHSBA(b.hue, 0.9, 1.0, 1.0)
	ctx.SetLineWidth(1.2)

	// Draw sleek triangle pointing upwards
	ctx.DrawTriangle(0, -9, -5, 7, 5, 7)
	ctx.Fill()
	ctx.Stroke()

	ctx.Pop()
}

func main() {
	width, height := 800, 500
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     width,
		Height:    height,
		FrameRate: 60,
		Title:     "Boids Flocking Simulation",
	})

	boids := make([]*Boid, 120)

	c.Setup(func(ctx *canvas.Context) {
		for i := 0; i < len(boids); i++ {
			boids[i] = newBoid(canvas.Random(0, float64(width)), canvas.Random(0, float64(height)))
		}
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.BackgroundRGB(0.06, 0.07, 0.12)

		w, h := float64(ctx.Width()), float64(ctx.Height())
		isClicked := ctx.IsMousePressed(canvas.MouseLeft)

		for _, b := range boids {
			b.flock(boids, ctx.Mouse, isClicked)
			b.update(w, h)
			b.draw(ctx)
		}

		ctx.FillRGB(0.7, 0.7, 0.7)
		ctx.DrawString("Click & hold mouse to attract boids", 15, 25)
	})
}
