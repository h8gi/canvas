package main

import (
	"testing"
)

func TestAllPresetsGrowth(t *testing.T) {
	for i, p := range presets {
		t.Run(p.Name, func(t *testing.T) {
			sim := NewSimulation(300, 200)
			sim.SetPreset(i)
			sim.Reset()

			// Run 300 frames (2400 steps)
			for frame := 0; frame < 300; frame++ {
				for step := 0; step < 8; step++ {
					sim.Step()
				}
			}

			activeCount := 0
			for _, v := range sim.v {
				if v > 0.1 {
					activeCount++
				}
			}

			t.Logf("Preset '%s' active cells: %d / %d", p.Name, activeCount, sim.w*sim.h)
			if activeCount == 0 {
				t.Errorf("Preset '%s' died out completely!", p.Name)
			}
		})
	}
}
