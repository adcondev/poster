package emulator

import (
	"sync"
	"testing"
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
			scale := float64(i) + 1.0
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
