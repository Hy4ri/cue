package tui

import (
	"errors"
	"testing"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/SuperCoolPencil/cue/internal/store"
	"github.com/SuperCoolPencil/cue/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

func requestTestModel(t *testing.T) (Model, *components.ListColumn) {
	t.Helper()
	st, _ := store.NewLibraryStore("", "", "")
	m := newVirtualTestModel(st)
	m.Inspector = components.NewInspector()
	m.ColumnStack.Reset(components.NewLibraryColumn(m.allLibraryEntries()))
	col := components.NewListColumn(components.ColumnTypePlaylists, "Playlists")
	col.SetContentID(playlistsLibraryID)
	m.ColumnStack.Push(col, 0)
	return m, col
}

func TestOlderViewResultCannotReplaceNewerRequest(t *testing.T) {
	m, col := requestTestModel(t)
	old := m.requestView(col, "content", func() tea.Msg { return PlaylistsLoadedMsg{Playlists: []*domain.Playlist{{ID: "old"}}} })
	newer := m.requestView(col, "content", func() tea.Msg { return PlaylistsLoadedMsg{Playlists: []*domain.Playlist{{ID: "new"}}} })
	updated, _ := m.Update(old())
	m = updated.(Model)
	if col.ItemCount() != 0 || !m.Loading {
		t.Fatal("old response overwrote the pending request")
	}
	updated, _ = m.Update(newer())
	m = updated.(Model)
	if col.SelectedItemID() != "new" || m.Loading || col.IsLoading() {
		t.Fatal("new response was not applied")
	}
}

func TestViewResponseCannotPopulateReopenedColumn(t *testing.T) {
	m, col := requestTestModel(t)
	cmd := m.requestView(col, "content", func() tea.Msg { return PlaylistsLoadedMsg{Playlists: []*domain.Playlist{{ID: "stale"}}} })
	m.ColumnStack.Pop()
	reopened := components.NewListColumn(components.ColumnTypePlaylists, "Playlists")
	reopened.SetContentID(playlistsLibraryID)
	m.ColumnStack.Push(reopened, 0)
	updated, _ := m.Update(cmd())
	m = updated.(Model)
	if reopened.ItemCount() != 0 || m.Loading {
		t.Fatal("old request populated a different instance of the view")
	}
}

func TestParentViewResponseUpdatesItsOwner(t *testing.T) {
	m, col := requestTestModel(t)
	cmd := m.requestView(col, "content", func() tea.Msg { return PlaylistsLoadedMsg{Playlists: []*domain.Playlist{{ID: "playlist"}}} })
	child := components.NewListColumn(components.ColumnTypePlaylistItems, "Items")
	child.SetContentID("playlist")
	m.ColumnStack.Push(child, 0)
	updated, _ := m.Update(cmd())
	m = updated.(Model)
	if col.SelectedItemID() != "playlist" || child.ItemCount() != 0 || m.ColumnStack.Top() != child {
		t.Fatal("response did not stay with its owning column")
	}
}

func TestStaleErrorCannotClearNewRequestLoading(t *testing.T) {
	m, col := requestTestModel(t)
	old := m.requestView(col, "content", func() tea.Msg { return ErrMsg{Err: errors.New("old failure")} })
	m.requestView(col, "content", func() tea.Msg { return PlaylistsLoadedMsg{} })
	updated, _ := m.Update(old())
	m = updated.(Model)
	if m.StatusIsErr || !m.Loading || !col.IsLoading() {
		t.Fatal("stale error changed the active request")
	}
}

func TestPlaybackLookupCancelledByNavigation(t *testing.T) {
	m, _ := requestTestModel(t)
	col := components.NewListColumn(components.ColumnTypeMixed, "Continue Watching")
	col.SetContentID(continueLibraryID)
	item := &domain.MediaItem{ID: "episode", Type: domain.MediaTypeEpisode, ViewOffset: time.Minute}
	col.SetItems([]*domain.MediaItem{item})
	m.ColumnStack.Push(col, 0)
	cmd := m.requestView(col, "playback", func() tea.Msg { return SeasonForPlaybackLoadedMsg{Item: item, Resume: true} })
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	// Even returning to the same selection must not resurrect the old intent.
	m.ColumnStack.Push(col, 0)
	updated, next := m.Update(cmd())
	m = updated.(Model)
	if m.State != StateBrowsing || m.pendingPlayback != nil || next != nil {
		t.Fatal("cancelled lookup opened confirmation or playback")
	}
}

func TestUnscopedPlaylistResponseCannotPopulateMovieView(t *testing.T) {
	m, _ := requestTestModel(t)
	col := components.NewListColumn(components.ColumnTypeMovies, "Movies")
	col.SetItems([]*domain.MediaItem{{ID: "movie", Type: domain.MediaTypeMovie}})
	m.ColumnStack.Push(col, 0)
	m.Update(PlaylistsLoadedMsg{Playlists: []*domain.Playlist{{ID: "playlist"}}})
	if col.SelectedItemID() != "movie" {
		t.Fatal("playlist load replaced movie view")
	}
}
