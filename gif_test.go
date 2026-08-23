package canvas

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestEncodeGIF(t *testing.T) {
	// Create 3 simple dummy frames (10x10) with different colors
	frames := make([]image.Image, 3)
	colors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},
		{R: 0, G: 255, B: 0, A: 255},
		{R: 0, G: 0, B: 255, A: 255},
	}

	for i := 0; i < 3; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				img.Set(x, y, colors[i])
			}
		}
		frames[i] = img
	}

	t.Run("default options", func(t *testing.T) {
		var buf bytes.Buffer
		err := EncodeGIF(&buf, frames, nil)
		if err != nil {
			t.Fatalf("EncodeGIF failed: %v", err)
		}

		decoded, err := gif.DecodeAll(&buf)
		if err != nil {
			t.Fatalf("failed to decode generated GIF: %v", err)
		}

		if len(decoded.Image) != 3 {
			t.Errorf("expected 3 frames, got %d", len(decoded.Image))
		}
		if decoded.LoopCount != 0 {
			t.Errorf("expected LoopCount 0, got %d", decoded.LoopCount)
		}
	})

	t.Run("custom options with delay and loop count", func(t *testing.T) {
		var buf bytes.Buffer
		dither := false
		opts := &GIFOptions{
			Delay:     10,
			LoopCount: 5,
			Dither:    &dither,
		}
		err := EncodeGIF(&buf, frames, opts)
		if err != nil {
			t.Fatalf("EncodeGIF failed: %v", err)
		}

		decoded, err := gif.DecodeAll(&buf)
		if err != nil {
			t.Fatalf("failed to decode generated GIF: %v", err)
		}

		if len(decoded.Image) != 3 {
			t.Errorf("expected 3 frames, got %d", len(decoded.Image))
		}
		if decoded.LoopCount != 5 {
			t.Errorf("expected LoopCount 5, got %d", decoded.LoopCount)
		}
		for i, d := range decoded.Delay {
			if d != 10 {
				t.Errorf("expected frame %d delay to be 10, got %d", i, d)
			}
		}
	})

	t.Run("empty frames error", func(t *testing.T) {
		var buf bytes.Buffer
		err := EncodeGIF(&buf, nil, nil)
		if err == nil {
			t.Error("expected error for empty frames, got nil")
		}
	})
}

func TestSaveGIF(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "sub", "test.gif")

	img1 := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img2 := image.NewRGBA(image.Rect(0, 0, 8, 8))
	frames := []image.Image{img1, img2}

	err := SaveGIF(outPath, frames, nil)
	if err != nil {
		t.Fatalf("SaveGIF failed: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("GIF file was not created: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("GIF file is empty")
	}
}

func TestGIFRecorder(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "record.gif")

	rec := newGIFRecorder(60)
	var wg sync.WaitGroup
	wg.Add(1)

	var callbackPath string
	var callbackErr error

	opts := &GIFOptions{
		OnComplete: func(path string, err error) {
			callbackPath = path
			callbackErr = err
			wg.Done()
		},
	}

	rec.start(outPath, 3, opts)
	if !rec.isRecording() {
		t.Error("expected recorder to be active")
	}

	ctx := NewContext(10, 10)

	// Capture frame 1
	ctx.BackgroundRGB(1, 0, 0)
	rec.capture(ctx)
	if !rec.isRecording() {
		t.Error("expected recorder to still be active after frame 1")
	}

	// Capture frame 2
	ctx.BackgroundRGB(0, 1, 0)
	rec.capture(ctx)
	if !rec.isRecording() {
		t.Error("expected recorder to still be active after frame 2")
	}

	// Capture frame 3 (triggers complete)
	ctx.BackgroundRGB(0, 0, 1)
	rec.capture(ctx)
	if rec.isRecording() {
		t.Error("expected recorder to become inactive after reaching target frames")
	}

	// Wait for async encoding
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if callbackErr != nil {
			t.Fatalf("recording complete callback returned error: %v", callbackErr)
		}
		if callbackPath != outPath {
			t.Errorf("expected callback path %s, got %s", outPath, callbackPath)
		}
		// Verify file exists
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			t.Errorf("expected recorded GIF file to exist: %s", outPath)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for GIF recording completion")
	}
}

func TestCanvas_RecordGIF_Integration(t *testing.T) {
	c := NewCanvas(nil)
	if c.IsRecordingGIF() {
		t.Error("expected IsRecordingGIF to be false initially")
	}

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "canvas_record.gif")

	c.RecordGIF(outPath, 10)
	if !c.IsRecordingGIF() {
		t.Error("expected IsRecordingGIF to be true after RecordGIF call")
	}
}

func TestContext_CloneImage(t *testing.T) {
	ctx := NewContext(10, 10)
	ctx.BackgroundRGB(1, 0, 0)

	cloned := ctx.CloneImage()
	if cloned == nil {
		t.Fatal("expected cloned image to not be nil")
	}

	// Verify cloned pixels match
	if cloned.Pix[0] != 255 || cloned.Pix[1] != 0 || cloned.Pix[2] != 0 {
		t.Errorf("cloned image pixel mismatch")
	}

	// Modify original context, cloned should remain unchanged
	ctx.BackgroundRGB(0, 0, 1)
	if cloned.Pix[0] != 255 || cloned.Pix[2] != 0 {
		t.Errorf("cloned image was mutated when context changed")
	}
}
