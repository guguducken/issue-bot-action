//go:build issue_status
// +build issue_status

package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/guguducken/issue-bot/pkg/github"
	"github.com/guguducken/issue-bot/pkg/util"
)

var (
	now time.Time = time.Now().In(time.FixedZone(`UTC`, 0))
)

type UserIssue struct {
	Number      int         `json:"number"`
	Title       string      `json:"title"`
	Url         string      `json:"url"`
	WorkTime    util.Time_m `json:"wotktime"`
	HoliadyTime util.Time_m `json:"holidaytime"`
	StartTime   github.TimeProject
	EndTime     github.TimeProject
	Status      map[string]int
}

type Tables struct {
}

type TeamIssue struct {
	cor_issue map[string][]UserIssue
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
	teams := make(map[string]Team, 60)
	err := json.Unmarshal(([]byte)(corresponding), &teams)
	if err != nil {
		util.Error(err.Error())
	}

	//针对每个team进行issue统计
	tm := make(map[string]*TeamIssue, 10)

}

func isNewOpen(issue github.Issue) bool {

}
