# Graph Report - .  (2026-05-28)

## Corpus Check
- 89 files · ~62,699 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 940 nodes · 1450 edges · 64 communities (31 shown, 33 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 199 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_ListColumn Component Core|ListColumn Component Core]]
- [[_COMMUNITY_TUI App Model|TUI App Model]]
- [[_COMMUNITY_UI Components Suite|UI Components Suite]]
- [[_COMMUNITY_System Bootstrap|System Bootstrap]]
- [[_COMMUNITY_Column Stack Tests|Column Stack Tests]]
- [[_COMMUNITY_Player and Scrobbler|Player and Scrobbler]]
- [[_COMMUNITY_Jellyfin Client API|Jellyfin Client API]]
- [[_COMMUNITY_Jellyfin Metadata Mapping|Jellyfin Metadata Mapping]]
- [[_COMMUNITY_TUI Message Types|TUI Message Types]]
- [[_COMMUNITY_Library Storage|Library Storage]]
- [[_COMMUNITY_Plex Client API|Plex Client API]]
- [[_COMMUNITY_Plex Metadata Mapping|Plex Metadata Mapping]]
- [[_COMMUNITY_Fuzzy Search Engine|Fuzzy Search Engine]]
- [[_COMMUNITY_MediaItem Domain Entity|MediaItem Domain Entity]]
- [[_COMMUNITY_Library Service Layer|Library Service Layer]]
- [[_COMMUNITY_Integration Tests|Integration Tests]]
- [[_COMMUNITY_GlobalSearch Component|GlobalSearch Component]]
- [[_COMMUNITY_Season Group UI|Season Group UI]]
- [[_COMMUNITY_Sort Modal Component|Sort Modal Component]]
- [[_COMMUNITY_Project Metadata|Project Metadata]]
- [[_COMMUNITY_Season Domain Entity|Season Domain Entity]]
- [[_COMMUNITY_Plex API DTOs|Plex API DTOs]]
- [[_COMMUNITY_Library Domain Entity|Library Domain Entity]]
- [[_COMMUNITY_Playlist Domain Entity|Playlist Domain Entity]]
- [[_COMMUNITY_Show Domain Entity|Show Domain Entity]]
- [[_COMMUNITY_Jellyfin API DTOs|Jellyfin API DTOs]]
- [[_COMMUNITY_Playlist Service|Playlist Service]]
- [[_COMMUNITY_ColumnStack Navigation|ColumnStack Navigation]]
- [[_COMMUNITY_InputModal Component|InputModal Component]]
- [[_COMMUNITY_Mock Library Client|Mock Library Client]]
- [[_COMMUNITY_View Rendering|View Rendering]]
- [[_COMMUNITY_Plex Auth Flow|Plex Auth Flow]]
- [[_COMMUNITY_Keymap Bindings|Keymap Bindings]]
- [[_COMMUNITY_Store Cache Wrapper|Store Cache Wrapper]]
- [[_COMMUNITY_Auth Adapter Layer|Auth Adapter Layer]]
- [[_COMMUNITY_Watch Status Domain|Watch Status Domain]]
- [[_COMMUNITY_Local View Queries|Local View Queries]]
- [[_COMMUNITY_MPV IPC Connection|MPV IPC Connection]]
- [[_COMMUNITY_Virtual Model Data|Virtual Model Data]]
- [[_COMMUNITY_Queue Operations|Queue Operations]]
- [[_COMMUNITY_Playback Interfaces|Playback Interfaces]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 46|Community 46]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 49|Community 49]]
- [[_COMMUNITY_Community 50|Community 50]]
- [[_COMMUNITY_Community 51|Community 51]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 53|Community 53]]
- [[_COMMUNITY_Community 54|Community 54]]
- [[_COMMUNITY_Community 55|Community 55]]
- [[_COMMUNITY_Community 56|Community 56]]
- [[_COMMUNITY_Community 57|Community 57]]
- [[_COMMUNITY_Community 58|Community 58]]
- [[_COMMUNITY_Community 59|Community 59]]

## God Nodes (most connected - your core abstractions)
1. `ListColumn` - 71 edges
2. `LibraryStore` - 33 edges
3. `Model` - 31 edges
4. `Client` - 25 edges
5. `Client` - 23 edges
6. `MediaItem` - 21 edges
7. `Truncate()` - 20 edges
8. `NewLibraryStore()` - 18 edges
9. `Service` - 17 edges
10. `NewListColumn()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `TestCLIHelp()` --calls--> `runCLI()`  [INFERRED]
  main_test.go → main.go
- `run()` --calls--> `NullLogger()`  [INFERRED]
  main.go → internal/log/logger.go
- `run()` --calls--> `NewLauncher()`  [INFERRED]
  main.go → internal/player/launcher.go
- `run()` --calls--> `NewLibraryStore()`  [INFERRED]
  main.go → internal/store/store.go
- `run()` --calls--> `NewModel()`  [INFERRED]
  main.go → internal/tui/app.go

## Hyperedges (group relationships)
- **Media Server Backend Support** — cue_project, plex_backend, jellyfin_backend [EXTRACTED 1.00]
- **Release Pipeline** — release_yml, goreleaser_yaml, goreleaser_tool [EXTRACTED 1.00]
- **Terminal UI Layout Components** — screenshot1_png, three_column_layout, watch_status_indicators, metadata_info_panel [EXTRACTED 1.00]

## Communities (64 total, 33 thin omitted)

### Community 0 - "ListColumn Component Core"
Cohesion: 0.05
Nodes (14): appendSortTag(), formatMonthYear(), mediaItemWatchIndicator(), watchIndicator(), WrapEpisodes(), WrapLibraries(), WrapMovies(), WrapPlaylistItems() (+6 more)

### Community 1 - "TUI App Model"
Cohesion: 0.05
Nodes (38): Model, Model, playlistsLibraryEntry(), virtualLibraryEntries(), ApplicationState, AddToPlaylistCmd(), AddToQueueCmd(), ClearLibraryStatusCmd() (+30 more)

### Community 2 - "UI Components Suite"
Cohesion: 0.05
Nodes (21): NewGlobalSearch(), TestGlobalSearchSelectionAndQueryChange(), Inspector, formatRemainingDuration(), NewInspector(), renderMediaBody(), renderMediaFooter(), renderMediaHeader() (+13 more)

### Community 3 - "System Bootstrap"
Cohesion: 0.06
Nodes (36): ClearCache(), ClearServerConfig(), DefaultCachePath(), DefaultConfig(), defaultConfigPath(), defaultLogPath(), LoadConfig(), SaveConfig() (+28 more)

### Community 4 - "Column Stack Tests"
Cohesion: 0.08
Nodes (21): NewLibraryColumn(), NewListColumn(), TestGetSeasonPlaylist(), TestListColumnHideWatched(), TestListColumnSetItemsPreservesSelection(), TestListColumnSetSelectedByIDHonorsSort(), TestListColumnSortSelectionAndFilter(), Model (+13 more)

### Community 5 - "Player and Scrobbler"
Cohesion: 0.08
Nodes (13): Launcher, NewLauncher(), NewService(), subFileArgs(), TestLookupSeekFlag(), mockPlaybackClient, PlaybackHandle, PlayerDef (+5 more)

### Community 6 - "Jellyfin Client API"
Cohesion: 0.10
Nodes (6): buildAuthHeader(), AuthFlow, AuthResult, Client, appendAPIKey(), firstNonEmpty()

### Community 7 - "Jellyfin Metadata Mapping"
Cohesion: 0.13
Nodes (28): extractAudioInfo(), extractBitrate(), extractResolution(), extractVideoCodec(), formatJellyfinDate(), mapEpisode(), MapEpisodes(), MapLibraries() (+20 more)

### Community 8 - "TUI Message Types"
Cohesion: 0.06
Nodes (29): ClearLibraryStatusMsg, ClearStatusMsg, ContinueWatchingLoadedMsg, EpisodesLoadedMsg, ErrMsg, LibrariesLoadedMsg, LibrarySyncProgressMsg, LogoutCompleteMsg (+21 more)

### Community 11 - "Plex Metadata Mapping"
Cohesion: 0.15
Nodes (19): findBestMedia(), mapEpisode(), MapEpisodes(), MapLibraryContent(), MapMediaItem(), mapMovie(), MapMovies(), MapOnDeck() (+11 more)

### Community 12 - "Fuzzy Search Engine"
Cohesion: 0.13
Nodes (18): FilterItem, FilterResult, allowedTypos(), compareFuzzyMatches(), dedupeAndSort(), findBestTokenMatch(), FuzzySearch(), levenshteinWithPositions() (+10 more)

### Community 15 - "Integration Tests"
Cohesion: 0.15
Nodes (12): TestLocalViews(), TestRecentlyAddedMixed(), TestSmartFilteredShows(), TestFetchContinueWatching(), TestFetchLibrariesSavesToStore(), TestFetchMoviesPaginatesAndReportsProgress(), TestSyncLibraryUsesFreshCache(), TestQueueOperations() (+4 more)

### Community 18 - "Sort Modal Component"
Cohesion: 0.16
Nodes (9): DefaultDirection(), EpisodeSortOptions(), MixedSortOptions(), MovieSortOptions(), ShowSortOptions(), SortDirection, SortField, SortModal (+1 more)

### Community 19 - "Project Metadata"
Cohesion: 0.24
Nodes (16): Cue Project README, SuperCoolPencil (GitHub Owner), CI GitHub Workflow, Cue Example Configuration, Cue Media Player, GoReleaser, GoReleaser Build Configuration, Jellyfin Media Server (+8 more)

### Community 21 - "Plex API DTOs"
Cohesion: 0.13
Nodes (13): APIResponse, Directory, flexBool, Guid, Media, MediaContainer, Metadata, Part (+5 more)

### Community 25 - "Jellyfin API DTOs"
Cohesion: 0.15
Nodes (12): AuthResponse, ImageTags, Item, ItemsResponse, MediaSource, MediaStream, PlaybackInfoResponse, SearchHint (+4 more)

### Community 28 - "InputModal Component"
Cohesion: 0.20
Nodes (3): NewInputModal(), TestInputModalSubmitAndEscape(), InputModal

### Community 29 - "Mock Library Client"
Cohesion: 0.18
Nodes (4): fakeLibraryClient, flattenMixed(), flattenMovies(), flattenShows()

### Community 30 - "View Rendering"
Cohesion: 0.31
Nodes (4): Model, helpEntry, renderConfirmDialog(), RenderSpinner()

### Community 31 - "Plex Auth Flow"
Cohesion: 0.29
Nodes (5): NewAuthClient(), NewAuthFlow(), AuthClient, AuthFlow, AuthResult

### Community 32 - "Keymap Bindings"
Cohesion: 0.22
Nodes (4): GlobalSearchKeyMap, ListColumnKeyMap, PlaylistModalKeyMap, SortModalKeyMap

### Community 33 - "Store Cache Wrapper"
Cohesion: 0.25
Nodes (5): listItemWrapper, cleanupLegacyJSONCache(), hashServerURL(), unwrapListItems(), wrapListItems()

### Community 34 - "Auth Adapter Layer"
Cohesion: 0.29
Nodes (4): AuthFlow, AuthResult, jellyfinAuthAdapter, plexAuthAdapter

### Community 40 - "Playback Interfaces"
Cohesion: 0.50
Nodes (3): PlayableMedia, PlaybackClient, Subtitle

## Knowledge Gaps
- **111 isolated node(s):** `columnLayout`, `LibrariesLoadedMsg`, `MoviesLoadedMsg`, `ShowsLoadedMsg`, `MixedLibraryLoadedMsg` (+106 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewLibraryStore()` connect `Integration Tests` to `Store Cache Wrapper`, `Community 42`, `System Bootstrap`, `Column Stack Tests`?**
  _High betweenness centrality (0.075) - this node is a cross-community bridge._
- **Why does `run()` connect `System Bootstrap` to `UI Components Suite`, `Player and Scrobbler`, `Integration Tests`?**
  _High betweenness centrality (0.075) - this node is a cross-community bridge._
- **What connects `columnLayout`, `LibrariesLoadedMsg`, `MoviesLoadedMsg` to the rest of the system?**
  _111 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `ListColumn Component Core` be split into smaller, more focused modules?**
  _Cohesion score 0.05088919288645691 - nodes in this community are weakly interconnected._
- **Should `TUI App Model` be split into smaller, more focused modules?**
  _Cohesion score 0.05495151337055539 - nodes in this community are weakly interconnected._
- **Should `UI Components Suite` be split into smaller, more focused modules?**
  _Cohesion score 0.05442176870748299 - nodes in this community are weakly interconnected._
- **Should `System Bootstrap` be split into smaller, more focused modules?**
  _Cohesion score 0.06205673758865248 - nodes in this community are weakly interconnected._