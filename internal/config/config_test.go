package config

import "testing"

func TestLoad(t *testing.T) {
	allVars := map[string]string{
		"LINEAR_API_KEY": "test-api-key",
		"LINEAR_TEAM_ID": "test-team-id",
		"DATABASE_URL":   "postgres://localhost/test",
	}

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
		wantCfg *Config
	}{
		{
			name:    "missing LINEAR_API_KEY returns error",
			env:     map[string]string{"LINEAR_TEAM_ID": "tid", "DATABASE_URL": "url"},
			wantErr: true,
		},
		{
			name:    "missing LINEAR_TEAM_ID returns error",
			env:     map[string]string{"LINEAR_API_KEY": "key", "DATABASE_URL": "url"},
			wantErr: true,
		},
		{
			name:    "missing DATABASE_URL returns error",
			env:     map[string]string{"LINEAR_API_KEY": "key", "LINEAR_TEAM_ID": "tid"},
			wantErr: true,
		},
		{
			name:    "all variables set returns correct Config",
			env:     allVars,
			wantErr: false,
			wantCfg: &Config{
				APIKey:      "test-api-key",
				TeamID:      "test-team-id",
				DatabaseURL: "postgres://localhost/test",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("LINEAR_API_KEY", "")
			t.Setenv("LINEAR_TEAM_ID", "")
			t.Setenv("DATABASE_URL", "")
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			cfg, err := Load()
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.APIKey != tc.wantCfg.APIKey || cfg.TeamID != tc.wantCfg.TeamID || cfg.DatabaseURL != tc.wantCfg.DatabaseURL {
				t.Errorf("got %+v, want %+v", cfg, tc.wantCfg)
			}
		})
	}
}
