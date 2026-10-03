package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/SuperCoolPencil/cue/internal/config"
	"github.com/SuperCoolPencil/cue/internal/mediaserver"
	"github.com/SuperCoolPencil/cue/internal/segments"
)

// runAnalyze explicitly analyzes one server season, never during playback.
func runAnalyze(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var season string
	var window, audio int
	fs.StringVar(&season, "season", "", "server season ID to analyze")
	fs.IntVar(&window, "window", 300, "seconds from the beginning and end (30-900)")
	fs.IntVar(&audio, "audio-track", 0, "zero-based audio track to fingerprint")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if season == "" || fs.NArg() != 0 || window < 30 || window > 900 || audio < 0 {
		fmt.Fprintln(stderr, "Usage: cue analyze --season <server-season-id> [--window 300] [--audio-track 0]")
		return 2
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	client, err := mediaserver.NewClient(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	analyzer := segments.Analyzer{Client: client, Cache: segments.DefaultCache(), Server: cfg.Server.URL, User: cfg.Server.UserID, WindowSeconds: window, AudioTrack: audio}
	count, err := analyzer.AnalyzeSeason(ctx, season)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Saved %d manual skip suggestions. Playback will use them for the same media revision.\n", count)
	return 0
}
