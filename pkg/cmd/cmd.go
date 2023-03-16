//go:build issue_time || issue_status

package main

import (
	"encoding/json"
	"os"

	"github.com/guguducken/issue-bot/pkg/util"
)

var (
	repos          string
	label_check    string
	time_check     string
	label_skip     string
	time_skip      string
	mentions       string
	corresponding  string
	cor_milestones string
	teams          map[string]Team
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

	if repos == "" || label_check == "" || time_check == "" || label_skip == "" || time_skip == "" || corresponding == "" {
		panic(`repos, labels or time settings is error, please check again`)
	}

	teams = make(map[string]Team, 60)
	err := json.Unmarshal(([]byte)(corresponding), &teams)
	if err != nil {
		util.Error(err.Error())
	}

}

type Team struct {
	Team  string `json:"team"`
	Wecom string `json:"wecom"`
}
