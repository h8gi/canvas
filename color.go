package canvas

import (
	"fmt"
	"image/color"
	"math"
	"strings"
)

// RGB creates a new color.RGBA with the given red, green, blue values (0-255) and full opacity (255).
func RGB(r, g, b uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// RGBA creates a new color.RGBA with the given red, green, blue, and alpha values (0-255).
func RGBA(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: a}
}

// Hex parses a hexadecimal color string (e.g. "#ff0055", "#fff", "38bdf8") and returns color.RGBA.
func Hex(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b, a uint8 = 0, 0, 0, 255
	switch len(hex) {
	case 3: // RGB (e.g. "fff")
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
	case 4: // RGBA (e.g. "ffff")
		fmt.Sscanf(hex, "%1x%1x%1x%1x", &r, &g, &b, &a)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
		a = (a << 4) | a
	case 6: // RRGGBB (e.g. "ffffff")
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	case 8: // RRGGBBAA (e.g. "ffffffff")
		fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}
	return color.RGBA{R: r, G: g, B: b, A: a}
}

// HSB converts Hue (0-360), Saturation (0-1), and Brightness (0-1) to color.RGBA with full alpha.
func HSB(h, s, b float64) color.RGBA {
	return HSBA(h, s, b, 1.0)
}

// HSBA converts Hue (0-360), Saturation (0-1), Brightness (0-1), and Alpha (0-1) to color.RGBA.
func HSBA(h, s, b, a float64) color.RGBA {
	h = math.Mod(h, 360.0)
	if h < 0 {
		h += 360.0
	}
	s = Constrain(s, 0.0, 1.0)
	b = Constrain(b, 0.0, 1.0)
	a = Constrain(a, 0.0, 1.0)

	c := b * s
	x := c * (1.0 - math.Abs(math.Mod(h/60.0, 2.0)-1.0))
	m := b - c

	var rPrime, gPrime, bPrime float64
	switch {
	case h < 60:
		rPrime, gPrime, bPrime = c, x, 0
	case h < 120:
		rPrime, gPrime, bPrime = x, c, 0
	case h < 180:
		rPrime, gPrime, bPrime = 0, c, x
	case h < 240:
		rPrime, gPrime, bPrime = 0, x, c
	case h < 300:
		rPrime, gPrime, bPrime = x, 0, c
	default:
		rPrime, gPrime, bPrime = c, 0, x
	}

	return color.RGBA{
		R: uint8(math.Round((rPrime + m) * 255.0)),
		G: uint8(math.Round((gPrime + m) * 255.0)),
		B: uint8(math.Round((bPrime + m) * 255.0)),
		A: uint8(math.Round(a * 255.0)),
	}
}

// HSL converts Hue (0-360), Saturation (0-1), and Lightness (0-1) to color.RGBA with full alpha.
func HSL(h, s, l float64) color.RGBA {
	return HSLA(h, s, l, 1.0)
}

// HSLA converts Hue (0-360), Saturation (0-1), Lightness (0-1), and Alpha (0-1) to color.RGBA.
func HSLA(h, s, l, a float64) color.RGBA {
	h = math.Mod(h, 360.0)
	if h < 0 {
		h += 360.0
	}
	s = Constrain(s, 0.0, 1.0)
	l = Constrain(l, 0.0, 1.0)
	a = Constrain(a, 0.0, 1.0)

	c := (1.0 - math.Abs(2.0*l-1.0)) * s
	x := c * (1.0 - math.Abs(math.Mod(h/60.0, 2.0)-1.0))
	m := l - c/2.0

	var rPrime, gPrime, bPrime float64
	switch {
	case h < 60:
		rPrime, gPrime, bPrime = c, x, 0
	case h < 120:
		rPrime, gPrime, bPrime = x, c, 0
	case h < 180:
		rPrime, gPrime, bPrime = 0, c, x
	case h < 240:
		rPrime, gPrime, bPrime = 0, x, c
	case h < 300:
		rPrime, gPrime, bPrime = x, 0, c
	default:
		rPrime, gPrime, bPrime = c, 0, x
	}

	return color.RGBA{
		R: uint8(math.Round((rPrime + m) * 255.0)),
		G: uint8(math.Round((gPrime + m) * 255.0)),
		B: uint8(math.Round((bPrime + m) * 255.0)),
		A: uint8(math.Round(a * 255.0)),
	}
}

// LerpColor calculates a color between two colors at a specific increment (amt between 0.0 and 1.0).
func LerpColor(c1, c2 color.Color, amt float64) color.RGBA {
	amt = Constrain(amt, 0.0, 1.0)

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	// RGBA() returns alpha-premultiplied uint32 in range [0, 65535]
	r := uint8(math.Round((float64(r1>>8) + amt*(float64(r2>>8)-float64(r1>>8)))))
	g := uint8(math.Round((float64(g1>>8) + amt*(float64(g2>>8)-float64(g1>>8)))))
	b := uint8(math.Round((float64(b1>>8) + amt*(float64(b2>>8)-float64(b1>>8)))))
	a := uint8(math.Round((float64(a1>>8) + amt*(float64(a2>>8)-float64(a1>>8)))))

	return color.RGBA{R: r, G: g, B: b, A: a}
}
