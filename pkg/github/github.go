package github

import "os"

var (
	githubRestAPI    string
	githubGraphqlAPI string
	token_github     string
)

func init() {
	githubRestAPI = os.Getenv(`GITHUB_API_URL`)
	githubGraphqlAPI = os.Getenv(`GITHUB_GRAPHQL_URL`)
	token_github = os.Getenv(`INPUT_TOKEN_ACTION`)

	if len(token_github) == 0 || len(githubRestAPI) == 0 {
		panic(`GitHub Settings Invalid: invalid GitHub API or Token, please check again`)
	}
}
