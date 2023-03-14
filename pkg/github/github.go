package github

var (
	githubRestAPI    string = `https://api.github.com`
	githubGraphqlAPI string = `https://api.github.com/graphql`
	token_github     string = `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`
)

func init() {
	// githubRestAPI = os.Getenv(`GITHUB_API_URL`)
	// githubGraphqlAPI = os.Getenv(`GITHUB_API_URL`)
	// token_github = os.Getenv(`INPUT_TOKEN_ACTION`)

	// if len(token_github) == 0 || len(githubRestAPI) == 0 {
	// 	panic(`GitHub Settings Invalid: invalid GitHub API or Token, please check again`)
	// }
}
