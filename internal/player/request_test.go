package player

import (
	"context"
	"testing"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
)

func TestPlaybackRequestStartsAtSelectedItemAndSnapshotsQueue(t *testing.T) {
	item := domain.MediaItem{ID: "selected", ViewOffset: 12 * time.Minute}
	queue := []domain.MediaItem{{ID: "previous"}, item, {ID: "next"}}
	request, err := NewRequest(item, true, true, queue...)
	if err != nil {
		t.Fatal(err)
	}
	queue[1].ID = "changed"
	if request.StartIndex != 1 || request.Queue[1].ID != item.ID || request.Offset != item.ViewOffset {
		t.Fatalf("invalid playback snapshot: %#v", request)
	}
}

func TestPlaybackRequestWithAutoplayOffContainsOnlySelection(t *testing.T) {
	item := domain.MediaItem{ID: "movie", ViewOffset: time.Minute}
	request, err := NewRequest(item, false, false, domain.MediaItem{ID: "unrelated"})
	if err != nil {
		t.Fatal(err)
	}
	if len(request.Queue) != 1 || request.Queue[0].ID != item.ID || request.Offset != 0 || request.StartIndex != 0 {
		t.Fatalf("unexpected request: %#v", request)
	}
}

func TestPlaybackRequestRejectsMissingSelectionBeforeNetworkOrLaunch(t *testing.T) {
	item := domain.MediaItem{ID: "movie"}
	if _, err := NewRequest(item, true, true, domain.MediaItem{ID: "unrelated"}); err == nil {
		t.Fatal("accepted unrelated queue")
	}
	service := &Service{} // No client or launcher: validation must happen first.
	if _, err := service.Start(context.Background(), Request{Item: item, Queue: []domain.MediaItem{{ID: "unrelated"}}, StartIndex: 0}); err == nil {
		t.Fatal("accepted invalid start index")
	}
}
