package github

import (
	"os"
)

var (
	githubAPI    string
	token_github string
)

func init() {
	githubAPI = os.Getenv(`GITHUB_API_URL`)
	token_github = os.Getenv(`INPUT_TOKEN_ACTION`)

	if len(token_github) == 0 || len(githubAPI) == 0 {
		panic(`GitHub Settings Invalid: invalid GitHub API or Token, please check again`)
	}
}
