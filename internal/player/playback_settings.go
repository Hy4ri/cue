package player

import (
	"crypto/sha256"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SuperCoolPencil/cue/internal/segments"
)

//go:embed playback_settings.lua
var playbackSettingsLua string

// preparePlaybackSettings gives mpv its own event-driven preference store.
func preparePlaybackSettings(showID string) (string, error) {
	return prepareScopedPlaybackSettings("", "", showID)
}

func prepareScopedPlaybackSettings(server, user, showID string, mapPath ...func(string) string) (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	legacyRoot := filepath.Join(root, "cue", "playback-settings")
	root = segments.DefaultCache().ShowDir(server, user, showID)
	if err := os.MkdirAll(root, 0700); err != nil {
		return "", err
	}
	settings := filepath.Join(root, "settings.json")
	if _, err := os.Stat(settings); os.IsNotExist(err) {
		legacy := filepath.Join(legacyRoot, fmt.Sprintf("%x.json", sha256.Sum256([]byte(showID))))
		if value, err := os.ReadFile(legacy); err == nil {
			if err := os.WriteFile(settings, value, 0600); err != nil {
				return "", err
			}
		}
	}

	script, err := os.CreateTemp(root, "cue-settings-*.lua")
	if err != nil {
		return "", err
	}
	playerSettings := settings
	if len(mapPath) > 0 {
		playerSettings = mapPath[0](settings)
	}
	source := "local settings_path = " + strconv.Quote(playerSettings) + "\n" + playbackSettingsLua
	_, writeErr := script.WriteString(source)
	closeErr := script.Close()
	if writeErr != nil {
		os.Remove(script.Name())
		return "", writeErr
	}
	if closeErr != nil {
		os.Remove(script.Name())
		return "", closeErr
	}
	return script.Name(), nil
}
