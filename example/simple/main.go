package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/fogleman/gg"
	"github.com/h8gi/boids/canvas"
)

func main() {
	c := canvas.New()
	c.Option(
		canvas.FrameRate(60),
		canvas.Size(400, 400),
	)
	w := NewWorld(5)
	counter := 0
	c.Main(func(dc *gg.Context) {
		if counter > 100 {
			os.Exit(0)
		}
		w.Draw(dc)
		dc.SavePNG(fmt.Sprintf("%04d.png", counter))
		w.Update()
		counter++
	})
}

type Vector [2]float64

func (v Vector) Add(other Vector) Vector {
	return Vector{v[0] + other[0], v[1] + other[1]}
}

type Boid struct {
	Position Vector
	Velocity Vector
}

func NewRandomBoid() *Boid {
	b := &Boid{
		Position: Vector{rand.Float64(), rand.Float64()},
		Velocity: Vector{rand.Float64(), rand.Float64()},
	}
	return b
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
		b.Position = b.Position.Add(b.Velocity)
	}
}

func (w *World) Draw(dc *gg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(0, 1, 0)
	for i := range w.Boids {
		b := w.Boids[i]
		dc.DrawCircle(b.Position[0], b.Position[1], 10)
		dc.Fill()
	}
}
