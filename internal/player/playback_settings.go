package player

import (
	"crypto/sha256"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

//go:embed playback_settings.lua
var playbackSettingsLua string

// preparePlaybackSettings gives mpv its own event-driven preference store.
func preparePlaybackSettings(showID string) (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	root = filepath.Join(root, "cue", "playback-settings")
	if err := os.MkdirAll(root, 0700); err != nil {
		return "", err
	}
	settings := filepath.Join(root, fmt.Sprintf("%x.json", sha256.Sum256([]byte(showID))))
	script, err := os.CreateTemp("", "cue-settings-*.lua")
	if err != nil {
		return "", err
	}
	source := "local settings_path = " + strconv.Quote(settings) + "\n" + playbackSettingsLua
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
