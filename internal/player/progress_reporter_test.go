package player

import (
	"context"
	"log/slog"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/SuperCoolPencil/cue/internal/domain"
)

type delayedProgressClient struct {
	domain.PlaybackClient
	started, release chan struct{}
	mu               sync.Mutex
	positions        []int64
}

func (c *delayedProgressClient) UpdateProgress(ctx context.Context, _ string, position int64) error {
	if position == 30000 {
		close(c.started)
		select {
		case <-c.release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.positions = append(c.positions, position)
	return nil
}

func TestProgressReporterDrainsOldWritesBeforeFinalBackwardSeek(t *testing.T) {
	client := &delayedProgressClient{started: make(chan struct{}), release: make(chan struct{})}
	reporter := newProgressReporter(client, slog.Default())
	item := domain.MediaItem{ID: "movie"}
	reporter.queue(item, 30000)
	select {
	case <-client.started:
	case <-time.After(time.Second):
		t.Fatal("first write did not start")
	}
	reporter.queue(item, 20000)
	finished := make(chan struct{})
	go func() { reporter.finish(&progressUpdate{item, 10000}); close(finished) }()
	select {
	case <-finished:
		t.Fatal("finish returned before the old write completed")
	default:
	}
	close(client.release)
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("progress writer did not drain")
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	if len(client.positions) != 3 || client.positions[0] != 30000 || client.positions[1] != 20000 || client.positions[2] != 10000 {
		t.Fatalf("writes out of order: %v", client.positions)
	}
}

type shutdownTimelineClient struct {
	lifecyclePlaybackClient
	stopping, release chan struct{}
}

func (c *shutdownTimelineClient) ReportTimeline(ctx context.Context, state, _ string, _, _ int64) error {
	if state == "stopped" {
		close(c.stopping)
		select {
		case <-c.release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func TestPlaybackCloseWaitsForFinalReport(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test uses a POSIX shell")
	}
	client := &shutdownTimelineClient{stopping: make(chan struct{}), release: make(chan struct{})}
	launcher := NewLauncher("/bin/sh", []string{"-c", "exec sleep 10"}, "", nil)
	service := NewService(launcher, client, nil)
	if _, err := service.Play(context.Background(), domain.MediaItem{ID: "movie"}); err != nil {
		t.Fatal(err)
	}
	finished := make(chan struct{})
	go func() { service.Close(); close(finished) }()
	select {
	case <-client.stopping:
	case <-time.After(time.Second):
		close(client.release)
		t.Fatal("no stopped report")
	}
	select {
	case <-finished:
		close(client.release)
		t.Fatal("shutdown returned before final report")
	default:
	}
	close(client.release)
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("shutdown did not finish")
	}
	if _, err := service.Play(context.Background(), domain.MediaItem{ID: "new"}); err == nil {
		t.Fatal("closed service accepted new playback")
	}
}
