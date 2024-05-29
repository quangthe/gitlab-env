package main

import (
	"github.com/urfave/cli/v2"
	"testing"
)

func Test_buildEnvVarsExportUrl(t *testing.T) {
	tests := []struct {
		name string
		args envVarsExportArgs
		want string
	}{
		{
			name: "build url",
			args: envVarsExportArgs{
				gitlabToken: "test-token",
				gitlabHost:  "gitlab.example.com",
				projectId:   "1",
				scopes:      *cli.NewStringSlice("dev"),
				pageSize:    150,
				pageNumber:  1,
			},
			want: "https://gitlab.example.com/api/v4/projects/1/variables?page=1&per_page=150",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildEnvVarsExportUrl(tt.args); got != tt.want {
				t.Errorf("buildEnvVarsExportUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
