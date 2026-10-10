package tui

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/SuperCoolPencil/cue/internal/library"
	"github.com/SuperCoolPencil/cue/internal/store"
	"github.com/SuperCoolPencil/cue/internal/tui/components"
)

type playbackSeasonClient struct {
	domain.LibraryClient
	items []*domain.MediaItem
	calls int
}

func (c *playbackSeasonClient) GetEpisodes(context.Context, string) ([]*domain.MediaItem, error) {
	c.calls++
	return c.items, nil
}

func TestPlaybackContextUsesExplicitQueueSource(t *testing.T) {
	for _, id := range []string{queueLibraryID, "playlist-1", continueLibraryID} {
		for _, autoplay := range []bool{false, true} {
			t.Run(id+"/"+map[bool]string{false: "single", true: "autoplay"}[autoplay], func(t *testing.T) {
				m, _ := requestTestModel(t)
				m.UIConfig.Autoplay = autoplay
				kind := components.ColumnTypeMixed
				if id == "playlist-1" {
					kind = components.ColumnTypePlaylistItems
				}
				col := components.NewListColumn(kind, "Items")
				col.SetContentID(id)
				movie := &domain.MediaItem{ID: "movie", Type: domain.MediaTypeMovie, ViewOffset: time.Minute}
				col.SetItems([]*domain.MediaItem{movie, {ID: "episode", Type: domain.MediaTypeEpisode}})
				col.SetSelectedByID(movie.ID)
				m.ColumnStack.Push(col, 0)
				updated, _ := m.playItemWithContext(col, movie, true)
				m = updated.(Model)
				want := 0
				if autoplay && id != continueLibraryID {
					want = 2
				}
				if m.pendingPlayback == nil || m.pendingPlayback.ID != movie.ID || len(m.pendingPlaylist) != want {
					t.Fatalf("wrong queue for source: %#v", m.pendingPlaylist)
				}
				movie.ID = "mutated"
				if m.pendingPlayback.ID != "movie" {
					t.Fatal("confirmation did not snapshot selection")
				}
			})
		}
	}
}

func TestFlatEpisodeLoadsFullSeasonOnlyForAutoplay(t *testing.T) {
	for _, autoplay := range []bool{false, true} {
		t.Run(map[bool]string{false: "single", true: "autoplay"}[autoplay], func(t *testing.T) {
			m, _ := requestTestModel(t)
			m.UIConfig.Autoplay = autoplay
			st, _ := store.NewLibraryStore("", "", "")
			item := &domain.MediaItem{ID: "ep2", Type: domain.MediaTypeEpisode, ShowID: "show", ParentID: "season", LibraryID: "lib", ViewOffset: time.Minute}
			client := &playbackSeasonClient{items: []*domain.MediaItem{{ID: "ep1", Type: domain.MediaTypeEpisode}, item, {ID: "ep3", Type: domain.MediaTypeEpisode}}}
			m.LibraryService = library.NewService(client, st, slog.Default())
			col := components.NewListColumn(components.ColumnTypeMixed, "Continue Watching")
			col.SetContentID(continueLibraryID)
			col.SetItems([]*domain.MediaItem{item})
			m.ColumnStack.Push(col, 0)
			updated, cmd := m.playItemWithContext(col, item, true)
			m = updated.(Model)
			if autoplay {
				if cmd == nil {
					t.Fatal("missing season lookup")
				}
				updated, _ = m.Update(cmd())
				m = updated.(Model)
				if client.calls != 1 || len(m.pendingPlaylist) != 3 {
					t.Fatal("full season was not used")
				}
			} else if client.calls != 0 || cmd != nil || len(m.pendingPlaylist) != 0 {
				t.Fatal("single playback fetched a season")
			}
			if m.State != StateConfirmResume || m.pendingPlayback.ID != item.ID {
				t.Fatal("wrong resume selection")
			}
		})
	}
}
