package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func TestResumeDialogShowsEpisodeAndPosition(t *testing.T) {
	m := Model{Width: 100, Height: 30, pendingPlayback: &domain.MediaItem{Type: domain.MediaTypeEpisode, ShowTitle: "Breaking Bad", Title: "Pilot", SeasonNum: 1, EpisodeNum: 1, ViewOffset: 12*time.Minute + 34*time.Second, Duration: 58*time.Minute + 9*time.Second}}
	rendered := ansi.Strip(m.renderResumeConfirmation())
	for _, want := range []string{"Breaking Bad", "S01E01 · Pilot", "Resume from 12:34 / 58:09", "Y Resume", "N Start Over", "Esc Cancel"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("missing %q in %s", want, rendered)
		}
	}
}
func TestResumeDialogMovieAndLongTimestamp(t *testing.T) {
	m := Model{Width: 100, Height: 30, pendingPlayback: &domain.MediaItem{Type: domain.MediaTypeMovie, Title: "Movie", ViewOffset: time.Hour + 2*time.Minute + 3*time.Second}}
	rendered := ansi.Strip(m.renderResumeConfirmation())
	if !strings.Contains(rendered, "Resume from 1:02:03") || strings.Contains(rendered, "S00E00") {
		t.Fatal(rendered)
	}
}
func TestConfirmationFitsNarrowTerminal(t *testing.T) {
	m := Model{Width: 40, Height: 24, pendingPlayback: &domain.MediaItem{Title: "Movie", ViewOffset: time.Minute}}
	rendered := m.renderResumeConfirmation()
	if lipgloss.Width(rendered) > m.Width {
		t.Fatalf("dialog width %d exceeds %d", lipgloss.Width(rendered), m.Width)
	}
	for _, want := range []string{"Y Resume", "N Start Over", "Esc Cancel"} {
		if !strings.Contains(ansi.Strip(rendered), want) {
			t.Errorf("missing %q", want)
		}
	}
}
