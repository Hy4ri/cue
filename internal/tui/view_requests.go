package tui

import (
	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/SuperCoolPencil/cue/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

type viewRequestKey struct {
	owner    *components.ListColumn
	resource string
}

type viewResultMsg struct {
	key         viewRequestKey
	generation  uint64
	contentID   string
	selectionID string
	message     tea.Msg
}

func (m *Model) requestView(col *components.ListColumn, resource string, cmd tea.Cmd) tea.Cmd {
	if col == nil || cmd == nil {
		return nil
	}
	if m.viewRequests == nil {
		m.viewRequests = make(map[viewRequestKey]uint64)
	}
	m.requestGeneration++
	key := viewRequestKey{col, resource}
	m.viewRequests[key] = m.requestGeneration
	result := viewResultMsg{key: key, generation: m.requestGeneration, contentID: col.ContentID()}
	if resource == "playback" || resource == "modal" {
		result.selectionID = col.SelectedItemID()
	}
	m.Loading = true
	if col.IsEmpty() {
		col.SetLoading(true)
	} else {
		col.SetRefreshing(true)
	}
	return func() tea.Msg {
		result.message = cmd()
		return result
	}
}

func (m Model) ownsColumn(col *components.ListColumn) bool {
	if m.ColumnStack == nil {
		return false
	}
	for i := 0; i < m.ColumnStack.Len(); i++ {
		if m.ColumnStack.Get(i) == col {
			return true
		}
	}
	return false
}

func (m Model) acceptsViewResult(result viewResultMsg) bool {
	if m.viewRequests[result.key] != result.generation || !m.ownsColumn(result.key.owner) || result.key.owner.ContentID() != result.contentID {
		return false
	}
	if result.key.resource == "playback" || result.key.resource == "modal" {
		return m.State == StateBrowsing && m.ColumnStack.Top() == result.key.owner && result.key.owner.SelectedItemID() == result.selectionID
	}
	return true
}

func (m *Model) cancelPlaybackRequests() {
	for key := range m.viewRequests {
		if key.resource == "playback" || key.resource == "modal" {
			delete(m.viewRequests, key)
			m.finishColumnRequest(key.owner)
		}
	}
	m.Loading = m.SyncingCount > 0 || m.hasActiveViewRequests()
}

func (m *Model) cancelViewRequests() {
	for key := range m.viewRequests {
		delete(m.viewRequests, key)
		m.finishColumnRequest(key.owner)
	}
}

func (m *Model) finishColumnRequest(col *components.ListColumn) {
	for key := range m.viewRequests {
		if key.owner == col {
			return
		}
	}
	col.SetLoading(false)
	col.SetRefreshing(false)
}

func (m Model) hasActiveViewRequests() bool {
	for key := range m.viewRequests {
		if m.ownsColumn(key.owner) {
			return true
		}
	}
	return false
}

func (m *Model) requestContent(id, resource string, cmd tea.Cmd) tea.Cmd {
	return m.requestView(m.ColumnStack.FindColumn(id), resource, cmd)
}

func (m *Model) loadMovies(lib domain.Library) tea.Cmd {
	return m.requestContent(lib.ID, "content", LoadMoviesCmd(m.LibraryService, lib))
}
func (m *Model) loadShows(lib domain.Library) tea.Cmd {
	return m.requestContent(lib.ID, "content", LoadShowsCmd(m.LibraryService, lib))
}
func (m *Model) loadMixed(lib domain.Library) tea.Cmd {
	return m.requestContent(lib.ID, "content", LoadMixedLibraryCmd(m.LibraryService, lib))
}
func (m *Model) loadSeasons(libID, showID string) tea.Cmd {
	col := m.ColumnStack.FindColumn(showID)
	// Reloading season groups supersedes their outstanding episode loads.
	for key := range m.viewRequests {
		if key.owner == col {
			delete(m.viewRequests, key)
		}
	}
	return m.requestView(col, "content", LoadSeasonsCmd(m.LibraryService, libID, showID))
}
func (m *Model) loadEpisodes(libID, showID, seasonID string) tea.Cmd {
	col := m.ColumnStack.FindColumn(seasonID)
	if col == nil {
		col = m.ColumnStack.FindColumn(showID)
	}
	return m.requestView(col, "episodes:"+seasonID, LoadEpisodesCmd(m.LibraryService, libID, showID, seasonID))
}
func (m *Model) loadPlaylists() tea.Cmd {
	return m.requestContent(playlistsLibraryID, "content", LoadPlaylistsCmd(m.PlaylistService))
}
func (m *Model) loadPlaylistItems(id string) tea.Cmd {
	return m.requestContent(id, "content", LoadPlaylistItemsCmd(m.PlaylistService, id))
}
func (m *Model) loadContinueWatching() tea.Cmd {
	return m.requestContent(continueLibraryID, "content", LoadContinueWatchingCmd(m.LibraryService))
}
func (m *Model) loadLibraries(refresh bool) tea.Cmd {
	return m.requestView(m.ColumnStack.Get(0), "root", loadLibrariesCmd(m.LibraryService, refresh))
}
