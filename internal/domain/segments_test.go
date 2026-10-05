package domain

import "testing"

func TestValidSkipSegments(t *testing.T) {
	input := []SkipSegment{
		{Kind: "intro", StartMs: 200, EndMs: 500}, {Kind: "outro", StartMs: 400, EndMs: 600},
		{Kind: "intro", StartMs: -1, EndMs: 10}, {Kind: "recap", StartMs: 700, EndMs: 800},
		{Kind: "outro", StartMs: 800, EndMs: 1100}, {Kind: "intro", StartMs: 650, EndMs: 700},
	}
	got := ValidSkipSegments(input, 1000)
	if len(got) != 1 || got[0].StartMs != 650 {
		t.Fatalf("invalid/conflicting intervals survived: %+v", got)
	}
	if input[0].StartMs != 200 {
		t.Fatal("modified input")
	}
}
