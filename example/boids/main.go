// very simple boids simulation.
package main

import (
	"math"
	"math/rand"

	"github.com/h8gi/canvas"
)

func main() {
	c := canvas.New()
	c.Option(
		canvas.FrameRate(60),
		canvas.Size(400, 400),
	)
	w := NewWorld(5)

	c.Setup(func(dc *canvas.Context) {
		dc.SetRGB(0, 0, 1)
		dc.Clear()
	})

	c.Draw(func(dc *canvas.Context) {
		w.Draw(dc)
		w.Update()
	})
}

type Vec2D [2]float64

func (a Vec2D) Add(b Vec2D) Vec2D {
	return Vec2D{a[0] + b[0], a[1] + b[1]}
}

func (a Vec2D) Mag() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1])
}

func (a Vec2D) Scale(k float64) Vec2D {
	return Vec2D{k * a[0], k * a[1]}
}

func (a Vec2D) Unit() Vec2D {
	return a.Scale(1 / a.Mag())
}

func (a Vec2D) Dot(b Vec2D) float64 {
	return a[0]*b[0] + a[1]*b[1]
}

func (a Vec2D) Rotate(rad float64) Vec2D {
	return Vec2D{
		a[0]*math.Cos(rad) - a[1]*math.Sin(rad),
		a[0]*math.Sin(rad) + a[1]*math.Cos(rad),
	}

}

type Boid struct {
	Position Vec2D
	Velocity Vec2D
}

func NewRandomBoid() *Boid {
	b := &Boid{
		Position: Vec2D{0.5, 0.5},
		Velocity: Vec2D{rand.Float64(), rand.Float64()},
	}
	return b
}

func (b *Boid) Forward() {
	b.Position = b.Position.Add(b.Velocity)
}

func (b *Boid) Wiggle(s float64) {
	rad := (2*math.Pi*rand.Float64() - math.Pi) * s
	b.Velocity = b.Velocity.Rotate(rad)
}

type World struct {
	Boids  []*Boid
	Width  int
	Height int
}

func NewWorld(n int) (w *World) {
	boids := make([]*Boid, n)
	for i := range boids {
		boids[i] = NewRandomBoid()
	}
	return &World{
		Boids: boids,
	}
}

func (w *World) Update() {
	for i := range w.Boids {
		b := w.Boids[i]
		b.Wiggle(0.05)
		b.Forward()
	}
}

func (w *World) Draw(dc *canvas.Context) {
	dc.SetRGB(0, 1, 0)
	for i := range w.Boids {
		b := w.Boids[i]
		dc.DrawCircle(b.Position[0], b.Position[1], 10)
		dc.Fill()
	}
}
