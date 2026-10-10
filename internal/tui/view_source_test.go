package tui

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/SuperCoolPencil/cue/internal/playlist"
	"github.com/SuperCoolPencil/cue/internal/store"
	"github.com/SuperCoolPencil/cue/internal/tui/components"
)

func TestLocalViewsRefreshWithoutServerLibraryRequests(t *testing.T) {
	for _, id := range []string{recentLibraryID, queueLibraryID, filterUnwatchedID, filterInProgressID, filter4KID, filterHighlyRatedID} {
		for _, manual := range []bool{false, true} {
			t.Run(id+"/"+map[bool]string{true: "manual", false: "playback"}[manual], func(t *testing.T) {
				st, _ := store.NewLibraryStore("", "", "")
				movie := &domain.MediaItem{ID: "movie", Type: domain.MediaTypeMovie, ViewOffset: time.Minute, Width: 3840, Height: 2160, Rating: 9}
				if id == filterUnwatchedID {
					movie.ViewOffset = 0
				}
				st.SaveLibraries([]domain.Library{{ID: "movies", Type: "movie"}})
				st.SaveMovies("movies", []*domain.MediaItem{movie}, 1)
				st.SaveQueueItems([]*domain.MediaItem{movie})
				m := newVirtualTestModel(st) // No client: a server call would fail the test.
				m.Inspector = components.NewInspector()
				m.ColumnStack.Reset(components.NewLibraryColumn(m.allLibraryEntries()))
				col := components.NewListColumn(components.ColumnTypeMixed, "Local view")
				col.SetContentID(id)
				m.ColumnStack.Push(col, 0)
				if manual {
					m, _ = m.refreshLibraryContent(col)
				} else {
					m.refreshCurrentView()
				}
				if col.ItemCount() != 1 || col.SelectedItemID() != movie.ID || col.IsRefreshing() {
					t.Fatalf("local projection not refreshed: count=%d refreshing=%v", col.ItemCount(), col.IsRefreshing())
				}
				if _, ok := st.GetMovies("movies"); !ok {
					t.Fatal("refresh erased local projection input")
				}
			})
		}
	}
}

type viewPlaylistClient struct {
	domain.PlaylistClient
	requested string
}

func (c *viewPlaylistClient) GetPlaylistItems(_ context.Context, id string) ([]*domain.MediaItem, error) {
	c.requested = id
	return []*domain.MediaItem{{ID: "movie", Type: domain.MediaTypeMovie}}, nil
}

func TestPlaylistPlaybackRefreshKeepsPlaylistView(t *testing.T) {
	st, _ := store.NewLibraryStore("", "", "")
	m := newVirtualTestModel(st)
	client := &viewPlaylistClient{}
	m.PlaylistService = playlist.NewService(client, st, slog.Default())
	m.ColumnStack.Reset(components.NewLibraryColumn(m.allLibraryEntries()))
	col := components.NewListColumn(components.ColumnTypePlaylistItems, "Playlist")
	col.SetContentID("playlist-1")
	m.ColumnStack.Push(col, 0)
	cmd := m.refreshCurrentView()
	if cmd == nil {
		t.Fatal("no refresh")
	}
	if _, ok := cmd().(PlaylistItemsLoadedMsg); !ok || client.requested != "playlist-1" {
		t.Fatalf("wrong playlist refresh: %q", client.requested)
	}
	if m.ColumnStack.Top() != col {
		t.Fatal("playlist view replaced during refresh")
	}
}
