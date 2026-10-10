package tui

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/SuperCoolPencil/cue/internal/config"
	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/SuperCoolPencil/cue/internal/library"
	"github.com/SuperCoolPencil/cue/internal/playlist"
	"github.com/SuperCoolPencil/cue/internal/store"
	"github.com/SuperCoolPencil/cue/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

type continueRefreshClient struct {
	domain.LibraryClient
	calls int
	items []*domain.MediaItem
}

func (c *continueRefreshClient) GetContinueWatching(context.Context) ([]*domain.MediaItem, error) {
	c.calls++
	return c.items, nil
}

func TestContinueWatchingRefreshUsesServerContinueEndpoint(t *testing.T) {
	for _, manual := range []bool{false, true} {
		name := "after playback"
		if manual {
			name = "manual"
		}
		t.Run(name, func(t *testing.T) {
			st, _ := store.NewLibraryStore("", "", "")
			m := newVirtualTestModel(st)
			movie := &domain.MediaItem{ID: "movie", Title: "Horrible Bosses", Type: domain.MediaTypeMovie}
			client := &continueRefreshClient{items: []*domain.MediaItem{movie}}
			m.LibraryService = library.NewService(client, st, slog.Default())
			m.Inspector = components.NewInspector()
			m.ColumnStack.Reset(components.NewLibraryColumn(m.allLibraryEntries()))
			col := components.NewListColumn(components.ColumnTypeMixed, "Continue Watching")
			col.SetContentID(continueLibraryID)
			col.SetItems([]*domain.MediaItem{{ID: "other", Type: domain.MediaTypeEpisode}, movie})
			col.SetSelectedByID(movie.ID)
			m.ColumnStack.Push(col, 0)
			var cmd tea.Cmd
			if manual {
				m, cmd = m.refreshLibraryContent(col)
			} else {
				updated, next := m.Update(PlaybackFinishedMsg{Title: movie.Title})
				m, cmd = updated.(Model), next
			}
			if cmd == nil {
				t.Fatal("missing refresh command")
			}
			msg := cmd()
			if batch, ok := msg.(tea.BatchMsg); ok {
				if len(batch) != 1 {
					t.Fatalf("unexpected refresh batch size: %d", len(batch))
				}
				msg = batch[0]()
			}
			result, owned := msg.(viewResultMsg)
			if _, ok := result.message.(ContinueWatchingLoadedMsg); !owned || !ok || client.calls != 1 {
				t.Fatalf("wrong refresh endpoint: message %T, calls %d", msg, client.calls)
			}
			updated, _ := m.Update(msg)
			m = updated.(Model)
			if m.Loading || col.IsRefreshing() || col.SelectedItemID() != movie.ID {
				t.Fatalf("refresh did not finish with selected movie: loading=%v selection=%s", m.Loading, col.SelectedItemID())
			}
		})
	}
}

func TestVirtualLibraryEntries(t *testing.T) {
	entries := virtualLibraryEntries()
	if len(entries) != 7 {
		t.Fatalf("virtual entries = %d", len(entries))
	}
	if entries[0].ID != continueLibraryID || entries[0].Name != "Continue Watching" {
		t.Fatalf("first virtual entry = %#v", entries[0])
	}
}

func TestContinueWatchingAvailableWithoutCachedProgress(t *testing.T) {
	st, _ := store.NewLibraryStore("", "", "")
	m := newVirtualTestModel(st)
	entries := m.allLibraryEntries()
	if len(entries) == 0 || entries[0].ID != continueLibraryID {
		t.Fatalf("Continue Watching missing with empty cache: %#v", entries)
	}
}

func TestDrillVirtualContinueWatching(t *testing.T) {
	st, _ := store.NewLibraryStore("", "", "")
	libs := []domain.Library{{ID: "movies", Name: "Movies", Type: "movie"}}
	_ = st.SaveLibraries(libs)
	_ = st.SaveMovies("movies", []*domain.MediaItem{{ID: "m1", Title: "Movie", ViewOffset: time.Minute}}, 1)

	model := newVirtualTestModel(st)
	root := components.NewLibraryColumn(model.allLibraryEntries())
	model.ColumnStack.Reset(root)

	result := model.drillVirtualLibrary(domain.Library{ID: continueLibraryID, Name: "Continue Watching", Type: "cue"}, 0)
	if result == nil || model.ColumnStack.Len() != 2 || result.Cmd == nil {
		t.Fatalf("virtual drill failed: %#v len=%d", result, model.ColumnStack.Len())
	}
	top := model.ColumnStack.Top()
	if top.ContentID() != continueLibraryID {
		t.Fatalf("content id = %q", top.ContentID())
	}

	// Verify that the column is in a loading state
	if !model.Loading {
		t.Fatal("expected model to be in loading state")
	}

	// Simulate the command completion
	msg := ContinueWatchingLoadedMsg{Items: []*domain.MediaItem{{ID: "m1", Title: "Movie"}}}
	m, _ := model.Update(msg)
	model = m.(Model)

	if top := model.ColumnStack.Top(); top.ItemCount() != 1 {
		t.Fatalf("item count = %d", top.ItemCount())
	}
}

func newVirtualTestModel(st domain.Store) Model {
	libSvc := library.NewService(nil, st, slog.Default())
	playlistSvc := playlist.NewService(nil, st, slog.Default())
	return Model{
		State:           StateBrowsing,
		Store:           st,
		LibraryService:  libSvc,
		PlaylistService: playlistSvc,
		ColumnStack:     NewColumnStack(),
		LibraryStates:   make(map[string]components.LibrarySyncState),
		AppConfig:       config.DefaultConfig(),
	}
}
