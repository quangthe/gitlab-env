package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
)

type envVarsImportArgs struct {
	gitlabToken string
	gitlabHost  string
	projectId   string
	scopes      cli.StringSlice
	inputFile   string
}

func (args envVarsImportArgs) validate() error {
	if args.gitlabHost == "" {
		return fmt.Errorf("gitlab host cannot be empty")
	}
	if args.gitlabToken == "" {
		return fmt.Errorf("gitlab token cannot be empty")
	}
	if args.projectId == "" {
		return fmt.Errorf("gitlab project id cannot be empty")
	}
	if args.inputFile == "" {
		return fmt.Errorf("input file cannot be empty")
	}
	return nil
}

func importProjectEnvVars(ctx context.Context, args envVarsImportArgs) error {
	logrus.Info("run project env vars import: args=", args)

	input, err := os.Open(args.inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer input.Close()

	bytes, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	var envVars []GitlabEnvVar
	if err := json.Unmarshal(bytes, &envVars); err != nil { // Parse []byte to the go struct pointer
		return fmt.Errorf("cannot unmarshal json: %v", err)
	}
	logrus.Info("number of env vars: ", len(envVars))
	for _, item := range envVars {
		if len(args.scopes.Value()) == 0 || slices.Contains(args.scopes.Value(), item.EnvironmentScope) {
			logrus.Info("processing ", item.Key, ", scope=", item.EnvironmentScope)
			err := createEnvVar(args, item)
			if err != nil {
				logrus.Info("try to update env var: ", item.Key)
				updateEnvVar(args, item)
			}
		}
	}
	return nil
}

func createEnvVar(args envVarsImportArgs, item GitlabEnvVar) error {
	requestUrl := buildCreateEnvVarUrl(args)
	logrus.Info("create env var: url=", requestUrl)

	data := url.Values{}
	data.Set("key", item.Key)
	data.Set("value", item.Value)
	data.Set("description", item.Description)
	data.Set("environment_scope", item.EnvironmentScope)
	data.Set("protected", strconv.FormatBool(item.Protected))
	data.Set("masked", strconv.FormatBool(item.Masked))
	data.Set("raw", strconv.FormatBool(item.Raw))
	data.Set("variable_type", item.VariableType)

	req, err := http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}
	req.Header.Set("PRIVATE-TOKEN", args.gitlabToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making http request: %s\n", err)
	}
	logrus.Info("response status code=", res.StatusCode)
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(res.Body)
		logrus.Info("response body: ", string(b))
		return fmt.Errorf("failed to create env var")
	}

	return nil
}

func updateEnvVar(args envVarsImportArgs, item GitlabEnvVar) error {
	requestUrl := buildUpdateEnvVarUrl(args, item.Key, item.EnvironmentScope)
	logrus.Info("update env var: url=", requestUrl)

	data := url.Values{}
	data.Set("key", item.Key)
	data.Set("value", item.Value)
	data.Set("description", item.Description)
	data.Set("environment_scope", item.EnvironmentScope)
	data.Set("protected", strconv.FormatBool(item.Protected))
	data.Set("masked", strconv.FormatBool(item.Masked))
	data.Set("raw", strconv.FormatBool(item.Raw))
	data.Set("variable_type", item.VariableType)

	req, err := http.NewRequest(http.MethodPut, requestUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}
	req.Header.Set("PRIVATE-TOKEN", args.gitlabToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making http request: %s\n", err)
	}
	logrus.Info("response status code=", res.StatusCode)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		logrus.Info("response body: ", string(b))
	}
	return nil
}

func buildCreateEnvVarUrl(args envVarsImportArgs) string {
	return fmt.Sprintf("https://%s/api/v4/projects/%s/variables", args.gitlabHost, args.projectId)
}

func buildUpdateEnvVarUrl(args envVarsImportArgs, envKey string, scope string) string {
	return fmt.Sprintf("https://%s/api/v4/projects/%s/variables/%s?filter[environment_scope]=%s", args.gitlabHost, args.projectId, envKey, scope)
}
