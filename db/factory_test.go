package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		conf    Config
		wantErr error
	}{
		{name: "sqlite", conf: Config{Engine: "sqlite", File: "zbot.db"}},
		{name: "no engine", conf: Config{}, wantErr: ErrNoEngine},
		{name: "unsupported engine", conf: Config{Engine: "mysql"}, wantErr: ErrUnsupportedEngine},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.conf)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.conf.Engine, got.Engine())
			assert.IsType(t, &ZbotDatabaseSqlite{}, got)
			assert.Equal(t, tt.conf.File, got.(*ZbotDatabaseSqlite).File)
		})
	}
}
