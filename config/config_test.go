package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectedErr bool
	}{
		{
			name: "Valid config",
			config: Config{
				Server: struct {
					Address    string `default:"127.0.0.1"`
					Port       string `default:"8080"`
					MaxClients int    `default:"10"`
				}{
					Address:    "127.0.0.1",
					Port:       "8080",
					MaxClients: 10,
				},
			},
			expectedErr: false,
		},
		{
			name: "Invalid port number",
			config: Config{
				Server: struct {
					Address    string `default:"127.0.0.1"`
					Port       string `default:"8080"`
					MaxClients int    `default:"10"`
				}{
					Address:    "127.0.0.1",
					Port:       "invalid",
					MaxClients: 10,
				},
			},
			expectedErr: true,
		},
		{
			name: "Port out of range",
			config: Config{
				Server: struct {
					Address    string `default:"127.0.0.1"`
					Port       string `default:"8080"`
					MaxClients int    `default:"10"`
				}{
					Address:    "127.0.0.1",
					Port:       "70000",
					MaxClients: 10,
				},
			},
			expectedErr: true,
		},
		{
			name: "Max clients less than 1",
			config: Config{
				Server: struct {
					Address    string `default:"127.0.0.1"`
					Port       string `default:"8080"`
					MaxClients int    `default:"10"`
				}{
					Address:    "127.0.0.1",
					Port:       "8080",
					MaxClients: 0,
				},
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
