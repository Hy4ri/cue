package tui

import (
	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/SuperCoolPencil/cue/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

// A source identifies the data behind a view independently of its presentation.
// In particular, a mixed column need not represent a server library.
type viewKind uint8

const (
	viewNone viewKind = iota
	viewLibrary
	viewContinue
	viewRecent
	viewQueue
	viewFilter
	viewSeasons
	viewEpisodes
	viewPlaylists
	viewPlaylistItems
)

type viewSource struct {
	kind          viewKind
	id            string
	library       *domain.Library
	libID, showID string
}

func (m Model) sourceFor(col *components.ListColumn) viewSource {
	if col == nil {
		return viewSource{}
	}
	source := viewSource{id: col.ContentID()}
	switch source.id {
	case continueLibraryID:
		source.kind = viewContinue
		return source
	case recentLibraryID:
		source.kind = viewRecent
		return source
	case queueLibraryID:
		source.kind = viewQueue
		return source
	}
	if filterKey(source.id) != "" {
		source.kind = viewFilter
		return source
	}
	switch col.ColumnType() {
	case components.ColumnTypeMovies, components.ColumnTypeShows, components.ColumnTypeMixed:
		source.library = m.findLibrary(source.id)
		if source.library != nil {
			source.kind = viewLibrary
		}
	case components.ColumnTypeSeasons, components.ColumnTypeSeasonEpisodes:
		source.kind, source.libID, source.showID = viewSeasons, m.currentLibID, source.id
	case components.ColumnTypeEpisodes:
		source.kind, source.libID, source.showID = viewEpisodes, m.currentLibID, m.currentShowID
		if item := col.SelectedMediaItem(); item != nil {
			if item.LibraryID != "" {
				source.libID = item.LibraryID
			}
			if item.ShowID != "" {
				source.showID = item.ShowID
			}
		}
	case components.ColumnTypePlaylists:
		source.kind = viewPlaylists
	case components.ColumnTypePlaylistItems:
		source.kind = viewPlaylistItems
	}
	return source
}

// reloadView is shared by manual refresh and playback completion. Only server
// sources invalidate server caches; local projections retain their inputs.
func (m *Model) reloadView(col *components.ListColumn) tea.Cmd {
	source := m.sourceFor(col)
	var cmd tea.Cmd
	switch source.kind {
	case viewContinue:
		cmd = LoadContinueWatchingCmd(m.LibraryService)
	case viewRecent:
		col.ReplaceItems(m.LibraryService.RecentlyAdded(0))
	case viewQueue:
		col.ReplaceItems(m.PlaylistService.QueueItems())
	case viewFilter:
		col.ReplaceItems(m.LibraryService.SmartFiltered(filterKey(source.id), 0))
	case viewLibrary:
		m.LibraryService.InvalidateLibrary(source.library.ID)
		switch source.library.Type {
		case "movie":
			cmd = LoadMoviesCmd(m.LibraryService, *source.library)
		case "show":
			cmd = LoadShowsCmd(m.LibraryService, *source.library)
		default:
			cmd = LoadMixedLibraryCmd(m.LibraryService, *source.library)
		}
	case viewSeasons:
		m.LibraryService.InvalidateShow(source.libID, source.showID)
		cmd = LoadSeasonsCmd(m.LibraryService, source.libID, source.showID)
	case viewEpisodes:
		m.LibraryService.InvalidateSeason(source.libID, source.showID, source.id)
		cmd = LoadEpisodesCmd(m.LibraryService, source.libID, source.showID, source.id)
	case viewPlaylists:
		m.PlaylistService.InvalidatePlaylists()
		cmd = LoadPlaylistsCmd(m.PlaylistService)
	case viewPlaylistItems:
		m.PlaylistService.InvalidatePlaylistItems(source.id)
		cmd = LoadPlaylistItemsCmd(m.PlaylistService, source.id)
	default:
		return nil
	}
	if cmd != nil {
		col.SetRefreshing(true)
		m.Loading = true
		return cmd
	}
	m.Loading = m.SyncingCount > 0
	return m.updateInspector()
}

// Once a full sync finishes, project its new data into the preserved view.
func (m *Model) refreshSyncedView() tea.Cmd {
	root := m.ColumnStack.Get(0)
	if root == nil {
		return nil
	}
	root.ReplaceItems(m.allLibraryEntries())
	col := m.ColumnStack.Top()
	if col == root {
		return nil
	}
	source := m.sourceFor(col)
	if source.kind == viewLibrary {
		switch source.library.Type {
		case "movie":
			if items, ok := m.Store.GetMovies(source.id); ok {
				col.ReplaceItems(items)
			}
		case "show":
			if items, ok := m.Store.GetShows(source.id); ok {
				col.ReplaceItems(items)
			}
		default:
			if items, ok := m.Store.GetMixedContent(source.id); ok {
				col.ReplaceItems(items)
			}
		}
		return m.updateInspector()
	}
	return m.reloadView(col)
}
