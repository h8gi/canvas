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

	t.Run("RandomGaussian", func(t *testing.T) {
		RandomSeed(42)
		val := RandomGaussian(0, 1)
		if math.IsNaN(val) || math.IsInf(val, 0) {
			t.Errorf("RandomGaussian returned invalid value: %v", val)
		}
	})

	t.Run("Vector Helpers", func(t *testing.T) {
		v0 := FromAngle(0)
		if math.Abs(v0.X-1.0) > 1e-9 || math.Abs(v0.Y-0.0) > 1e-9 {
			t.Errorf("FromAngle(0) expected (1, 0), got %+v", v0)
		}

		vPi2 := FromAngle(math.Pi / 2)
		if math.Abs(vPi2.X-0.0) > 1e-9 || math.Abs(vPi2.Y-1.0) > 1e-9 {
			t.Errorf("FromAngle(pi/2) expected (0, 1), got %+v", vPi2)
		}

		heading := Heading(vPi2)
		if math.Abs(heading-math.Pi/2) > 1e-9 {
			t.Errorf("Heading expected pi/2, got %v", heading)
		}

		vRand := Random2D()
		length := math.Hypot(vRand.X, vRand.Y)
		if math.Abs(length-1.0) > 1e-9 {
			t.Errorf("Random2D expected length 1.0, got %v", length)
		}
	})
}

