package segments

import (
	"context"
	"math/rand"
	"testing"
)

func randomFingerprint(seed int64, n int) []uint32 {
	r := rand.New(rand.NewSource(seed))
	fp := make([]uint32, n)
	for i := range fp {
		fp[i] = r.Uint32()
	}
	return fp
}
func TestDetectShiftedRecurringSequences(t *testing.T) {
	theme := randomFingerprint(5, 500)
	var episodes [][]uint32
	for i := 0; i < 3; i++ {
		fp := randomFingerprint(int64(i+20), 1000)
		copy(fp[50+i*80:], theme)
		episodes = append(episodes, fp)
	}
	got, err := Detect(context.Background(), episodes, []int64{200000, 200000, 200000}, []int64{0, 0, 0}, "intro")
	if err != nil {
		t.Fatal(err)
	}
	for i, segments := range got {
		if len(segments) != 1 || !segments[0].ManualOnly {
			t.Fatalf("episode %d: %+v", i, segments)
		}
		want := int64((float64(50+i*80)*FrameSeconds + boundaryMargin) * 1000)
		if segments[0].StartMs != want {
			t.Fatalf("episode %d offset: %+v want %d", i, segments, want)
		}
	}
}
func TestDetectAbstainsForSilenceAndInsufficientSupport(t *testing.T) {
	for _, episodes := range [][][]uint32{{make([]uint32, 600), make([]uint32, 600), make([]uint32, 600)}, {randomFingerprint(1, 600), randomFingerprint(1, 600), randomFingerprint(2, 600)}} {
		got, err := Detect(context.Background(), episodes, []int64{200000, 200000, 200000}, []int64{0, 0, 0}, "intro")
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range got {
			if len(s) != 0 {
				t.Fatalf("false match: %+v", s)
			}
		}
	}
}
func TestDetectCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Detect(ctx, [][]uint32{randomFingerprint(1, 600)}, []int64{200000}, []int64{0}, "intro")
	if err == nil {
		t.Fatal("ignored cancellation")
	}
}
func TestCacheRevision(t *testing.T) {
	c := Cache{Root: t.TempDir()}
	id := Identity{Server: "server", Item: "episode", Source: "version", Revision: "old", DurationMs: 100000}
	entry := Cached{Identity: id, Version: Version, Head: []uint32{1, 2}}
	if err := c.Save(entry); err != nil {
		t.Fatal(err)
	}
	if _, ok := c.Load(id); !ok {
		t.Fatal("cache miss")
	}
	id.Revision = "new"
	if _, ok := c.Load(id); ok {
		t.Fatal("stale revision reused")
	}
}

func TestDetectToleratesBriefFingerprintDifferences(t *testing.T) {
	theme := randomFingerprint(5, 500)
	episodes := [][]uint32{append([]uint32(nil), theme...), append([]uint32(nil), theme...), append([]uint32(nil), theme...)}
	// Interrupt every twelve seconds: no uninterrupted passage reaches 25s.
	for start := 90; start < 450; start += 90 {
		for n := start; n < start+5; n++ {
			episodes[0][n] ^= 0xffffffff
		}
	}
	got, err := Detect(context.Background(), episodes, []int64{200000, 200000, 200000}, []int64{0, 0, 0}, "intro")
	if err != nil {
		t.Fatal(err)
	}
	for i, s := range got {
		if len(s) != 1 || s[0].EndMs < 60000 {
			t.Fatalf("episode %d: %+v", i, s)
		}
	}
}

func TestDetectRejectsMostlyMatchingShortFragments(t *testing.T) {
	theme := randomFingerprint(8, 600)
	altered := append([]uint32(nil), theme...)
	for i := 0; i < len(altered); i++ {
		if i%10 < 2 {
			altered[i] ^= 0xffffffff
		}
	}
	got, err := Detect(context.Background(), [][]uint32{altered, theme, theme}, []int64{200000, 200000, 200000}, []int64{0, 0, 0}, "intro")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range got {
		if len(s) > 0 {
			t.Fatalf("accepted 20%% mismatch: %+v", s)
		}
	}
}
