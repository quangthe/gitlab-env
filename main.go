package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	flagGitlabToken = "gitlab-token"
	flagGitlabHost  = "gitlab-host"
	flagProjectId   = "project-id"
	flagEnvScopes   = "env-scopes"
	flagPageSize    = "page-size"
	flagPageNumber  = "page-number"
	flagOutputFile  = "output-file"
	flagInputFile   = "input-file"
)

const (
	defaultJsonFile = "gitlab-env-vars.json"
)

type GitlabEnvVar struct {
	VariableType     string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	Raw              bool   `json:"raw"`
	EnvironmentScope string `json:"environment_scope"`
	Description      string `json:"description"`
}

// for app version
var version string

func main() {
	var envVarsExportArgs envVarsExportArgs
	var envVarsImportArgs envVarsImportArgs

	app := &cli.App{
		Name:    "gitlab-env",
		Usage:   "A handy tool for exporting and importing Gitlab CICD environment variables",
		Version: version,
	}

	exportFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        flagGitlabToken,
			Aliases:     []string{"t"},
			Value:       "",
			Usage:       "Gitlab token",
			Destination: &envVarsExportArgs.gitlabToken,
			EnvVars:     []string{"GITLAB_TOKEN"},
		},
		&cli.StringFlag{
			Name:        flagGitlabHost,
			Aliases:     []string{"host"},
			Value:       "gitlab.com",
			Usage:       "Gitlab host domain. Example: gitlab.example.com",
			Destination: &envVarsExportArgs.gitlabHost,
			EnvVars:     []string{"GITLAB_HOST"},
		},
		&cli.StringFlag{
			Name:        flagProjectId,
			Aliases:     nil,
			Value:       "",
			Usage:       "Gitlab project ID",
			Destination: &envVarsExportArgs.projectId,
			EnvVars:     []string{"PROJECT_ID"},
		},
		&cli.StringSliceFlag{
			Name:        flagEnvScopes,
			Aliases:     []string{"s"},
			Value:       nil,
			Usage:       "Filter env vars by scopes",
			Destination: &envVarsExportArgs.scopes,
			EnvVars:     []string{"SCOPES"},
		},
		&cli.UintFlag{
			Name:        flagPageSize,
			Aliases:     nil,
			Value:       1000,
			Usage:       "Page size of return result",
			Destination: &envVarsExportArgs.pageSize,
			EnvVars:     []string{"ENV_VAR_PAGE_SIZE"},
		},
		&cli.UintFlag{
			Name:        flagPageNumber,
			Aliases:     nil,
			Value:       1,
			Usage:       "Page number of return result",
			Destination: &envVarsExportArgs.pageNumber,
			EnvVars:     []string{"ENV_VAR_PAGE_NUMBER"},
		},
		&cli.StringFlag{
			Name:        flagOutputFile,
			Aliases:     []string{"f"},
			Value:       defaultJsonFile,
			Usage:       "Path to output file",
			Destination: &envVarsExportArgs.outputFile,
			EnvVars:     []string{"OUTPUT_FILE"},
		},
	}

	importFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        flagGitlabToken,
			Aliases:     []string{"t"},
			Value:       "",
			Usage:       "Gitlab token",
			Destination: &envVarsImportArgs.gitlabToken,
			EnvVars:     []string{"GITLAB_TOKEN"},
		},
		&cli.StringFlag{
			Name:        flagGitlabHost,
			Aliases:     []string{"host"},
			Value:       "gitlab.com",
			Usage:       "Gitlab host domain. Example: gitlab.example.com",
			Destination: &envVarsImportArgs.gitlabHost,
			EnvVars:     []string{"GITLAB_HOST"},
		},
		&cli.StringFlag{
			Name:        flagProjectId,
			Aliases:     nil,
			Value:       "",
			Usage:       "Gitlab project ID",
			Destination: &envVarsImportArgs.projectId,
			EnvVars:     []string{"PROJECT_ID"},
		},
		&cli.StringSliceFlag{
			Name:        flagEnvScopes,
			Aliases:     []string{"s"},
			Value:       nil,
			Usage:       "Filter imported env vars by scopes",
			Destination: &envVarsImportArgs.scopes,
			EnvVars:     []string{"SCOPES"},
		},
		&cli.StringFlag{
			Name:        flagInputFile,
			Aliases:     []string{"f"},
			Value:       defaultJsonFile,
			Usage:       "Path to input file",
			Destination: &envVarsImportArgs.inputFile,
			EnvVars:     []string{"INPUT_FILE"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "export",
			Usage: "Export Gitlab env vars of the given project",
			Flags: exportFlags,
			Action: func(ctx *cli.Context) error {
				if err := envVarsExportArgs.validate(); err != nil {
					return fmt.Errorf("invalid argument: %w", err)
				}
				return exportProjectEnvVars(ctx.Context, envVarsExportArgs)
			},
		},
		{
			Name:  "import",
			Usage: "Import Gitlab env vars into the given project",
			Flags: importFlags,
			Action: func(ctx *cli.Context) error {
				if err := envVarsImportArgs.validate(); err != nil {
					return fmt.Errorf("invalid argument: %w", err)
				}
				return importProjectEnvVars(ctx.Context, envVarsImportArgs)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
