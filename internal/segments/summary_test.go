package segments

import (
	"github.com/SuperCoolPencil/cue/internal/domain"
	"os"
	"testing"
	"time"
)

func TestShowIntroSummary(t *testing.T) {
	c := Cache{Root: t.TempDir()}
	id := Identity{Server: "server", User: "user", Show: "show", Item: "one", DurationMs: 100000}
	save := func(e Cached) {
		t.Helper()
		if err := c.Save(e); err != nil {
			t.Fatal(err)
		}
	}
	old := Cached{Identity: id, Version: Version, Membership: "season", Segments: []domain.SkipSegment{{Kind: "intro", StartMs: 1000, EndMs: 30000}}}
	save(old)
	past := time.Now().Add(-time.Hour)
	if err := os.Chtimes(c.path(id), past, past); err != nil {
		t.Fatal(err)
	}
	id.Revision = "new"
	save(Cached{Identity: id, Version: Version, Membership: "season"})
	id.Item = "two"
	save(Cached{Identity: id, Version: Version, Membership: "season", Segments: old.Segments})
	id.Item = "three"
	save(Cached{Identity: id, Version: Version})
	id.Item = "four"
	save(Cached{Identity: id, Version: "obsolete", Membership: "season", Segments: old.Segments})
	got, err := c.ShowIntroSummary("server", "user", "show")
	if err != nil {
		t.Fatal(err)
	}
	if got != (IntroSummary{Detected: 1, Analyzed: 2, Pending: 1}) {
		t.Fatalf("summary = %+v", got)
	}
	other, err := c.ShowIntroSummary("server", "another-user", "show")
	if err != nil || other != (IntroSummary{}) {
		t.Fatalf("other profile: %+v, %v", other, err)
	}
}
