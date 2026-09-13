package model

import (
	"strings"
	"testing"

	"github.com/bjarneo/cliamp/playlist"
	"github.com/bjarneo/cliamp/provider"
	"github.com/bjarneo/cliamp/ui"
)

func dated(title, date string) playlist.Track {
	return playlist.Track{Path: "/" + title, Title: title, DurationSecs: 60,
		ProviderMeta: map[string]string{provider.MetaPodcastPublished: date}}
}

func TestEpisodeDateColumn(t *testing.T) {
	tests := []struct {
		name   string
		on     bool
		width  int
		tracks []playlist.Track
		want   bool
	}{
		{"on, wide, dated", true, 80, []playlist.Track{dated("a", "2026-09-12")}, true},
		{"off", false, 80, []playlist.Track{dated("a", "2026-09-12")}, false},
		{"too narrow", true, 60, []playlist.Track{dated("a", "2026-09-12")}, false},
		{"nothing dated", true, 80, []playlist.Track{{Path: "/radio", Title: "Radio"}}, false},
		{"one dated among many", true, 80, []playlist.Track{{Path: "/radio"}, dated("a", "2026-09-12")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := ui.PanelWidth
			ui.PanelWidth = tt.width
			t.Cleanup(func() { ui.PanelWidth = old })
			m := Model{showEpisodeDates: tt.on}
			if got := m.episodeDateColumn(tt.tracks); got != tt.want {
				t.Errorf("episodeDateColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Undated rows get a blank cell of the same width, so titles line up.
func TestEpisodeDateCellKeepsAlignment(t *testing.T) {
	with := stripAnsi(episodeDateCell(dated("a", "2026-09-12")))
	without := episodeDateCell(playlist.Track{})
	if len(with) != len(without) {
		t.Errorf("cell widths differ: %q (%d) vs %q (%d)", with, len(with), without, len(without))
	}
	if !strings.HasPrefix(with, "2026-09-12") {
		t.Errorf("dated cell = %q, want the date first", with)
	}
}

func TestRenderQueueBodyShowsEpisodeDateBeforeTitle(t *testing.T) {
	old := ui.PanelWidth
	ui.PanelWidth = 80
	t.Cleanup(func() { ui.PanelWidth = old })
	m := &Model{playlist: playlist.New(), plVisible: 5, showEpisodeDates: true}
	m.playlist.Replace([]playlist.Track{dated("Dated", "2026-09-10")})
	m.playlist.Queue(0)

	body := stripAnsi(m.renderQueueBody())

	if !strings.Contains(body, "1. 2026-09-10  Dated") {
		t.Errorf("queue body = %q, want the date between the number and the title", body)
	}
}

func TestToggleEpisodeDatesPersists(t *testing.T) {
	saver := &recordingSaver{}
	m := &Model{showEpisodeDates: true, configSaver: saver}

	m.toggleEpisodeDates()

	if m.showEpisodeDates {
		t.Error("toggle did not turn the column off")
	}
	if got := saver.saved["show_episode_dates"]; got != "false" {
		t.Errorf("saved = %q, want \"false\"", got)
	}
	m.toggleEpisodeDates()
	if got := saver.saved["show_episode_dates"]; got != "true" {
		t.Errorf("saved = %q, want \"true\"", got)
	}
}
