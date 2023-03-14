package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/guguducken/issue-bot/pkg/github"
	"github.com/guguducken/issue-bot/pkg/util"
)

// set env variable
var (
	repos          string
	label_check    string
	time_check     string
	label_skip     string
	time_skip      string
	mentions       string
	corresponding  string
	cor_milestones string
)

func init() {
	repos = os.Getenv(`INPUT_REPOS`)
	label_check = os.Getenv(`INPUT_LABEL_CHECK`)
	time_check = os.Getenv(`INPUT_TIME_CHECK`)
	label_skip = os.Getenv(`INPUT_LABEL_SKIP`)
	time_skip = os.Getenv(`INPUT_TIME_SKIP`)
	cor_milestones = os.Getenv(`INPUT_COR_MILESTONES`)

	mentions = os.Getenv(`INPUT_MENTIONS`)
	corresponding = os.Getenv(`INPUT_CORRESPONDING`)

	// if repos == "" || label_check == "" || time_check == "" || label_skip == "" || time_skip == "" || corresponding == "" {
	// 	panic(`repos, labels or time settings is error, please check again`)
	// }

}

type Team struct {
	Name    string `json:"name"`
	Peoples []struct {
		Login  string `json:"login"`
		Weocom string `json:"wecom"`
	} `json:"peoples"`
}

func main() {
	//generate all array
	arr_repos := strings.Split(repos, `,`)
	arr_label_check := strings.Split(label_check, `,`)
	arr_time_check := strings.Split(time_check, `,`)
	arr_label_skip := strings.Split(label_skip, `,`)
	arr_time_skip := strings.Split(time_skip, `,`)
	arr_mentions := strings.Split(mentions, `,`)
	arr_milestones := strings.Split(cor_milestones, `,`)
	teams := make([]Team, 0, 10)
	err := json.Unmarshal(([]byte)(corresponding), &teams)
	if err != nil {
		util.Error(err.Error())
	}
	//query for different repo
	for i := 0; i < len(arr_repos); i++ {

	}
	//get all issues
}

func processWithRepo(repo_full string, milestone string, labels []string, times []string, labels_skip, times_skip []string) {
	type user struct {
		login   string
		wecom   string
		content string
		issues  map[string]struct {
			id          int
			title       string
			workTime    util.Time_m
			holiadyTime util.Time_m
		}
		team   Team
		tables struct {
		}
	}
	m := make(map[string]user, len(labels))
	repo := strings.Split(repo_full, `/`)
	for i := 0; i < len(labels); i++ {
		q_issue := github.NewIssuesQuery(repo[0], repo[1], milestone, `open`, ``, ``, ``, ``, ``, nil, []string{labels[i]})
		issues, err := q_issue.GetAllIssues()
		if err != nil {
			util.Error(err.Error())
			continue
		}
		for j := 0; j < len(issues); j++ {
			expTime := times[j]
			skipTimeInd := checkLabelExit(issues[j], labels_skip)
			if skipTimeInd != -1 {
				expTime = times_skip[skipTimeInd]
			}

		}

	}
}

func checkLabelExit(issue github.Issue, labels []string) int {
	for i := 0; i < len(labels); i++ {
		for j := 0; j < len(issue.Labels); j++ {
			if labels[i] == issue.Labels[j].Name {
				return i
			}
		}
	}
	return -1
}
