package main

import (
	"testing"
)

func Test_buildCreateEnvVarUrl(t *testing.T) {
	tests := []struct {
		name string
		args envVarsImportArgs
		want string
	}{
		{
			name: "build url",
			args: envVarsImportArgs{
				gitlabToken: "test-token",
				gitlabHost:  "gitlab.example.com",
				projectId:   "2",
			},
			want: "https://gitlab.example.com/api/v4/projects/2/variables",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildCreateEnvVarUrl(tt.args); got != tt.want {
				t.Errorf("buildCreateEnvVarUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildUpdateEnvVarUrl(t *testing.T) {
	tests := []struct {
		name     string
		args     envVarsImportArgs
		envKey   string
		envScope string
		want     string
	}{
		{
			name: "build url",
			args: envVarsImportArgs{
				gitlabToken: "test-token",
				gitlabHost:  "gitlab.example.com",
				projectId:   "2",
			},
			envKey:   "SERVICE_TOKEN",
			envScope: "dev",
			want:     "https://gitlab.example.com/api/v4/projects/2/variables/SERVICE_TOKEN?filter[environment_scope]=dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildUpdateEnvVarUrl(tt.args, tt.envKey, tt.envScope); got != tt.want {
				t.Errorf("buildUpdateEnvVarUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
