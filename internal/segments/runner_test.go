package segments

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/SuperCoolPencil/cue/internal/domain"
)

type inventoryClient struct {
	domain.LibraryClient
	domain.PlaybackClient
	fail    bool
	removed bool
}

func (c *inventoryClient) GetLibraries(context.Context) ([]domain.Library, error) {
	return []domain.Library{{ID: "tv", Type: "show"}}, nil
}
func (c *inventoryClient) GetShows(_ context.Context, _ string, offset, limit int) ([]*domain.Show, int, error) {
	if c.fail {
		return nil, 0, errors.New("offline")
	}
	return []*domain.Show{{ID: "kept"}}, 1, nil
}
func (c *inventoryClient) GetSeasons(context.Context, string) ([]*domain.Season, error) {
	return []*domain.Season{{ID: "season"}}, nil
}
func (c *inventoryClient) GetEpisodes(context.Context, string) ([]*domain.MediaItem, error) {
	out := []*domain.MediaItem{{ID: "a", ShowID: "kept"}, {ID: "b", ShowID: "kept"}, {ID: "c", ShowID: "kept"}}
	if c.removed {
		return out[:2], nil
	}
	return out, nil
}
func (c *inventoryClient) ResolvePlayable(_ context.Context, id string) (domain.PlayableMedia, error) {
	return domain.PlayableMedia{URL: id, SourceID: "source", Revision: "1", DurationMs: 200000}, nil
}
func TestStartupAnalysisCacheAndCleanup(t *testing.T) {
	client := &inventoryClient{}
	cache := Cache{Root: t.TempDir()}
	calls := 0
	analyzer := Analyzer{Client: client, Cache: cache, Server: "server", User: "user", Extractor: func(_ context.Context, _ string, _ int64, _ int, _ int) ([]uint32, error) {
		calls++
		return randomFingerprint(5, 500), nil
	}}
	obsolete := cache.ShowDir("server", "user", "removed")
	other := cache.ShowDir("other", "user", "removed")
	for _, dir := range []string{obsolete, other} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := analyzer.Startup(context.Background()); err != nil {
		t.Fatal(err)
	}
	if calls != 6 {
		t.Fatalf("extraction calls=%d", calls)
	}
	if _, err := os.Stat(obsolete); !os.IsNotExist(err) {
		t.Fatal("removed show retained")
	}
	if _, err := os.Stat(other); err != nil {
		t.Fatal("other server pruned")
	}
	if err := analyzer.Startup(context.Background()); err != nil {
		t.Fatal(err)
	}
	if calls != 6 {
		t.Fatal("unchanged media re-fingerprinted")
	}
	client.removed = true
	if err := analyzer.Startup(context.Background()); err != nil {
		t.Fatal(err)
	}
	id := Identity{Server: "server", User: "user", Show: "kept", Item: "c", Source: "source", Revision: "1", DurationMs: 200000}
	if _, ok := cache.Load(id); ok {
		t.Fatal("removed episode retained")
	}
}
func TestStartupFailureDoesNotPrune(t *testing.T) {
	cache := Cache{Root: t.TempDir()}
	obsolete := cache.ShowDir("server", "user", "removed")
	os.MkdirAll(obsolete, 0700)
	analyzer := Analyzer{Client: &inventoryClient{fail: true}, Cache: cache, Server: "server", User: "user"}
	if err := analyzer.Startup(context.Background()); err == nil {
		t.Fatal("expected inventory failure")
	}
	if _, err := os.Stat(obsolete); err != nil {
		t.Fatal("pruned on incomplete inventory")
	}
}
