package segments

import (
	"encoding/xml"
	"github.com/SuperCoolPencil/cue/internal/domain"
	"testing"
)

func TestChapterXMLBoundariesAndRemainder(t *testing.T) {
	b, err := ChapterXML([]domain.SkipSegment{{Kind: "intro", StartMs: 0, EndMs: 90000}, {Kind: "outro", StartMs: 1300000, EndMs: 1390000}}, 1440000)
	if err != nil {
		t.Fatal(err)
	}
	var value chaptersXML
	if err := xml.Unmarshal(b, &value); err != nil {
		t.Fatal(err)
	}
	expected := []string{"00:00:00.000000000", "00:01:30.000000000", "00:21:40.000000000", "00:23:10.000000000"}
	if len(value.Atoms) != len(expected) {
		t.Fatalf("chapters=%+v", value.Atoms)
	}
	for i, want := range expected {
		if value.Atoms[i].Start != want {
			t.Fatalf("boundary=%s want=%s", value.Atoms[i].Start, want)
		}
	}
	if value.Atoms[3].Display.Name != "Episode" {
		t.Fatal("remainder discarded")
	}
}
