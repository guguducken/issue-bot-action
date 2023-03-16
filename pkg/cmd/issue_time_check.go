//go:build issue_time

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/guguducken/issue-bot/pkg/github"
	"github.com/guguducken/issue-bot/pkg/util"
	"github.com/guguducken/issue-bot/pkg/wecom"
)

// set env variable

func main() {
	//generate all array
	arr_repos := strings.Split(repos, `,`)
	arr_label_check := strings.Split(label_check, `,`)
	arr_time_check := strings.Split(time_check, `,`)
	arr_label_skip := strings.Split(label_skip, `,`)
	arr_time_skip := strings.Split(time_skip, `,`)
	arr_mentions := strings.Split(mentions, `,`)
	arr_milestones := strings.Split(cor_milestones, `,`)

	//针对每个用户进行issue统计
	m := make(map[string]*User, 100)

	//对issue进行过期时间检测
	for i := 0; i < len(arr_repos); i++ {
		util.Info(`Start to check repo ` + arr_repos[i] + ` with milestone ` + arr_milestones[i])
		checkIssueExpiredTimeWithRepo(arr_repos[i], arr_milestones[i], arr_label_check, arr_time_check, arr_label_skip, arr_time_skip, arr_mentions, teams, m)
	}
	//expired Time Notice
	total := make(map[string]map[string]int, len(arr_repos))
	for i := 0; i < len(arr_repos); i++ {
		total[arr_repos[i]] = make(map[string]int)
	}
	for key, value := range m {
		message := ``
		for k, v := range value.Repo {
			message += "`" + k + "`\n"
			for i := 0; i < len(arr_label_check); i++ {
				message += v.labelNoticeMess[arr_label_check[i]]
				total[k][arr_label_check[i]] += v.labelIssueCount[arr_label_check[i]]
			}
		}
		message += "Total: " + strconv.Itoa(value.Total)
		message += "Assignee: @<" + value.Wecom + ">"
		finalMess, _ := json.Marshal(wecom.GenMarkdownMessage(message, []string{})) //.SendWecomMessage()
		fmt.Printf(" %s finalMess: %v\n", key, string(finalMess))
	}
	message := ``
	for k, v := range total {
		message += "`" + k + "`\n"
		for vk, vv := range v {
			message += vk + " total: `" + strconv.Itoa(vv) + "`\n"
		}
	}
	wecom.GenMarkdownMessage(message, arr_mentions)                               //.SendWecomMessage()
	finalMess, _ := json.Marshal(wecom.GenMarkdownMessage(message, arr_mentions)) //.SendWecomMessage()
	fmt.Printf("Final total finalMess: %v\n", string(finalMess))
}

type User struct {
	Login    string `json:"login"`
	Wecom    string `json:"wecom"`
	Repo     map[string]*Content
	Total    int    `json:"total"`
	TeamName string `json:"teamname"`
}

type Content struct {
	labelNoticeMess map[string]string
	labelIssueCount map[string]int
}

func checkIssueExpiredTimeWithRepo(repo_full string, milestone string, labels []string, times []string, labels_skip, times_skip []string, arr_mentions []string, team map[string]Team, m map[string]*User) {
	repo := strings.Split(repo_full, `/`)
	for i := 0; i < len(labels); i++ {
		q_issue := github.NewIssuesQuery(repo[0], repo[1], milestone, `open`, ``, ``, ``, ``, ``, nil, []string{labels[i]})
		issues, err := q_issue.GetAllIssues()
		if err != nil {
			util.Error(err.Error())
			continue
		}
		for j := 0; j < len(issues); j++ {
			login := issues[j].Assignee.Login

			util.Info(`Start to check issue ->>>> number: ` + strconv.Itoa(issues[j].Number) + ` title: ` + issues[j].Title)
			t := times[i]
			skipInd := checkLabelExit(issues[j], labels_skip)
			if skipInd != -1 {
				t = times_skip[skipInd]
			}
			expTime, err := strconv.ParseInt(t, 10, 64)
			if err != nil {
				panic(`Input time_check or time_skip invalid: ` + err.Error())
			}

			la := github.GetLastUpdateTime(repo[0], repo[1], issues[j])
			if la.After(*issues[j].CreatedAt) {
				issues[j].UpdatedAt = &la
			} else {
				issues[j].UpdatedAt = issues[j].CreatedAt
			}

			work, holiday, err := util.GetPassedTimeWithoutWeekend(*issues[j].UpdatedAt)
			if err != nil {
				util.Error(err.Error())
			}
			util.Info(`Check: work--> ` + strconv.FormatInt(work.Days, 10) + `d-` + strconv.FormatInt(work.Hours, 10) + `h:` + strconv.FormatInt(work.Minute, 10) + `m:` + strconv.FormatInt(work.Second, 10) + `s == ` + strconv.FormatInt(work.TotalMillSecond, 10) + `ms || holiday--> ` + strconv.FormatInt(holiday.Days, 10) + `d-` + strconv.FormatInt(holiday.Hours, 10) + `h:` + strconv.FormatInt(holiday.Minute, 10) + `m:` + strconv.FormatInt(holiday.Second, 10) + `s` + strconv.FormatInt(holiday.MillSecond, 10) + ` == ` + strconv.FormatInt(holiday.TotalMillSecond, 10) + `ms || target--> ` + strconv.FormatInt(expTime, 10) + `h`)

			if work.TotalMillSecond >= expTime*3600000 {
				util.Warning(`>>> ` + repo_full + ` Warning issue: ` + strconv.Itoa(issues[j].Number) + ` - ` + issues[j].Title + ` update time: ` + issues[j].UpdatedAt.Format(time.RFC3339))
				if _, ok := m[login]; !ok {
					m[login] = &User{
						Login:    login,
						Wecom:    team[login].Wecom,
						TeamName: team[login].Team,
						Total:    0,
						Repo:     make(map[string]*Content, 5),
					}
				}
				m[login].Total++

				if _, ok := m[login].Repo[repo_full]; !ok {
					m[login].Repo[repo_full] = &Content{
						labelNoticeMess: make(map[string]string, len(labels)),
						labelIssueCount: make(map[string]int, len(labels)),
					}
				}
				if _, ok := m[login].Repo[repo_full].labelNoticeMess[labels[i]]; !ok {
					str := "**<font color=\"warning\">" + labels[i] + "</font>**\n"
					m[login].Repo[repo_full].labelNoticeMess[labels[i]] = str
					m[login].Repo[repo_full].labelIssueCount[labels[i]] = 0
				}
				m[login].Repo[repo_full].labelNoticeMess[labels[i]] += "-------------------------------------\n[" + issues[j].Title + "](" + issues[j].URL + ")\nUpdateAt: " + issues[j].UpdatedAt.Format(time.RFC3339) + "\nWorked: " + strconv.FormatInt(work.Days, 10) + "d-" + strconv.FormatInt(work.Hours, 10) + "h:" + strconv.FormatInt(work.Minute, 10) + "m:" + strconv.FormatInt(work.Second, 10) + "s\n"
				//对应的issue number+1
				m[login].Repo[repo_full].labelIssueCount[labels[i]]++
			} else {
				util.Info(`The last update time of this issue time is` + issues[j].UpdatedAt.Format(time.RFC3339) + ` which is not expired, so skip it`)
			}
			fmt.Println()
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
