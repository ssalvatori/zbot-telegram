package command

import (
	"container/list"
	"errors"
	"testing"
)

func TestHelpCommand(t *testing.T) {
	commands := &CommandsList{List: list.New()}
	commands.Chain("learn", &LearnCommand{}, 0, LearnHelpMessage)
	commands.Chain("who", &WhoCommand{}, 0, WhoHelpMessage)

	tests := []struct {
		name    string
		helper  *HelpCommand
		text    string
		private bool
		want    string
		wantErr error
	}{
		{
			name:   "registered command",
			helper: &HelpCommand{Commands: commands},
			text:   "!help learn",
			want:   LearnHelpMessage,
		},
		{
			name:   "case-insensitive command name",
			helper: &HelpCommand{Commands: commands},
			text:   "!help LEARN ",
			want:   LearnHelpMessage,
		},
		{
			name:   "who command",
			helper: &HelpCommand{Commands: commands},
			text:   "!help who",
			want:   WhoHelpMessage,
		},
		{
			name:    "unknown command",
			helper:  &HelpCommand{Commands: commands},
			text:    "!help unknown",
			wantErr: ErrNextCommand,
		},
		{
			name:    "missing command name",
			helper:  &HelpCommand{Commands: commands},
			text:    "!help",
			wantErr: ErrNextCommand,
		},
		{
			name:    "private message",
			helper:  &HelpCommand{Commands: commands},
			text:    "!help learn",
			private: true,
			wantErr: ErrNextCommand,
		},
		{
			name:    "unset command list",
			helper:  &HelpCommand{},
			text:    "!help learn",
			wantErr: ErrNextCommand,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			message, err := test.helper.ProcessText(test.text, userTest, 0, "test", test.private)
			if message != test.want {
				t.Errorf("HelpCommand.ProcessText(%q) = %q, want %q", test.text, message, test.want)
			}
			if !errors.Is(err, test.wantErr) {
				t.Errorf("HelpCommand.ProcessText(%q) error = %v, want %v", test.text, err, test.wantErr)
			}
		})
	}
}
