package canvas

import "math"

// EaseLinear performs linear interpolation (constant speed).
func EaseLinear(t float64) float64 {
	return t
}

// EaseInSine accelerates using a sine curve.
func EaseInSine(t float64) float64 {
	return 1.0 - math.Cos((t*math.Pi)/2.0)
}

// EaseOutSine decelerates using a sine curve.
func EaseOutSine(t float64) float64 {
	return math.Sin((t * math.Pi) / 2.0)
}

// EaseInOutSine accelerates then decelerates using a sine curve.
func EaseInOutSine(t float64) float64 {
	return -(math.Cos(math.Pi*t) - 1.0) / 2.0
}

// EaseInQuad accelerates using a quadratic curve (t^2).
func EaseInQuad(t float64) float64 {
	return t * t
}

// EaseOutQuad decelerates using a quadratic curve.
func EaseOutQuad(t float64) float64 {
	return 1.0 - (1.0-t)*(1.0-t)
}

// EaseInOutQuad accelerates then decelerates using a quadratic curve.
func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2.0 * t * t
	}
	return 1.0 - math.Pow(-2.0*t+2.0, 2.0)/2.0
}

// EaseInCubic accelerates using a cubic curve (t^3).
func EaseInCubic(t float64) float64 {
	return t * t * t
}

// EaseOutCubic decelerates using a cubic curve.
func EaseOutCubic(t float64) float64 {
	return 1.0 - math.Pow(1.0-t, 3.0)
}

// EaseInOutCubic accelerates then decelerates using a cubic curve.
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4.0 * t * t * t
	}
	return 1.0 - math.Pow(-2.0*t+2.0, 3.0)/2.0
}

// EaseInQuart accelerates using a quartic curve (t^4).
func EaseInQuart(t float64) float64 {
	return t * t * t * t
}

// EaseOutQuart decelerates using a quartic curve.
func EaseOutQuart(t float64) float64 {
	return 1.0 - math.Pow(1.0-t, 4.0)
}

// EaseInOutQuart accelerates then decelerates using a quartic curve.
func EaseInOutQuart(t float64) float64 {
	if t < 0.5 {
		return 8.0 * t * t * t * t
	}
	return 1.0 - math.Pow(-2.0*t+2.0, 4.0)/2.0
}

// EaseInExpo accelerates exponentially.
func EaseInExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	return math.Pow(2.0, 10.0*t-10.0)
}

// EaseOutExpo decelerates exponentially.
func EaseOutExpo(t float64) float64 {
	if t >= 1.0 {
		return 1.0
	}
	return 1.0 - math.Pow(2.0, -10.0*t)
}

// EaseInOutExpo accelerates then decelerates exponentially.
func EaseInOutExpo(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1.0 {
		return 1.0
	}
	if t < 0.5 {
		return math.Pow(2.0, 20.0*t-10.0) / 2.0
	}
	return (2.0 - math.Pow(2.0, -20.0*t+10.0)) / 2.0
}

// EaseInCirc accelerates using a circular curve.
func EaseInCirc(t float64) float64 {
	t = Constrain(t, 0.0, 1.0)
	return 1.0 - math.Sqrt(1.0-math.Pow(t, 2.0))
}

// EaseOutCirc decelerates using a circular curve.
func EaseOutCirc(t float64) float64 {
	t = Constrain(t, 0.0, 1.0)
	return math.Sqrt(1.0 - math.Pow(t-1.0, 2.0))
}

// EaseInOutCirc accelerates then decelerates using a circular curve.
func EaseInOutCirc(t float64) float64 {
	t = Constrain(t, 0.0, 1.0)
	if t < 0.5 {
		return (1.0 - math.Sqrt(1.0-math.Pow(2.0*t, 2.0))) / 2.0
	}
	return (math.Sqrt(1.0-math.Pow(-2.0*t+2.0, 2.0)) + 1.0) / 2.0
}

// EaseInBack accelerates with an initial backward overshoot.
func EaseInBack(t float64) float64 {
	c1 := 1.70158
	c3 := c1 + 1.0
	return c3*t*t*t - c1*t*t
}

// EaseOutBack decelerates with a final overshoot.
func EaseOutBack(t float64) float64 {
	c1 := 1.70158
	c3 := c1 + 1.0
	return 1.0 + c3*math.Pow(t-1.0, 3.0) + c1*math.Pow(t-1.0, 2.0)
}

// EaseInOutBack accelerates with backward overshoot then decelerates with overshoot.
func EaseInOutBack(t float64) float64 {
	c1 := 1.70158
	c2 := c1 * 1.525
	if t < 0.5 {
		return (math.Pow(2.0*t, 2.0) * ((c2+1.0)*2.0*t - c2)) / 2.0
	}
	return (math.Pow(2.0*t-2.0, 2.0)*((c2+1.0)*(t*2.0-2.0)+c2) + 2.0) / 2.0
}

// EaseInElastic accelerates with an elastic bounce at the start.
func EaseInElastic(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1.0 {
		return 1.0
	}
	c4 := (2.0 * math.Pi) / 3.0
	return -math.Pow(2.0, 10.0*t-10.0) * math.Sin((t*10.0-10.75)*c4)
}

// EaseOutElastic decelerates with an elastic bounce at the end.
func EaseOutElastic(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1.0 {
		return 1.0
	}
	c4 := (2.0 * math.Pi) / 3.0
	return math.Pow(2.0, -10.0*t)*math.Sin((t*10.0-0.75)*c4) + 1.0
}

// EaseInOutElastic accelerates with elastic bounce, then decelerates with bounce.
func EaseInOutElastic(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1.0 {
		return 1.0
	}
	c5 := (2.0 * math.Pi) / 4.5
	if t < 0.5 {
		return -(math.Pow(2.0, 20.0*t-10.0) * math.Sin((20.0*t-11.125)*c5)) / 2.0
	}
	return (math.Pow(2.0, -20.0*t+10.0)*math.Sin((20.0*t-11.125)*c5))/2.0 + 1.0
}

// EaseOutBounce decelerates with a bouncing effect.
func EaseOutBounce(t float64) float64 {
	n1 := 7.5625
	d1 := 2.75

	if t < 1.0/d1 {
		return n1 * t * t
	} else if t < 2.0/d1 {
		t -= 1.5 / d1
		return n1*t*t + 0.75
	} else if t < 2.5/d1 {
		t -= 2.25 / d1
		return n1*t*t + 0.9375
	} else {
		t -= 2.625 / d1
		return n1*t*t + 0.984375
	}
}

// EaseInBounce accelerates with a bouncing effect.
func EaseInBounce(t float64) float64 {
	return 1.0 - EaseOutBounce(1.0-t)
}

// EaseInOutBounce accelerates with bouncing then decelerates with bouncing.
func EaseInOutBounce(t float64) float64 {
	if t < 0.5 {
		return (1.0 - EaseOutBounce(1.0-2.0*t)) / 2.0
	}
	return (1.0 + EaseOutBounce(2.0*t-1.0)) / 2.0
}
