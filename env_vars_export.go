package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
)

type envVarsExportArgs struct {
	gitlabToken string
	gitlabHost  string
	projectId   string
	scopes      cli.StringSlice
	pageSize    uint
	pageNumber  uint
	exportAll   bool
	outputFile  string
}

func (args envVarsExportArgs) validate() error {
	if args.gitlabHost == "" {
		return fmt.Errorf("gitlab host cannot be empty")
	}
	if args.gitlabToken == "" {
		return fmt.Errorf("gitlab token cannot be empty")
	}
	if args.projectId == "" {
		return fmt.Errorf("gitlab project id cannot be empty")
	}
	if args.pageNumber < 1 {
		return fmt.Errorf("page number should be greater than or equal 1")
	}
	if args.pageSize < 1 {
		return fmt.Errorf("page size should be greater than or equal 1")
	}
	if args.pageSize > gitlabMaxPageSize {
		return fmt.Errorf("page size should be less than or equal %d", gitlabMaxPageSize)
	}
	if args.outputFile == "" {
		return fmt.Errorf("output file cannot be empty")
	}
	return nil
}

func exportProjectEnvVars(ctx context.Context, args envVarsExportArgs) error {
	logrus.Info("run project env vars export: args=", args)
	var result []GitlabEnvVar

	if args.exportAll {
		logrus.Info("export-all flag turned on, start with pageNumber=1, pageSize=", gitlabMaxPageSize)
		args.pageNumber = 1
		args.pageSize = gitlabMaxPageSize
	}

	for {
		requestUrl := buildEnvVarsExportUrl(args)
		logrus.Info("env vars url=", requestUrl)

		req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
		if err != nil {
			return fmt.Errorf("could not create request: %s", err)
		}
		req.Header.Set("PRIVATE-TOKEN", args.gitlabToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error making http request: %s\n", err)
		}
		logrus.Info("response status code=", res.StatusCode)

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body) // response body is []byte
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		var envVars []GitlabEnvVar
		if err := json.Unmarshal(body, &envVars); err != nil { // Parse []byte to the go struct pointer
			return fmt.Errorf("cannot unmarshal json: %v", err)
		}

		logrus.Info("number of env vars: ", len(envVars))

		for _, item := range envVars {
			if len(args.scopes.Value()) == 0 || slices.Contains(args.scopes.Value(), item.EnvironmentScope) {
				logrus.Info("found env=", item.Key, ", scope=", item.EnvironmentScope)
				result = append(result, item)
			}
		}
		if args.exportAll == false {
			break
		}

		// Check the 'x-next-page' header
		nextPage := res.Header.Get("x-next-page")
		if nextPage == "" {
			logrus.Info("no more pages to fetch. exiting loop")
			break
		} else {
			logrus.Info("fetching nextPage=", nextPage)
			args.pageNumber = StringToUint(nextPage)
		}
	}

	logrus.Info("writing result to output file: ", args.outputFile)
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result to json: %v", err)
	}
	err = os.WriteFile(args.outputFile, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed write to output file: %v", err)
	}
	return nil
}

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}

func buildEnvVarsExportUrl(args envVarsExportArgs) string {
	return fmt.Sprintf("https://%s/api/v4/projects/%s/variables?page=%d&per_page=%d", args.gitlabHost, args.projectId, args.pageNumber, args.pageSize)
}
