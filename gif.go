package canvas

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// GIFOptions holds configuration options for recording and encoding animated GIFs.
type GIFOptions struct {
	// Delay is the frame delay in 100ths of a second (10ms units).
	// If 0, delay is automatically calculated based on the canvas FrameRate (e.g. 60fps -> 2, 30fps -> 3).
	Delay int
	// LoopCount specifies the number of animation loops. 0 means infinite loop (default).
	LoopCount int
	// Dither enables Floyd-Steinberg dithering when converting RGBA frames to paletted images.
	// Default is true.
	Dither *bool
	// OnComplete is an optional callback invoked when GIF encoding finishes.
	OnComplete func(path string, err error)
}

type gifRecorder struct {
	mu             sync.Mutex
	active         bool
	path           string
	targetFrames   int
	capturedFrames []*image.RGBA
	opts           *GIFOptions
	frameRate      int
}

func newGIFRecorder(frameRate int) *gifRecorder {
	return &gifRecorder{
		frameRate: frameRate,
	}
}

func (r *gifRecorder) start(path string, frames int, opts *GIFOptions) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.path = path
	r.targetFrames = frames
	r.capturedFrames = make([]*image.RGBA, 0, frames)
	r.opts = opts
	r.active = true
}

func (r *gifRecorder) isRecording() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.active
}

// capture copies the current context RGBA image and triggers encoding when targetFrames is reached.
func (r *gifRecorder) capture(ctx *Context) {
	r.mu.Lock()
	if !r.active {
		r.mu.Unlock()
		return
	}

	src := ctx.Image().(*image.RGBA)
	// Deep copy the pixel buffer
	bounds := src.Bounds()
	frameCopy := image.NewRGBA(bounds)
	copy(frameCopy.Pix, src.Pix)
	r.capturedFrames = append(r.capturedFrames, frameCopy)

	if len(r.capturedFrames) >= r.targetFrames {
		r.active = false
		captured := r.capturedFrames
		r.capturedFrames = nil
		path := r.path
		opts := r.opts
		frameRate := r.frameRate
		r.mu.Unlock()

		go func() {
			err := saveCapturedFramesAsGIF(path, captured, frameRate, opts)
			if opts != nil && opts.OnComplete != nil {
				opts.OnComplete(path, err)
			}
		}()
		return
	}
	r.mu.Unlock()
}

func saveCapturedFramesAsGIF(path string, frames []*image.RGBA, frameRate int, opts *GIFOptions) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create GIF file %s: %w", path, err)
	}
	defer f.Close()

	images := make([]image.Image, len(frames))
	for i, frame := range frames {
		images[i] = frame
	}

	delay := 2
	if frameRate > 0 {
		delay = 100 / frameRate
		if delay < 1 {
			delay = 1
		}
	}
	if opts != nil && opts.Delay > 0 {
		delay = opts.Delay
	}

	return EncodeGIF(f, images, opts, delay)
}

// EncodeGIF encodes a slice of image.Image into an animated GIF written to w.
func EncodeGIF(w io.Writer, images []image.Image, opts *GIFOptions, defaultDelay ...int) error {
	if len(images) == 0 {
		return fmt.Errorf("no frames provided for GIF encoding")
	}

	delayVal := 3
	if len(defaultDelay) > 0 && defaultDelay[0] > 0 {
		delayVal = defaultDelay[0]
	}
	if opts != nil && opts.Delay > 0 {
		delayVal = opts.Delay
	}

	loopCount := 0
	if opts != nil {
		loopCount = opts.LoopCount
	}

	dither := true
	if opts != nil && opts.Dither != nil {
		dither = *opts.Dither
	}

	palettedFrames := make([]*image.Paletted, len(images))
	delays := make([]int, len(images))

	for i, img := range images {
		palettedFrames[i] = ImageToPaletted(img, dither)
		delays[i] = delayVal
	}

	anim := &gif.GIF{
		Image:     palettedFrames,
		Delay:     delays,
		LoopCount: loopCount,
	}

	return gif.EncodeAll(w, anim)
}

// SaveGIF saves a sequence of images as an animated GIF file.
func SaveGIF(path string, images []image.Image, opts *GIFOptions) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create GIF file %s: %w", path, err)
	}
	defer f.Close()

	return EncodeGIF(f, images, opts)
}

// ImageToPaletted converts an arbitrary image.Image to *image.Paletted using the standard Plan9 256-color palette.
func ImageToPaletted(img image.Image, dither bool) *image.Paletted {
	if paletted, ok := img.(*image.Paletted); ok {
		return paletted
	}

	bounds := img.Bounds()
	paletted := image.NewPaletted(bounds, palette.Plan9)

	if dither {
		draw.FloydSteinberg.Draw(paletted, bounds, img, bounds.Min)
	} else {
		draw.Src.Draw(paletted, bounds, img, bounds.Min)
	}

	return paletted
}
