package command

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := tt.helper.ProcessText(tt.text, userTest, 0, "test", tt.private)
			assert.Equal(t, tt.want, message, "HelpCommand.ProcessText(%q)", tt.text)
			if tt.wantErr == nil {
				assert.NoError(t, err, "HelpCommand.ProcessText(%q)", tt.text)
			} else {
				assert.ErrorIs(t, err, tt.wantErr, "HelpCommand.ProcessText(%q)", tt.text)
			}
		})
	}
}
