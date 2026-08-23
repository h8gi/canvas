# Reaction-Diffusion Example (Gray-Scott Model)

An interactive Gray-Scott reaction-diffusion simulation running in real-time with `canvas`.

The Gray-Scott model simulates the chemical reaction and diffusion of two virtual substances, producing intricate self-organizing Turing patterns such as spots, coral structures, mitosis (cell division), maze-like fingerprints, and propagating spirals.

## How to Run

```sh
go run main.go
```

## Controls

- **Left Click / Drag**: Inject chemical $V$ onto the canvas to seed or disturb patterns.
- **Right Click / Drag**: Erase chemical $V$ from the canvas.
- **1 – 7 Keys**: Switch pattern presets:
  - **1**: Mitosis (Cell Division)
  - **2**: Coral Growth
  - **3**: Maze / Fingerprint
  - **4**: Soliton Spots
  - **5**: Pulsating Spirals
  - **6**: Dynamic Chaos
  - **7**: U-Skate Worlds
- **C Key**: Cycle through color palettes (*Neon Cyberpunk*, *Deep Ocean*, *Magma*, *Toxic Emerald*, *Monochrome*).
- **Space Key**: Reset to a random seeded initial state.
- **R Key**: Clear canvas.
- **P Key**: Pause / resume simulation.
- **H Key**: Toggle on-screen HUD (controls & parameters).
- **S Key**: Save current frame to `reaction_diffusion.png`.
