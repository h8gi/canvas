package canvas

import (
	"image/color"
	"testing"
)

func TestColorUtilities(t *testing.T) {
	t.Run("RGB and RGBA", func(t *testing.T) {
		c1 := RGB(255, 128, 0)
		if c1.R != 255 || c1.G != 128 || c1.B != 0 || c1.A != 255 {
			t.Errorf("unexpected RGB result: %+v", c1)
		}

		c2 := RGBA(10, 20, 30, 40)
		if c2.R != 10 || c2.G != 20 || c2.B != 30 || c2.A != 40 {
			t.Errorf("unexpected RGBA result: %+v", c2)
		}
	})

	t.Run("Hex", func(t *testing.T) {
		cRed := Hex("#ff0000")
		if cRed.R != 255 || cRed.G != 0 || cRed.B != 0 || cRed.A != 255 {
			t.Errorf("Hex(#ff0000) expected red, got %+v", cRed)
		}

		cShort := Hex("#0f0")
		if cShort.R != 0 || cShort.G != 255 || cShort.B != 0 || cShort.A != 255 {
			t.Errorf("Hex(#0f0) expected green, got %+v", cShort)
		}
	})

	t.Run("HSB to RGBA", func(t *testing.T) {
		// Red: H=0, S=1, B=1
		red := HSB(0, 1, 1)
		if red.R != 255 || red.G != 0 || red.B != 0 || red.A != 255 {
			t.Errorf("HSB(0, 1, 1) expected red [255 0 0 255], got %+v", red)
		}

		// Green: H=120, S=1, B=1
		green := HSB(120, 1, 1)
		if green.R != 0 || green.G != 255 || green.B != 0 || green.A != 255 {
			t.Errorf("HSB(120, 1, 1) expected green [0 255 0 255], got %+v", green)
		}

		// Blue: H=240, S=1, B=1
		blue := HSB(240, 1, 1)
		if blue.R != 0 || blue.G != 0 || blue.B != 255 || blue.A != 255 {
			t.Errorf("HSB(240, 1, 1) expected blue [0 0 255 255], got %+v", blue)
		}

		// White: H=0, S=0, B=1
		white := HSB(0, 0, 1)
		if white.R != 255 || white.G != 255 || white.B != 255 || white.A != 255 {
			t.Errorf("HSB(0, 0, 1) expected white, got %+v", white)
		}

		// Alpha support
		semiRed := HSBA(0, 1, 1, 0.5)
		if semiRed.R != 255 || semiRed.A != 128 && semiRed.A != 127 {
			t.Errorf("HSBA alpha expected ~128, got %d", semiRed.A)
		}
	})

	t.Run("HSL to RGBA", func(t *testing.T) {
		red := HSL(0, 1, 0.5)
		if red.R != 255 || red.G != 0 || red.B != 0 || red.A != 255 {
			t.Errorf("HSL(0, 1, 0.5) expected red, got %+v", red)
		}
	})

	t.Run("LerpColor", func(t *testing.T) {
		black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

		c0 := LerpColor(black, white, 0.0)
		if c0 != black {
			t.Errorf("LerpColor amt 0.0 expected black, got %+v", c0)
		}

		c1 := LerpColor(black, white, 1.0)
		if c1 != white {
			t.Errorf("LerpColor amt 1.0 expected white, got %+v", c1)
		}

		cMid := LerpColor(black, white, 0.5)
		if cMid.R != 127 && cMid.R != 128 {
			t.Errorf("LerpColor amt 0.5 expected ~128, got %+v", cMid)
		}
	})
}

func TestContext_ColorHelpers(t *testing.T) {
	ctx := NewContext(10, 10)

	ctx.BackgroundHSB(120, 1, 1)
	pix := ctx.pix()
	if pix[0] != 0 || pix[1] != 255 || pix[2] != 0 {
		t.Errorf("BackgroundHSB expected green, got [%d %d %d]", pix[0], pix[1], pix[2])
	}

	ctx.FillHSB(0, 1, 1)
	ctx.StrokeHSB(240, 1, 1)
	ctx.FillHSBA(60, 1, 1, 0.8)
	ctx.StrokeHSBA(180, 1, 1, 0.5)
}
