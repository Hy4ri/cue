package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/SuperCoolPencil/cue/internal/segments"
	"github.com/SuperCoolPencil/cue/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

type introStatusLoadedMsg struct{ Key, Label string }

func introStatusLabel(summary segments.IntroSummary, progress segments.AnalysisProgress) string {
	if progress.Active {
		if progress.Matching {
			return fmt.Sprintf("Matching… (%d episodes; %d intros detected)", progress.Total, summary.Detected)
		}
		return fmt.Sprintf("Analyzing… (%d/%d fingerprinted; %d intros detected)", progress.Completed, progress.Total, summary.Detected)
	}
	if summary.Pending > 0 {
		return fmt.Sprintf("%d detected · %d analyzed · %d pending", summary.Detected, summary.Analyzed, summary.Pending)
	}
	if summary.Analyzed > 0 {
		return fmt.Sprintf("Analysis complete · intros found in %d/%d episodes", summary.Detected, summary.Analyzed)
	}
	return "Not analyzed yet"
}

// updateIntroStatus performs disk reads asynchronously and refreshes while the
// show remains selected so background analysis appears without navigation.
func (m *Model) updateIntroStatus(item interface{}) tea.Cmd {
	season := introSeason(item)
	if season == nil && m.ColumnStack != nil {
		for n := m.ColumnStack.Len() - 1; n >= 0; n-- {
			if parent := introSeason(m.ColumnStack.Get(n).SelectedItem()); parent != nil {
				season = parent
				break
			}
		}
	}
	if season == nil || m.AppConfig == nil {
		m.introStatusKey = ""
		m.introStatusPending = false
		m.introStatusLabel = ""
		m.Inspector.SetIntroStatus("")
		return nil
	}
	server, user := m.AppConfig.Server.URL, m.AppConfig.Server.UserID
	key := strings.Join([]string{server, user, season.ShowID, season.ID}, "\x00")
	if key != m.introStatusKey {
		m.introStatusKey = key
		m.introStatusPending = false
		m.introStatusChecked = time.Time{}
		m.Inspector.SetIntroStatus("Checking…")
		m.introStatusLabel = "Checking…"
	}
	if m.introStatusPending || time.Since(m.introStatusChecked) < 2*time.Second {
		return nil
	}
	m.introStatusPending = true
	m.introStatusChecked = time.Now()
	showID, seasonID := season.ShowID, season.ID
	return func() tea.Msg {
		summary, err := segments.DefaultCache().SeasonIntroSummary(server, user, showID, seasonID)
		progress := segments.SeasonAnalysisProgress(server, user, showID, seasonID)
		label := introStatusLabel(summary, progress)
		if err != nil {
			label = "Status unavailable"
		}
		return introStatusLoadedMsg{Key: key, Label: label}
	}
}

func introSeason(item interface{}) *domain.Season {
	switch v := item.(type) {
	case *domain.Season:
		return v
	case *components.SeasonHeader:
		return v.Season
	}
	return nil
}
