package player

import (
	"context"
	"log/slog"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
)

type progressUpdate struct {
	item     domain.MediaItem
	position int64
}

// A session has one progress writer. Periodic updates are bounded, and finish
// drains them before returning so an older write cannot undo the final seek.
type progressReporter struct {
	updates chan progressUpdate
	done    chan struct{}
}

func newProgressReporter(client domain.PlaybackClient, logger *slog.Logger) *progressReporter {
	r := &progressReporter{updates: make(chan progressUpdate, 1), done: make(chan struct{})}
	go func() {
		defer close(r.done)
		for update := range r.updates {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := client.UpdateProgress(ctx, update.item.ID, update.position); err != nil {
				logger.Warn("failed to update progress", "itemID", update.item.ID, "error", err)
			}
			cancel()
		}
	}()
	return r
}

func (r *progressReporter) queue(item domain.MediaItem, position int64) {
	select {
	case r.updates <- progressUpdate{item, position}:
	default:
	}
}

func (r *progressReporter) finish(final *progressUpdate) {
	if final != nil {
		r.updates <- *final
	}
	close(r.updates)
	<-r.done
}
