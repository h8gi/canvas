package canvas_test

import (
	"fmt"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func ExampleMap() {
	// Map 5 from range [0, 10] to [0, 100]
	mapped := canvas.Map(5, 0, 10, 0, 100)
	fmt.Printf("%.0f\n", mapped)
	// Output: 50
}

func ExampleLerp() {
	// Linear interpolation halfway between 10 and 20
	val := canvas.Lerp(10, 20, 0.5)
	fmt.Printf("%.0f\n", val)
	// Output: 15
}

func ExampleConstrain() {
	fmt.Printf("%.0f\n", canvas.Constrain(15, 0, 10))
	fmt.Printf("%.0f\n", canvas.Constrain(-5, 0, 10))
	fmt.Printf("%.0f\n", canvas.Constrain(5, 0, 10))
	// Output:
	// 10
	// 0
	// 5
}

func ExampleDist() {
	distance := canvas.Dist(0, 0, 3, 4)
	fmt.Printf("%.0f\n", distance)
	// Output: 5
}

func ExampleNewCanvas() {
	// Create canvas instance with configuration
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    400,
		FrameRate: 60,
		Title:     "Animation",
	})

	c.Setup(func(ctx *canvas.Context) {
		ctx.Background(colornames.White)
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.FillHex("#ff0000")
		ctx.DrawCircle(ctx.Mouse.X, ctx.Mouse.Y, 20)
		ctx.Fill()
	})

	// To start loop in your main():
	// c.Draw(...)
}
