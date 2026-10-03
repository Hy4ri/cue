package tui

import (
	"testing"

	"github.com/SuperCoolPencil/cue/internal/segments"
)

func TestIntroStatusShowsActivityAndPartialResults(t *testing.T) {
	for _, test := range []struct {
		name     string
		summary  segments.IntroSummary
		progress segments.AnalysisProgress
		want     string
	}{
		{"extracting", segments.IntroSummary{Detected: 4, Analyzed: 10}, segments.AnalysisProgress{Active: true, Completed: 2, Total: 10}, "Analyzing… (2/10 fingerprinted; 4 intros detected)"},
		{"matching", segments.IntroSummary{Pending: 10}, segments.AnalysisProgress{Active: true, Matching: true, Completed: 10, Total: 10}, "Matching… (10 episodes; 0 intros detected)"},
		{"interrupted", segments.IntroSummary{Detected: 4, Analyzed: 4, Pending: 6}, segments.AnalysisProgress{}, "4 detected · 4 analyzed · 6 pending"},
		{"finished", segments.IntroSummary{Detected: 4, Analyzed: 10}, segments.AnalysisProgress{}, "Analysis complete · intros found in 4/10 episodes"},
		{"no matches", segments.IntroSummary{Analyzed: 10}, segments.AnalysisProgress{}, "Analysis complete · intros found in 0/10 episodes"},
		{"not started", segments.IntroSummary{}, segments.AnalysisProgress{}, "Not analyzed yet"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := introStatusLabel(test.summary, test.progress); got != test.want {
				t.Fatalf("status = %q, want %q", got, test.want)
			}
		})
	}
}
