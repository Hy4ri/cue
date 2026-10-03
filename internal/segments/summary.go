package segments

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
)

// IntroSummary reports saved analysis, not an assertion about unscanned episodes.
type IntroSummary struct{ Detected, Analyzed, Pending int }

func (c Cache) ShowIntroSummary(server, user, show string) (IntroSummary, error) {
	return c.introSummary(server, user, show, "", nil)
}

func (c Cache) EpisodeIntroSummary(server, user, show string, episodes map[string]bool) (IntroSummary, error) {
	return c.introSummary(server, user, show, "", episodes)
}

func (c Cache) introSummary(server, user, show, season string, episodes map[string]bool) (IntroSummary, error) {
	var summary IntroSummary
	if c.Root == "" {
		return summary, nil
	}
	directory := filepath.Join(c.ShowDir(server, user, show), "analysis")
	files, err := os.ReadDir(directory)
	if os.IsNotExist(err) {
		return summary, nil
	}
	if err != nil {
		return summary, err
	}
	type record struct {
		entry    Cached
		modified time.Time
	}
	latest := map[string]record{}
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(directory, file.Name()))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return summary, err
		}
		var entry Cached
		if json.Unmarshal(b, &entry) != nil || entry.Version != Version || entry.Identity.Server != server || entry.Identity.User != user || entry.Identity.Show != show || entry.Identity.Item == "" {
			continue
		}
		if season != "" && entry.SeasonID != season {
			continue
		}
		info, err := file.Info()
		if episodes != nil && !episodes[entry.Identity.Item] {
			continue
		}
		if err != nil {
			continue
		}
		previous, ok := latest[entry.Identity.Item]
		if !ok || info.ModTime().After(previous.modified) {
			latest[entry.Identity.Item] = record{entry, info.ModTime()}
		}
	}
	for _, record := range latest {
		entry := record.entry
		if entry.Membership == "" {
			summary.Pending++
			continue
		}
		summary.Analyzed++
		for _, segment := range domain.ValidSkipSegments(entry.Segments, entry.Identity.DurationMs) {
			if segment.Kind == "intro" {
				summary.Detected++
				break
			}
		}
	}
	return summary, nil
}

func (c Cache) SeasonIntroSummary(server, user, show, season string) (IntroSummary, error) {
	return c.introSummary(server, user, show, season, nil)
}
