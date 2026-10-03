package segments

import (
	"context"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractWithPrivateProxy(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("FFmpeg unavailable")
	}
	capabilities, err := exec.Command("ffmpeg", "-hide_banner", "-muxers").CombinedOutput()
	if err != nil || !strings.Contains(string(capabilities), "chromaprint") {
		t.Skip("Chromaprint unavailable")
	}
	data := make([]byte, 44+40*11025*2)
	copy(data, "RIFF")
	binary.LittleEndian.PutUint32(data[4:], uint32(len(data)-8))
	copy(data[8:], "WAVEfmt ")
	binary.LittleEndian.PutUint32(data[16:], 16)
	binary.LittleEndian.PutUint16(data[20:], 1)
	binary.LittleEndian.PutUint16(data[22:], 1)
	binary.LittleEndian.PutUint32(data[24:], 11025)
	binary.LittleEndian.PutUint32(data[28:], 22050)
	binary.LittleEndian.PutUint16(data[32:], 2)
	binary.LittleEndian.PutUint16(data[34:], 16)
	copy(data[36:], "data")
	binary.LittleEndian.PutUint32(data[40:], uint32(len(data)-44))
	path := filepath.Join(t.TempDir(), "test.wav")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("api_key") != "secret" {
			t.Error("source lost authentication")
		}
		http.ServeFile(w, r, path)
	}))
	defer server.Close()
	fp, err := Extract(context.Background(), server.URL+"?api_key=secret", 0, 30, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(fp) < 180 || len(fp) > 260 {
		t.Fatalf("unexpected fingerprint length %d", len(fp))
	}
}
