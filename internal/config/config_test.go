package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "missing LINEAR_API_KEY",
			env: map[string]string{
				"LINEAR_TEAM_ID": "team-123",
				"DATABASE_URL":   "postgres://localhost/test",
			},
			wantErr: true,
		},
		{
			name: "missing LINEAR_TEAM_ID",
			env: map[string]string{
				"LINEAR_API_KEY": "key-123",
				"DATABASE_URL":   "postgres://localhost/test",
			},
			wantErr: true,
		},
		{
			name: "missing DATABASE_URL",
			env: map[string]string{
				"LINEAR_API_KEY": "key-123",
				"LINEAR_TEAM_ID": "team-123",
			},
			wantErr: true,
		},
		{
			name: "valid config",
			env: map[string]string{
				"LINEAR_API_KEY": "key-123",
				"LINEAR_TEAM_ID": "team-123",
				"DATABASE_URL":   "postgres://localhost/test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("LINEAR_API_KEY")
			os.Unsetenv("LINEAR_TEAM_ID")
			os.Unsetenv("DATABASE_URL")

			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.env {
					os.Unsetenv(k)
				}
			}()

			cfg, err := Load()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, tt.env["LINEAR_API_KEY"], cfg.APIKey)
				assert.Equal(t, tt.env["LINEAR_TEAM_ID"], cfg.TeamID)
				assert.Equal(t, tt.env["DATABASE_URL"], cfg.DatabaseURL)
			}
		})
	}
}
