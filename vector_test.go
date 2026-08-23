package canvas

import (
	"math"
	"testing"
)

func TestVectorOperations(t *testing.T) {
	t.Run("VecAdd and VecSub", func(t *testing.T) {
		v1 := V(3, 4)
		v2 := V(1, 2)

		add := VecAdd(v1, v2)
		if add.X != 4 || add.Y != 6 {
			t.Errorf("VecAdd expected (4,6), got (%v,%v)", add.X, add.Y)
		}

		sub := VecSub(v1, v2)
		if sub.X != 2 || sub.Y != 2 {
			t.Errorf("VecSub expected (2,2), got (%v,%v)", sub.X, sub.Y)
		}
	})

	t.Run("VecMult and VecDiv", func(t *testing.T) {
		v := V(2, 3)

		mult := VecMult(v, 2.5)
		if mult.X != 5 || mult.Y != 7.5 {
			t.Errorf("VecMult expected (5, 7.5), got (%v,%v)", mult.X, mult.Y)
		}

		div := VecDiv(v, 2)
		if div.X != 1 || div.Y != 1.5 {
			t.Errorf("VecDiv expected (1, 1.5), got (%v,%v)", div.X, div.Y)
		}

		// Division by 0 returns unmodified
		divZero := VecDiv(v, 0)
		if divZero != v {
			t.Errorf("VecDiv by 0 expected original vector, got (%v,%v)", divZero.X, divZero.Y)
		}
	})

	t.Run("VecMag and VecMagSq", func(t *testing.T) {
		v := V(3, 4)
		if VecMag(v) != 5.0 {
			t.Errorf("VecMag expected 5.0, got %v", VecMag(v))
		}
		if VecMagSq(v) != 25.0 {
			t.Errorf("VecMagSq expected 25.0, got %v", VecMagSq(v))
		}
	})

	t.Run("VecNormalize and VecSetMag", func(t *testing.T) {
		v := V(0, 10)
		norm := VecNormalize(v)
		if norm.X != 0 || norm.Y != 1.0 {
			t.Errorf("VecNormalize expected (0,1), got (%v,%v)", norm.X, norm.Y)
		}

		zeroNorm := VecNormalize(V(0, 0))
		if zeroNorm.X != 0 || zeroNorm.Y != 0 {
			t.Errorf("VecNormalize of (0,0) expected (0,0), got (%v,%v)", zeroNorm.X, zeroNorm.Y)
		}

		setMag := VecSetMag(V(3, 4), 10.0)
		if math.Abs(VecMag(setMag)-10.0) > 1e-6 {
			t.Errorf("VecSetMag expected magnitude 10.0, got %v", VecMag(setMag))
		}
		if math.Abs(setMag.X-6.0) > 1e-6 || math.Abs(setMag.Y-8.0) > 1e-6 {
			t.Errorf("VecSetMag expected (6,8), got (%v,%v)", setMag.X, setMag.Y)
		}
	})

	t.Run("VecLimit", func(t *testing.T) {
		v := V(30, 40) // mag = 50
		limited := VecLimit(v, 10)
		if math.Abs(VecMag(limited)-10.0) > 1e-6 {
			t.Errorf("VecLimit expected magnitude 10.0, got %v", VecMag(limited))
		}

		// Within limit should not change
		unlimited := VecLimit(V(3, 4), 10)
		if unlimited.X != 3 || unlimited.Y != 4 {
			t.Errorf("VecLimit within limit changed vector: (%v,%v)", unlimited.X, unlimited.Y)
		}
	})

	t.Run("VecDist, VecDot, and VecCross", func(t *testing.T) {
		v1 := V(0, 0)
		v2 := V(3, 4)
		if VecDist(v1, v2) != 5.0 {
			t.Errorf("VecDist expected 5.0, got %v", VecDist(v1, v2))
		}

		// Dot product of orthogonal vectors is 0
		if VecDot(V(1, 0), V(0, 1)) != 0 {
			t.Errorf("VecDot expected 0, got %v", VecDot(V(1, 0), V(0, 1)))
		}
		// Dot product of parallel vectors
		if VecDot(V(2, 3), V(4, 5)) != 23 {
			t.Errorf("VecDot expected 23, got %v", VecDot(V(2, 3), V(4, 5)))
		}

		// Cross product: (1, 0) x (0, 1) = 1
		if VecCross(V(1, 0), V(0, 1)) != 1 {
			t.Errorf("VecCross expected 1, got %v", VecCross(V(1, 0), V(0, 1)))
		}
	})

	t.Run("VecRotate and VecAngleBetween", func(t *testing.T) {
		v := V(1, 0)
		rotated := VecRotate(v, math.Pi/2)
		if math.Abs(rotated.X) > 1e-6 || math.Abs(rotated.Y-1.0) > 1e-6 {
			t.Errorf("VecRotate by 90 deg expected (0,1), got (%v,%v)", rotated.X, rotated.Y)
		}

		angle := VecAngleBetween(V(1, 0), V(0, 1))
		if math.Abs(angle-math.Pi/2) > 1e-6 {
			t.Errorf("VecAngleBetween expected pi/2, got %v", angle)
		}
	})

	t.Run("VecLerp", func(t *testing.T) {
		v1 := V(0, 10)
		v2 := V(10, 20)
		lerped := VecLerp(v1, v2, 0.5)
		if lerped.X != 5 || lerped.Y != 15 {
			t.Errorf("VecLerp at 0.5 expected (5,15), got (%v,%v)", lerped.X, lerped.Y)
		}
	})
}
