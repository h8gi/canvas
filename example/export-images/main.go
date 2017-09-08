package main

import (
	"fmt"
	"math/rand"
	"os"

	"golang.org/x/image/colornames"

	"github.com/fogleman/gg"
	"github.com/h8gi/canvas"
)

func main() {
	c := canvas.New()

	c.Option(
		canvas.FrameRate(60),
		canvas.Size(200, 150),
	)

	world := NewWorld(200, 150)

	counter := 0
	c.Draw(func(dc *gg.Context) {
		if counter > 300 {
			os.Exit(0)
		}
		dc.SavePNG(fmt.Sprintf("%04d.png", counter))
		world.Update()
		for y := 0; y < dc.Height(); y++ {
			for x := 0; x < dc.Width(); x++ {
				if world.IsAliveAt(x, y) {
					dc.SetColor(colornames.White)
				} else {
					dc.SetColor(colornames.Black)
				}
				dc.SetPixel(x, y)
			}
		}
		counter++
	})
}

// lifegame
type World struct {
	matrix    [][]int
	neighbors [][]int
	width     int
	height    int
}

func NewWorld(w, h int) *World {
	m := make([][]int, h)
	n := make([][]int, h)
	for i := range m {
		m[i] = make([]int, w)
		n[i] = make([]int, w)
		for j := range m[i] {
			if rand.Float64() > 0.5 {
				m[i][j] = 1
			}
		}
	}
	return &World{
		matrix:    m,
		neighbors: n,
		width:     w,
		height:    h,
	}
}

func (world *World) IsAliveAt(x, y int) bool {
	return world.matrix[y][x] == 1
}

func (world *World) CountNeighbors(x, y int) int {
	var abs = func(n int) int {
		if n < 0 {
			return -n
		} else {
			return n
		}
	}
	top := abs((y - 1) % world.height)
	bottom := abs((y + 1) % world.height)
	left := abs((x - 1) % world.width)
	right := abs((x + 1) % world.width)
	return world.matrix[top][left] + world.matrix[top][x] + world.matrix[top][right] +
		world.matrix[y][left] + world.matrix[y][right] +
		world.matrix[bottom][left] + world.matrix[bottom][x] + world.matrix[bottom][right]
}

func (world *World) Update() {
	for x := 0; x < world.width; x++ {
		for y := 0; y < world.height; y++ {
			world.neighbors[y][x] = world.CountNeighbors(x, y)
		}
	}

	for x := 0; x < world.width; x++ {
		for y := 0; y < world.height; y++ {
			if world.IsAliveAt(x, y) {
				if 2 == world.neighbors[y][x] || 3 == world.neighbors[y][x] {
					world.matrix[y][x] = 1
				} else {
					world.matrix[y][x] = 0
				}
			} else {
				if world.neighbors[y][x] == 3 {
					world.matrix[y][x] = 1
				} else {
					world.matrix[y][x] = 0
				}
			}
		}
	}
}
