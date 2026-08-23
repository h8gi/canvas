package canvas

import "math"

// VecAdd adds two vectors v1 and v2 and returns the result.
func VecAdd(v1, v2 Vec) Vec {
	return V(v1.X+v2.X, v1.Y+v2.Y)
}

// VecSub subtracts vector v2 from v1 and returns the result.
func VecSub(v1, v2 Vec) Vec {
	return V(v1.X-v2.X, v1.Y-v2.Y)
}

// VecMult multiplies a vector by a scalar and returns the result.
func VecMult(v Vec, scalar float64) Vec {
	return V(v.X*scalar, v.Y*scalar)
}

// VecDiv divides a vector by a scalar and returns the result.
// If scalar is 0, it returns the original vector unchanged.
func VecDiv(v Vec, scalar float64) Vec {
	if scalar == 0 {
		return v
	}
	return V(v.X/scalar, v.Y/scalar)
}

// VecMag calculates the magnitude (length) of the vector.
func VecMag(v Vec) float64 {
	return math.Hypot(v.X, v.Y)
}

// VecMagSq calculates the squared magnitude of the vector.
func VecMagSq(v Vec) float64 {
	return v.X*v.X + v.Y*v.Y
}

// VecNormalize returns a unit vector pointing in the same direction as v.
// If v is a zero vector, it returns V(0, 0).
func VecNormalize(v Vec) Vec {
	m := VecMag(v)
	if m == 0 {
		return V(0, 0)
	}
	return V(v.X/m, v.Y/m)
}

// VecSetMag returns a new vector with the same direction as v and the specified magnitude.
func VecSetMag(v Vec, length float64) Vec {
	return VecMult(VecNormalize(v), length)
}

// VecLimit constrains the magnitude of the vector to the given max value.
func VecLimit(v Vec, max float64) Vec {
	magSq := VecMagSq(v)
	if magSq > max*max {
		return VecMult(VecNormalize(v), max)
	}
	return v
}

// VecDist calculates the Euclidean distance between two vectors.
func VecDist(v1, v2 Vec) float64 {
	return Dist(v1.X, v1.Y, v2.X, v2.Y)
}

// VecDot calculates the dot product of two vectors.
func VecDot(v1, v2 Vec) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// VecCross calculates the 2D cross product (determinant) of two vectors.
func VecCross(v1, v2 Vec) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// VecRotate rotates a 2D vector by an angle in radians.
func VecRotate(v Vec, rad float64) Vec {
	cos := math.Cos(rad)
	sin := math.Sin(rad)
	return V(v.X*cos-v.Y*sin, v.X*sin+v.Y*cos)
}

// VecAngleBetween calculates the angle (in radians) between two vectors.
func VecAngleBetween(v1, v2 Vec) float64 {
	m1 := VecMag(v1)
	m2 := VecMag(v2)
	if m1 == 0 || m2 == 0 {
		return 0
	}
	dot := VecDot(v1, v2) / (m1 * m2)
	dot = Constrain(dot, -1.0, 1.0)
	return math.Acos(dot)
}

// VecLerp linear-interpolates between two vectors at increment amt (0.0 to 1.0).
func VecLerp(v1, v2 Vec, amt float64) Vec {
	return V(
		Lerp(v1.X, v2.X, amt),
		Lerp(v1.Y, v2.Y, amt),
	)
}
