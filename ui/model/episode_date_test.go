package model

import (
	"strings"
	"testing"

	"github.com/bjarneo/cliamp/playlist"
	"github.com/bjarneo/cliamp/provider"
	"github.com/bjarneo/cliamp/ui"
)

func TestTrackTrailer(t *testing.T) {
	dated := playlist.Track{DurationSecs: 3768, ProviderMeta: map[string]string{provider.MetaPodcastPublished: "2026-09-10"}}
	tests := []struct {
		name  string
		width int
		track playlist.Track
		want  string
	}{
		{"date and duration", 80, dated, "2026-09-10  1:02:48"},
		{"no date", 80, playlist.Track{DurationSecs: 3768}, "1:02:48"},
		{"date without duration", 80, playlist.Track{ProviderMeta: map[string]string{provider.MetaPodcastPublished: "2026-09-10"}}, "2026-09-10"},
		{"too narrow for the date", 60, dated, "1:02:48"},
		{"nothing", 80, playlist.Track{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := ui.PanelWidth
			ui.PanelWidth = tt.width
			t.Cleanup(func() { ui.PanelWidth = old })
			if got := trackTrailer(tt.track); got != tt.want {
				t.Errorf("trackTrailer() = %q, want %q", got, tt.want)
			}
		})
	}
}

// The queue view shares the trailer, so a dated episode shows its date there too.
func TestRenderQueueBodyShowsEpisodeDate(t *testing.T) {
	old := ui.PanelWidth
	ui.PanelWidth = 80
	t.Cleanup(func() { ui.PanelWidth = old })
	m := &Model{playlist: playlist.New(), plVisible: 5}
	m.playlist.Replace([]playlist.Track{{
		Path: "/a.mp3", Title: "Dated", DurationSecs: 60,
		ProviderMeta: map[string]string{provider.MetaPodcastPublished: "2026-09-10"},
	}})
	m.playlist.Queue(0)

	if body := stripAnsi(m.renderQueueBody()); !strings.Contains(body, "2026-09-10  1:00") {
		t.Errorf("queue body = %q, want the date before the duration", body)
	}
}
