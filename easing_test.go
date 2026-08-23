package canvas

import (
	"math"
	"testing"
)

func TestEasingFunctions(t *testing.T) {
	easingFuncs := []struct {
		name string
		fn   func(float64) float64
	}{
		{"EaseLinear", EaseLinear},
		{"EaseInSine", EaseInSine},
		{"EaseOutSine", EaseOutSine},
		{"EaseInOutSine", EaseInOutSine},
		{"EaseInQuad", EaseInQuad},
		{"EaseOutQuad", EaseOutQuad},
		{"EaseInOutQuad", EaseInOutQuad},
		{"EaseInCubic", EaseInCubic},
		{"EaseOutCubic", EaseOutCubic},
		{"EaseInOutCubic", EaseInOutCubic},
		{"EaseInQuart", EaseInQuart},
		{"EaseOutQuart", EaseOutQuart},
		{"EaseInOutQuart", EaseInOutQuart},
		{"EaseInExpo", EaseInExpo},
		{"EaseOutExpo", EaseOutExpo},
		{"EaseInOutExpo", EaseInOutExpo},
		{"EaseInCirc", EaseInCirc},
		{"EaseOutCirc", EaseOutCirc},
		{"EaseInOutCirc", EaseInOutCirc},
		{"EaseInBack", EaseInBack},
		{"EaseOutBack", EaseOutBack},
		{"EaseInOutBack", EaseInOutBack},
		{"EaseInElastic", EaseInElastic},
		{"EaseOutElastic", EaseOutElastic},
		{"EaseInOutElastic", EaseInOutElastic},
		{"EaseInBounce", EaseInBounce},
		{"EaseOutBounce", EaseOutBounce},
		{"EaseInOutBounce", EaseInOutBounce},
	}

	for _, tt := range easingFuncs {
		t.Run(tt.name, func(t *testing.T) {
			// Boundary at t = 0
			val0 := tt.fn(0)
			if math.Abs(val0) > 1e-5 {
				t.Errorf("%s(0) expected ~0, got %v", tt.name, val0)
			}

			// Boundary at t = 1
			val1 := tt.fn(1)
			if math.Abs(val1-1.0) > 1e-5 {
				t.Errorf("%s(1) expected ~1, got %v", tt.name, val1)
			}

			// Midpoint check (should not be NaN or Inf)
			valMid := tt.fn(0.5)
			if math.IsNaN(valMid) || math.IsInf(valMid, 0) {
				t.Errorf("%s(0.5) returned invalid value %v", tt.name, valMid)
			}
		})
	}
}
