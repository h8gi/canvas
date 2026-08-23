package main

import (
	"math/rand"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     640,
		Height:    400,
		FrameRate: 30,
		Title:     "Static Generative Art (NoLoop Example)",
		Resizable: true,
	})

	c.Setup(func(ctx *canvas.Context) {
		// Stop continuous frame rendering (renders once on launch)
		c.NoLoop()
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.BackgroundHex("#0f172a")

		// Palette
		palette := []string{"#38bdf8", "#818cf8", "#c084fc", "#f472b6", "#fb7185", "#34d399"}

		// Seed for deterministic art rendering
		rand.Seed(12345)

		// Draw generative geometric grid pattern
		cols, rows := 8, 5
		cellW := float64(ctx.Width()) / float64(cols)
		cellH := (float64(ctx.Height()) - 50) / float64(rows)

		for r := 0; r < rows; r++ {
			for col := 0; col < cols; col++ {
				cx := float64(col)*cellW + cellW/2
				cy := float64(r)*cellH + cellH/2 + 10

				ctx.Push()
				ctx.Translate(cx, cy)
				ctx.Rotate(canvas.Radians(float64(rand.Intn(4) * 45)))

				colorHex := palette[rand.Intn(len(palette))]
				ctx.FillHex(colorHex)
				ctx.StrokeHex("#ffffff")
				ctx.SetLineWidth(1.2)

				shapeType := rand.Intn(3)
				size := cellW * 0.35
				switch shapeType {
				case 0:
					ctx.DrawCircle(0, 0, size)
					ctx.Fill()
					ctx.Stroke()
				case 1:
					ctx.DrawRectangle(-size, -size, size*2, size*2)
					ctx.Fill()
					ctx.Stroke()
				case 2:
					ctx.DrawRegularPolygon(6, 0, 0, size, 0)
					ctx.Fill()
					ctx.Stroke()
				}

				// Inner detail
				ctx.FillHex("#0f172a")
				ctx.DrawCircle(0, 0, size*0.4)
				ctx.Fill()

				ctx.Pop()
			}
		}

		// Title / Info
		ctx.FillColor(colornames.White)
		ctx.DrawStringAnchored("Generative Geometric Grid (Rendered once with NoLoop)", float64(ctx.Width())/2, 380, 0.5, 0.5)
	})
}

