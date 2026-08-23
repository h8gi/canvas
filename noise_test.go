package canvas

import (
	"math"
	"testing"
)

func TestNoise(t *testing.T) {
	t.Run("Range [0.0, 1.0]", func(t *testing.T) {
		for x := -5.0; x <= 5.0; x += 0.3 {
			for y := -5.0; y <= 5.0; y += 0.3 {
				val := Noise2D(x, y)
				if val < 0.0 || val > 1.0 {
					t.Fatalf("Noise2D(%v, %v) = %v is out of [0.0, 1.0] range", x, y, val)
				}
			}
		}
	})

	t.Run("Continuity", func(t *testing.T) {
		x := 1.234
		y := 5.678
		delta := 0.001
		v1 := Noise2D(x, y)
		v2 := Noise2D(x+delta, y)
		diff := math.Abs(v1 - v2)
		if diff > 0.05 {
			t.Errorf("Noise lacks continuity: v1=%v, v2=%v, diff=%v", v1, v2, diff)
		}
	})

	t.Run("1D and 3D shortcuts", func(t *testing.T) {
		v1D := Noise1D(2.5)
		if v1D < 0.0 || v1D > 1.0 {
			t.Errorf("Noise1D out of range: %v", v1D)
		}
		v3D := Noise(1.0, 2.0, 3.0)
		if v3D < 0.0 || v3D > 1.0 {
			t.Errorf("Noise (3D) out of range: %v", v3D)
		}
	})

	t.Run("Seed reproducibility", func(t *testing.T) {
		NoiseSeed(42)
		val1 := Noise2D(3.14, 2.71)

		NoiseSeed(42)
		val2 := Noise2D(3.14, 2.71)

		if val1 != val2 {
			t.Errorf("Expected same noise output for seed 42, got %v and %v", val1, val2)
		}

		NoiseSeed(999)
		val3 := Noise2D(3.14, 2.71)
		if val1 == val3 {
			t.Errorf("Expected different noise output for different seeds, both were %v", val1)
		}
	})

	t.Run("FBM Noise", func(t *testing.T) {
		for x := 0.0; x <= 5.0; x += 0.5 {
			for y := 0.0; y <= 5.0; y += 0.5 {
				v2D := NoiseFBM(x, y, 4, 0.5, 2.0)
				if v2D < 0.0 || v2D > 1.0 {
					t.Fatalf("NoiseFBM out of bounds: %v", v2D)
				}
				v3D := NoiseFBM3D(x, y, 1.0, 4, 0.5, 2.0)
				if v3D < 0.0 || v3D > 1.0 {
					t.Fatalf("NoiseFBM3D out of bounds: %v", v3D)
				}
			}
		}

		v1D := NoiseFBM1D(3.5, 4, 0.5, 2.0)
		if v1D < 0.0 || v1D > 1.0 {
			t.Fatalf("NoiseFBM1D out of bounds: %v", v1D)
		}

		// Edge case: octaves <= 0
		if NoiseFBM(1.0, 1.0, 0, 0.5, 2.0) != 0 {
			t.Errorf("Expected 0 for 0 octaves")
		}
	})
}
