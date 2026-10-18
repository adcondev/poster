package emulator

import (
	"sync"
	"testing"

	"golang.org/x/image/math/fixed"
)

// TestScaledFacesRace verifies that accessing scaledFaces from multiple goroutines
// does not cause a data race.
func TestScaledFacesRace(t *testing.T) {
	fm := NewFontManager()

	// Load a font first (assuming JetBrainsMono-Regular.ttf is available in embedded fs)
	// We use a dummy target size.
	err := fm.LoadFont("testfont", "JetBrainsMono-Regular.ttf", 12, 24)
	if err != nil {
		t.Skipf("Skipping test because font loading failed (likely missing embedded font): %v", err)
	}

	var wg sync.WaitGroup
	// Spawn multiple goroutines to access GetOrCreateScaledFace
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Use different scales to trigger writes to the map
			scale := float64(i%5) + 1.0
			_, _ = fm.GetOrCreateScaledFace("testfont", scale, scale)
		}(i)
	}

	// Concurrently clear the cache to stress test locking
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fm.ClearScaledFaceCache()
		}
	}()

	wg.Wait()
}

// TestClearScaledFaceCachePreservesBaseFace verifies that ClearScaledFaceCache
// does not close the base font face.
func TestClearScaledFaceCachePreservesBaseFace(t *testing.T) {
	fm := NewFontManager()
	err := fm.LoadFont("testfont", "JetBrainsMono-Regular.ttf", 12, 24)
	if err != nil {
		t.Fatalf("LoadFont failed: %v", err)
	}

	// Generate a scaled face
	face2x, err := fm.GetOrCreateScaledFace("testfont", 2.0, 2.0)
	if err != nil {
		t.Fatalf("GetOrCreateScaledFace failed: %v", err)
	}
	if face2x == nil {
		t.Fatal("expected non-nil face2x")
	}

	// Clear cache
	fm.ClearScaledFaceCache()

	// Base font face must still be valid and usable
	sf, err := fm.GetFont("testfont")
	if err != nil {
		t.Fatalf("GetFont failed: %v", err)
	}
	if sf.face == nil {
		t.Fatal("expected non-nil sf.face")
	}

	// Calling Glyph on sf.face must succeed
	_, _, _, _, ok := sf.face.Glyph(fixed.Point26_6{}, 'A')
	if !ok {
		t.Fatal("expected glyph for 'A' to be found in base face")
	}

	// Re-fetching 1.0x1.0 scaled face should succeed
	face1x, err := fm.GetOrCreateScaledFace("testfont", 1.0, 1.0)
	if err != nil {
		t.Fatalf("GetOrCreateScaledFace(1.0, 1.0) failed: %v", err)
	}
	if face1x == nil {
		t.Fatal("expected non-nil face1x")
	}
}
