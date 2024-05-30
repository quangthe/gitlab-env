# Gitlab Env

A handy tool to export and import a large number of CICD environment variables.

Command usage:
```shell
go install github.com/quangthe/gitlab-env@latest

gitlab-env -h
gitlab-env export -h
gitlab-env import -h
```

## Export Gitlab project env vars

```shell
gitlab-env export \
  --gitlab-token "your-gitlab-token" \
  --gitlab-host "gitlab.example.com"  \
  --project-id "55" \
  --output-file gitlab.json
```
> By default, the command will retrieve up to `1000` env vars. Use `--page-size` flag to change the maximum number of exported env vars. 

Filter by environment scopes `--env-scopes`:
```shell
gitlab-env export \
  --gitlab-token "your-gitlab-token" \
  --gitlab-host "gitlab.example.com"  \
  --project-id "55" \
  --env-scopes "*,dev,prod" \
  --output-file gitlab.json
```

The gitlab.json file will look like this
```json
[
  {
    "variable_type": "env_var",
    "key": "FORM_URL",
    "value": "jdbc:postgresql:\\/\\/shared-database:5432\\/forms",
    "protected": true,
    "masked": false,
    "raw": true,
    "environment_scope": "*",
    "description": ""
  },
  {
    "variable_type": "env_var",
    "key": "POLL_URL",
    "value": "jdbc:postgresql:\\/\\/shared-database:5432\\/polls",
    "protected": true,
    "masked": false,
    "raw": true,
    "environment_scope": "*",
    "description": ""
  }
]
```

## Import env vars into Gitlab project

```shell
gitlab-env import \
  --gitlab-token "your-gitlab-token" \
  --gitlab-host "gitlab.example.com"  \
  --project-id "56" \
  --input-file gitlab.json
```

> By default, the `import` command does not filter on the `Environment Scope`. Use `--env-scopes` flag to only import env vars with selected scopes. 
Example: `--env-scopes "*,dev,prod"` will import env vars with scope `*`, `dev` and `prod`. 