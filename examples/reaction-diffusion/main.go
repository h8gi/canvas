package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

const (
	simWidth  = 300
	simHeight = 200
	scale     = 2
	winWidth  = simWidth * scale
	winHeight = simHeight * scale
)

type Preset struct {
	Name string
	F    float64
	K    float64
}

var presets = []Preset{
	{Name: "1: Mitosis (Cell Division)", F: 0.0367, K: 0.0649},
	{Name: "2: Coral Growth", F: 0.0545, K: 0.0620},
	{Name: "3: Maze / Fingerprint", F: 0.0400, K: 0.0600},
	{Name: "4: Soliton Spots", F: 0.0300, K: 0.0620},
	{Name: "5: Pulsating Spirals", F: 0.0180, K: 0.0510},
	{Name: "6: Dynamic Chaos", F: 0.0260, K: 0.0550},
	{Name: "7: U-Skate Worlds", F: 0.0620, K: 0.0610},
}

type Palette struct {
	Name string
	LUT  [256]color.RGBA
}

func makeLUT(stops []struct {
	pos float64
	r   float64
	g   float64
	b   float64
}) [256]color.RGBA {
	var lut [256]color.RGBA
	for i := 0; i < 256; i++ {
		t := float64(i) / 255.0
		// Find surrounding stops
		var c0, c1 struct {
			pos float64
			r   float64
			g   float64
			b   float64
		}
		c0 = stops[0]
		c1 = stops[len(stops)-1]
		for j := 0; j < len(stops)-1; j++ {
			if t >= stops[j].pos && t <= stops[j+1].pos {
				c0 = stops[j]
				c1 = stops[j+1]
				break
			}
		}
		rangeLen := c1.pos - c0.pos
		var factor float64
		if rangeLen > 1e-6 {
			factor = (t - c0.pos) / rangeLen
		}
		r := uint8(math.Round(math.Max(0, math.Min(255, (c0.r+(c1.r-c0.r)*factor)*255))))
		g := uint8(math.Round(math.Max(0, math.Min(255, (c0.g+(c1.g-c0.g)*factor)*255))))
		b := uint8(math.Round(math.Max(0, math.Min(255, (c0.b+(c1.b-c0.b)*factor)*255))))
		lut[i] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return lut
}

var palettes = []Palette{
	{
		Name: "Neon Cyberpunk",
		LUT: makeLUT([]struct {
			pos, r, g, b float64
		}{
			{0.00, 0.04, 0.04, 0.10}, // Deep midnight
			{0.20, 0.10, 0.02, 0.35}, // Dark violet
			{0.45, 0.85, 0.05, 0.55}, // Hot magenta
			{0.70, 0.00, 0.85, 0.90}, // Neon cyan
			{1.00, 0.98, 0.95, 0.40}, // Electric yellow
		}),
	},
	{
		Name: "Deep Ocean",
		LUT: makeLUT([]struct {
			pos, r, g, b float64
		}{
			{0.00, 0.02, 0.05, 0.12}, // Abyssal blue
			{0.30, 0.05, 0.25, 0.45}, // Deep teal
			{0.60, 0.10, 0.65, 0.65}, // Bright seafoam
			{0.85, 0.40, 0.90, 0.85}, // Aquamarine
			{1.00, 0.95, 1.00, 0.98}, // Foam white
		}),
	},
	{
		Name: "Magma",
		LUT: makeLUT([]struct {
			pos, r, g, b float64
		}{
			{0.00, 0.02, 0.01, 0.03}, // Charcoal
			{0.25, 0.45, 0.05, 0.15}, // Crimson
			{0.55, 0.90, 0.25, 0.05}, // Vivid orange
			{0.80, 0.98, 0.75, 0.10}, // Bright amber
			{1.00, 1.00, 0.98, 0.85}, // Incandescent white
		}),
	},
	{
		Name: "Toxic Emerald",
		LUT: makeLUT([]struct {
			pos, r, g, b float64
		}{
			{0.00, 0.02, 0.06, 0.04}, // Forest black
			{0.30, 0.05, 0.35, 0.15}, // Deep jade
			{0.65, 0.15, 0.85, 0.35}, // Neon emerald
			{0.85, 0.65, 0.95, 0.45}, // Lime glow
			{1.00, 0.95, 1.00, 0.85}, // Pale yellow
		}),
	},
	{
		Name: "Monochrome",
		LUT: makeLUT([]struct {
			pos, r, g, b float64
		}{
			{0.00, 0.05, 0.05, 0.05},
			{0.50, 0.50, 0.50, 0.50},
			{1.00, 0.95, 0.95, 0.95},
		}),
	},
}

type Simulation struct {
	w, h       int
	u, v       []float64
	nextU      []float64
	nextV      []float64
	du, dv     float64
	f, k       float64
	renderImg  *image.RGBA
	paletteIdx int
	presetIdx  int
	paused     bool
}

func NewSimulation(w, h int) *Simulation {
	size := w * h
	sim := &Simulation{
		w:         w,
		h:         h,
		u:         make([]float64, size),
		v:         make([]float64, size),
		nextU:     make([]float64, size),
		nextV:     make([]float64, size),
		du:        1.0,
		dv:        0.5,
		renderImg: image.NewRGBA(image.Rect(0, 0, w*scale, h*scale)),
	}
	sim.SetPreset(0)
	sim.Reset()
	return sim
}

func (s *Simulation) SetPreset(idx int) {
	if idx < 0 || idx >= len(presets) {
		return
	}
	s.presetIdx = idx
	s.f = presets[idx].F
	s.k = presets[idx].K
}

func (s *Simulation) Reset() {
	size := s.w * s.h
	for i := 0; i < size; i++ {
		s.u[i] = 1.0
		s.v[i] = 0.0
	}

	cx, cy := s.w/2, s.h/2

	// Center seed block with noise for symmetry breaking
	blockSize := 20
	for dy := -blockSize; dy <= blockSize; dy++ {
		py := (cy + dy + s.h) % s.h
		for dx := -blockSize; dx <= blockSize; dx++ {
			px := (cx + dx + s.w) % s.w
			idx := py*s.w + px
			s.u[idx] = 0.50 + canvas.Random(-0.02, 0.02)
			s.v[idx] = 0.25 + canvas.Random(-0.02, 0.02)
		}
	}

	// Several satellite seed spots
	s.AddChemical(cx-60, cy-35, 12)
	s.AddChemical(cx+60, cy+35, 12)
	s.AddChemical(cx+55, cy-40, 10)
	s.AddChemical(cx-55, cy+40, 10)
}

func (s *Simulation) Clear() {
	size := s.w * s.h
	for i := 0; i < size; i++ {
		s.u[i] = 1.0
		s.v[i] = 0.0
		s.nextU[i] = 1.0
		s.nextV[i] = 0.0
	}
}

func (s *Simulation) AddChemical(cx, cy, radius int) {
	r2 := radius * radius
	for dy := -radius; dy <= radius; dy++ {
		py := (cy + dy + s.h) % s.h
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= r2 {
				px := (cx + dx + s.w) % s.w
				idx := py*s.w + px
				s.u[idx] = 0.50 + canvas.Random(-0.05, 0.05)
				s.v[idx] = 0.25 + canvas.Random(-0.05, 0.05)
			}
		}
	}
}

func (s *Simulation) RemoveChemical(cx, cy, radius int) {
	r2 := radius * radius
	for dy := -radius; dy <= radius; dy++ {
		py := (cy + dy + s.h) % s.h
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= r2 {
				px := (cx + dx + s.w) % s.w
				idx := py*s.w + px
				s.u[idx] = 1.0
				s.v[idx] = 0.0
			}
		}
	}
}

func (s *Simulation) Step() {
	w, h := s.w, s.h
	u, v := s.u, s.v
	nu, nv := s.nextU, s.nextV
	du, dv := s.du, s.dv
	f, k := s.f, s.k
	dt := 1.0

	for y := 0; y < h; y++ {
		yUp := ((y - 1 + h) % h) * w
		yCenter := y * w
		yDown := ((y + 1) % h) * w

		for x := 0; x < w; x++ {
			xLeft := (x - 1 + w) % w
			xCenter := x
			xRight := (x + 1) % w

			idx := yCenter + xCenter

			// 9-point Laplacian stencil
			lapU := 0.2*(u[yUp+xCenter]+u[yDown+xCenter]+u[yCenter+xLeft]+u[yCenter+xRight]) +
				0.05*(u[yUp+xLeft]+u[yUp+xRight]+u[yDown+xLeft]+u[yDown+xRight]) -
				u[idx]

			lapV := 0.2*(v[yUp+xCenter]+v[yDown+xCenter]+v[yCenter+xLeft]+v[yCenter+xRight]) +
				0.05*(v[yUp+xLeft]+v[yUp+xRight]+v[yDown+xLeft]+v[yDown+xRight]) -
				v[idx]

			uVal := u[idx]
			vVal := v[idx]
			uvv := uVal * vVal * vVal

			newU := uVal + (du*lapU-uvv+f*(1.0-uVal))*dt
			newV := vVal + (dv*lapV+uvv-(f+k)*vVal)*dt

			// Clamp to [0, 1]
			if newU < 0 {
				newU = 0
			} else if newU > 1 {
				newU = 1
			}
			if newV < 0 {
				newV = 0
			} else if newV > 1 {
				newV = 1
			}

			nu[idx] = newU
			nv[idx] = newV
		}
	}

	// Swap buffers
	copy(s.u, s.nextU)
	copy(s.v, s.nextV)
}

func (s *Simulation) UpdateTexture() {
	lut := palettes[s.paletteIdx].LUT
	pix := s.renderImg.Pix
	stride := s.w * scale * 4

	for y := 0; y < s.h; y++ {
		yOffset := y * s.w
		for x := 0; x < s.w; x++ {
			idx := yOffset + x
			// Value based on concentration of V
			val := s.v[idx]
			cIdx := int(val * 255.0)
			if cIdx < 0 {
				cIdx = 0
			} else if cIdx > 255 {
				cIdx = 255
			}
			col := lut[cIdx]

			// Write 2x2 pixels to upscale smoothly
			for dy := 0; dy < scale; dy++ {
				rowStart := (y*scale+dy)*stride + (x * scale * 4)
				for dx := 0; dx < scale; dx++ {
					pIdx := rowStart + dx*4
					pix[pIdx] = col.R
					pix[pIdx+1] = col.G
					pix[pIdx+2] = col.B
					pix[pIdx+3] = 255
				}
			}
		}
	}
}

func main() {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     winWidth,
		Height:    winHeight,
		FrameRate: 60,
		Title:     "Reaction-Diffusion (Gray-Scott System)",
	})

	sim := NewSimulation(simWidth, simHeight)
	showHelp := true

	c.Setup(func(ctx *canvas.Context) {
		ctx.BackgroundHex("#0a0a14")
	})

	c.KeyPressed(func(ctx *canvas.Context, k canvas.Key) {
		switch k {
		case canvas.Key1:
			sim.SetPreset(0)
		case canvas.Key2:
			sim.SetPreset(1)
		case canvas.Key3:
			sim.SetPreset(2)
		case canvas.Key4:
			sim.SetPreset(3)
		case canvas.Key5:
			sim.SetPreset(4)
		case canvas.Key6:
			sim.SetPreset(5)
		case canvas.Key7:
			sim.SetPreset(6)
		case canvas.KeyC:
			sim.paletteIdx = (sim.paletteIdx + 1) % len(palettes)
		case canvas.KeySpace:
			sim.Reset()
		case canvas.KeyR:
			sim.Clear()
		case canvas.KeyP:
			sim.paused = !sim.paused
		case canvas.KeyH:
			showHelp = !showHelp
		case canvas.KeyS:
			_ = ctx.SaveFrame("reaction_diffusion.png")
		}
	})

	c.Draw(func(ctx *canvas.Context) {
		// Mouse interaction: Left click to inject chemical V, Right click to erase
		if ctx.IsMousePressed(canvas.MouseLeft) || ctx.IsMouseDragged {
			mx := int(ctx.Mouse.X) / scale
			my := int(ctx.Mouse.Y) / scale
			sim.AddChemical(mx, my, 6)
		} else if ctx.IsMousePressed(canvas.MouseRight) {
			mx := int(ctx.Mouse.X) / scale
			my := int(ctx.Mouse.Y) / scale
			sim.RemoveChemical(mx, my, 10)
		}

		// Run multiple simulation sub-steps per render frame for smooth real-time evolution
		if !sim.paused {
			for step := 0; step < 8; step++ {
				sim.Step()
			}
		}

		// Update pixel buffer and render to canvas
		sim.UpdateTexture()
		ctx.DrawImage(sim.renderImg, 0, 0)

		// On-screen HUD / Controls overlay
		if showHelp {
			ctx.Push()
			// Background panel for text legibility
			ctx.SetRGBA(0.05, 0.05, 0.08, 0.8)
			ctx.DrawRoundedRectangle(12, 12, 360, 115, 6)
			ctx.Fill()

			ctx.FillColor(colornames.White)
			currentPreset := presets[sim.presetIdx]
			currentPalette := palettes[sim.paletteIdx]

			ctx.DrawString(fmt.Sprintf("Preset: %s", currentPreset.Name), 22, 32)
			ctx.SetRGB(0.7, 0.7, 0.8)
			ctx.DrawString(fmt.Sprintf("F: %.4f | k: %.4f | Palette: %s", currentPreset.F, currentPreset.K, currentPalette.Name), 22, 50)
			ctx.DrawString("Left Drag: Inject chemical V | Right Drag: Erase", 22, 70)
			ctx.DrawString("1-7: Presets | C: Palette | Space: Seed | R: Clear", 22, 88)
			ctx.DrawString("P: Pause | H: Toggle HUD | S: Save PNG", 22, 106)
			ctx.Pop()
		}
	})
}
