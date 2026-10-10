package player

import (
	"fmt"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
)

// Request is a snapshot of one playback intent. Queue order and the starting
// entry are explicit; the launcher never substitutes the first unrelated item.
type Request struct {
	Item       domain.MediaItem
	Queue      []domain.MediaItem
	StartIndex int
	Offset     time.Duration
}

func NewRequest(item domain.MediaItem, resume, autoplay bool, queue ...domain.MediaItem) (Request, error) {
	request := Request{Item: item}
	if resume {
		request.Offset = item.ViewOffset
	}
	if !autoplay || len(queue) == 0 {
		request.Queue = []domain.MediaItem{item}
	} else {
		request.Queue = append([]domain.MediaItem(nil), queue...)
		request.StartIndex = -1
		for i, entry := range request.Queue {
			if entry.ID == item.ID {
				request.StartIndex = i
				request.Queue[i] = item
				break
			}
		}
	}
	return request, request.Validate()
}

func (r Request) Validate() error {
	if r.Item.ID == "" {
		return fmt.Errorf("playback item has no ID")
	}
	if r.StartIndex < 0 || r.StartIndex >= len(r.Queue) || r.Queue[r.StartIndex].ID != r.Item.ID {
		return fmt.Errorf("selected playback item %s is missing from the queue at its start index", r.Item.ID)
	}
	if r.Offset < 0 {
		return fmt.Errorf("playback offset must not be negative")
	}
	for _, entry := range r.Queue {
		if entry.ID == "" {
			return fmt.Errorf("playback queue contains an item without an ID")
		}
	}
	return nil
}
