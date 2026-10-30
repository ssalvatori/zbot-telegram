package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogsCommandShow(t *testing.T) {
	c := &LogsCommand{}
	if got, want := c.ShowHelp(), LogsHelpMessage; got != want {
		t.Fatalf("ShowHelp = %q, want %q", got, want)
	}
	if got, want := c.ShowDescription(), LogsDescription; got != want {
		t.Fatalf("ShowDescription = %q, want %q", got, want)
	}
}

func TestLogsCommandReadsLastLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "zbot.log")
	entries := make([]string, 0, 25)
	for i := 1; i <= 25; i++ {
		entries = append(entries, "line "+string(rune('A'+i-1)))
	}
	if err := os.WriteFile(path, []byte(strings.Join(entries, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := &LogsCommand{Path: path}
	out, err := c.ProcessText("", userTest, 1, "chat", true, true)
	if err != nil {
		t.Fatalf("ProcessText returned error: %v", err)
	}

	gotLines := strings.Split(out, "\n")
	if len(gotLines) != logsDefaultLimit {
		t.Fatalf("expected %d lines, got %d", logsDefaultLimit, len(gotLines))
	}
	if want := "line " + string(rune('A'+25-logsDefaultLimit)); gotLines[0] != want {
		t.Fatalf("first kept line should be %q, got %q", want, gotLines[0])
	}

	// Default limit is logsDefaultLimit regardless of the file size.
	again, _ := (&LogsCommand{Path: path}).ProcessText("", userTest, 1, "chat", true, true)
	if again != out {
		t.Errorf("expected same output on a second read; got %d lines", strings.Count(again, "\n")+1)
	}
}

func TestLogsCommandEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "zbot.log")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	c := &LogsCommand{Path: path}
	out, err := c.ProcessText("", userTest, 1, "chat", true, true)
	if err != nil {
		t.Fatalf("ProcessText returned error: %v", err)
	}
	if out != "No log entries yet." {
		t.Errorf("got %q", out)
	}
}

func TestLogsCommandMissingFile(t *testing.T) {
	c := &LogsCommand{Path: filepath.Join(t.TempDir(), "does-not-exist.log")}
	out, err := c.ProcessText("", userTest, 1, "chat", true, true)
	if err != nil {
		t.Fatalf("ProcessText returned error on missing file: %v", err)
	}
	if out != "No log entries yet." {
		t.Errorf("got %q", out)
	}
}

func TestLogsCommandNotInvokedAsCommand(t *testing.T) {
	c := &LogsCommand{Path: filepath.Join(t.TempDir(), "zbot.log")}
	_, err := c.ProcessText("hello", userTest, 1, "chat", true, false)
	if wantErr := ErrNextCommand; !strings.Contains(err.Error(), wantErr.Error()) {
		t.Errorf("got error %v, expected %v", err, wantErr)
	}
}

func TestTruncateMessage(t *testing.T) {
	tests := []struct {
		name      string
		msg       string
		maxLength int
		want      string
	}{
		{
			name:      "shorter than max returns unchanged",
			msg:       "hello",
			maxLength: 10,
			want:      "hello",
		},
		{
			name:      "exactly max length returns unchanged",
			msg:       "12345",
			maxLength: 5,
			want:      "12345",
		},
		{
			name:      "longer than max is truncated from the end",
			msg:       "abcdefghij",
			maxLength: 5,
			want:      "fghij... (truncated)",
		},
		{
			name:      "empty string",
			msg:       "",
			maxLength: 10,
			want:      "",
		},
		{
			name:      "zero max length truncates everything",
			msg:       "hello",
			maxLength: 0,
			want:      "... (truncated)",
		},
		{
			name:      "multibyte characters counted by bytes",
			msg:       "héllo wörld",
			maxLength: 5,
			want:      "örld... (truncated)",
		},
		{
			name:      "single character longer than max",
			msg:       "x",
			maxLength: 0,
			want:      "... (truncated)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateMessage(tt.msg, tt.maxLength)
			assert.Equal(t, tt.want, got)
		})
	}
}
