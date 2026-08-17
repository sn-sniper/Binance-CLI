package config

import (
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "empty api key",
			cfg:     &Config{APIKey: "", SecretKey: "secret"},
			wantErr: true,
		},
		{
			name:    "empty secret key",
			cfg:     &Config{APIKey: "key", SecretKey: ""},
			wantErr: true,
		},
		{
			name:    "placeholder api key",
			cfg:     &Config{APIKey: "your_api_key_here", SecretKey: "secret"},
			wantErr: true,
		},
		{
			name:    "placeholder secret key",
			cfg:     &Config{APIKey: "key", SecretKey: "your_secret_key_here"},
			wantErr: true,
		},
		{
			name:    "valid credentials",
			cfg:     &Config{APIKey: "my_valid_key", SecretKey: "my_valid_secret"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.cfg.HasCredentials() == tt.wantErr {
				t.Errorf("Config.HasCredentials() = %v, want %v", tt.cfg.HasCredentials(), !tt.wantErr)
			}
		})
	}
}
