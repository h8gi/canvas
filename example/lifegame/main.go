package main

import (
	"fmt"
	"math/rand"

	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/key"

	"github.com/h8gi/canvas"
)

func main() {
	world := NewWorld(200, 150)

	c := canvas.New(&canvas.NewCanvasOptions{
		Width:     200,
		Height:    150,
		FrameRate: 30,
	})
	stop := false
	c.Draw(func(ctx *canvas.Context) {
		if ctx.KeyPressed() {
			fmt.Println(ctx.KeyEvent())
			if ctx.KeyCode() == key.CodeS {
				stop = !stop
			}
		}
		if !stop {
			world.Update()
		}
		for y := 0; y < ctx.Height(); y++ {
			for x := 0; x < ctx.Width(); x++ {
				if world.IsAliveAt(x, y) {
					ctx.SetColor(colornames.White)
				} else {
					ctx.SetColor(colornames.Black)
				}
				ctx.SetPixel(x, y)
			}
		}
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
