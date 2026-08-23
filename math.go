package canvas

import (
	"math"
	"math/rand"
)

// Map maps a value from one range [inMin, inMax] to another [outMin, outMax].
func Map(val, inMin, inMax, outMin, outMax float64) float64 {
	return outMin + (outMax-outMin)*((val-inMin)/(inMax-inMin))
}

// Lerp calculates a number between two numbers at a specific increment (amt: 0.0 to 1.0).
func Lerp(start, stop, amt float64) float64 {
	return start + amt*(stop-start)
}

// Constrain constrains a value to not exceed a maximum and minimum value.
func Constrain(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// Dist calculates the Euclidean distance between two points.
func Dist(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// Radians converts a degree measurement to its corresponding value in radians.
func Radians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// Degrees converts a radian measurement to its corresponding value in degrees.
func Degrees(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

// Random returns a pseudo-random number between min (inclusive) and max (exclusive).
func Random(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomRange returns a pseudo-random number between 0 (inclusive) and max (exclusive).
func RandomRange(max float64) float64 {
	return rand.Float64() * max
}

// RandomSeed sets the seed for the pseudo-random number generator.
func RandomSeed(seed int64) {
	rand.Seed(seed)
}

// RandomGaussian returns a pseudo-random float64 from a Gaussian (normal) distribution
// with the specified mean and standard deviation (stdDev).
func RandomGaussian(mean, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + mean
}

// FromAngle returns a 2D unit vector pointing in the direction of the given angle in radians.
func FromAngle(rad float64) Vec {
	return V(math.Cos(rad), math.Sin(rad))
}

// Random2D returns a random 2D unit vector.
func Random2D() Vec {
	return FromAngle(rand.Float64() * math.Pi * 2.0)
}

// Heading returns the angle of rotation (in radians) for the vector.
func Heading(v Vec) float64 {
	return math.Atan2(v.Y, v.X)
}
