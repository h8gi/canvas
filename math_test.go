package canvas

import (
	"math"
	"testing"
)

func TestMathUtilities(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		got := Map(5, 0, 10, 0, 100)
		if got != 50 {
			t.Errorf("expected 50, got %v", got)
		}
		got2 := Map(0.5, 0, 1, -10, 10)
		if got2 != 0 {
			t.Errorf("expected 0, got %v", got2)
		}
	})

	t.Run("Lerp", func(t *testing.T) {
		got := Lerp(10, 20, 0.5)
		if got != 15 {
			t.Errorf("expected 15, got %v", got)
		}
	})

	t.Run("Constrain", func(t *testing.T) {
		if got := Constrain(5, 0, 10); got != 5 {
			t.Errorf("expected 5, got %v", got)
		}
		if got := Constrain(-5, 0, 10); got != 0 {
			t.Errorf("expected 0, got %v", got)
		}
		if got := Constrain(15, 0, 10); got != 10 {
			t.Errorf("expected 10, got %v", got)
		}
	})

	t.Run("Dist", func(t *testing.T) {
		got := Dist(0, 0, 3, 4)
		if got != 5 {
			t.Errorf("expected 5, got %v", got)
		}
	})

	t.Run("Radians and Degrees", func(t *testing.T) {
		rad := Radians(180)
		if math.Abs(rad-math.Pi) > 1e-9 {
			t.Errorf("expected pi, got %v", rad)
		}
		deg := Degrees(math.Pi)
		if math.Abs(deg-180) > 1e-9 {
			t.Errorf("expected 180, got %v", deg)
		}
	})

	t.Run("Random", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			r := Random(10, 20)
			if r < 10 || r > 20 {
				t.Fatalf("Random out of range: %v", r)
			}
		}
	})
}
