package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/guguducken/issue-bot/pkg/github"
	"github.com/guguducken/issue-bot/pkg/util"
	"github.com/guguducken/issue-bot/pkg/wecom"
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

	if repos == "" || label_check == "" || time_check == "" || label_skip == "" || time_skip == "" || corresponding == "" {
		panic(`repos, labels or time settings is error, please check again`)
	}

}

type Team struct {
	Team  string `json:"team"`
	Wecom string `json:"wecom"`
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
	//针对每个用户进行issue统计
	m := make(map[string]*User, 100)

	//针对每个team进行issue统计
	tm := make(map[string]*TeamIssue, 10)

	//对issue进行过期时间检测
	for i := 0; i < len(arr_repos); i++ {
		util.Info(`Start to check repo ` + arr_repos[i] + ` with milestone ` + arr_milestones[i])
		processWithRepo(arr_repos[i], arr_milestones[i], arr_label_check, arr_time_check, arr_label_skip, arr_time_skip, arr_mentions, teams, m, tm)
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
		fmt.Printf(" %s finalMess: %v\n", key, finalMess)
	}
	message := ``
	for k, v := range total {
		message += "`" + k + "`\n"
		for vk, vv := range v {
			message += vk + " total: `" + strconv.Itoa(vv) + "`\n"
		}
	}
	wecom.GenMarkdownMessage(message, arr_mentions) //.SendWecomMessage()
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

type UserIssue struct {
	Number      int         `json:"number"`
	Title       string      `json:"title"`
	Url         string      `json:"url"`
	WorkTime    util.Time_m `json:"wotktime"`
	HoliadyTime util.Time_m `json:"holidaytime"`
	StartTime   github.TimeProject
	EndTime     github.TimeProject
	Status      string
}

type Tables struct {
}

type TeamIssue struct {
	cor_issue map[string][]UserIssue
}

func processWithRepo(repo_full string, milestone string, labels []string, times []string, labels_skip, times_skip []string, arr_mentions []string, team map[string]Team, m map[string]*User, tm map[string]*TeamIssue) {
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
			la, err := github.GetLastUpdateTime(repo[0], repo[1], issues[j].Number)
			if err != nil {
				util.Error(err.Error())
			} else if la.After(*issues[j].UpdatedAt) {
				issues[j].UpdatedAt = la
			}
			t := times[i]
			skipInd := checkLabelExit(issues[j], labels_skip)
			if skipInd != -1 {
				t = times_skip[skipInd]
			}
			expTime, err := strconv.ParseInt(t, 10, 64)
			if err != nil {
				panic(`Input time_check or time_skip invalid: ` + err.Error())
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
					str := `**<font color=\"warning\">` + labels[i] + "</font>**\n"
					m[login].Repo[repo_full].labelNoticeMess[labels[i]] = str
					m[login].Repo[repo_full].labelIssueCount[labels[i]] = 0
				}
				m[login].Repo[repo_full].labelNoticeMess[labels[i]] += "-------------------------------------\n[" + issues[j].Title + "](" + issues[j].URL + ")\nUpdateAt: " + issues[j].UpdatedAt.Format(time.RFC3339) + "\nWorked: " + strconv.FormatInt(work.Days, 10) + "d-" + strconv.FormatInt(work.Hours, 10) + "h:" + strconv.FormatInt(work.Minute, 10) + "m:" + strconv.FormatInt(work.Second, 10) + "s\n"
				//对应的issue number+1
				m[login].Repo[repo_full].labelIssueCount[labels[i]]++

				fmt.Printf("m[issues[j]].login: %v\n", m[login].Login)
				fmt.Printf("m[login].Content[repo_full].LabelContent[labels[i]]: %v\n", m[login].Repo[repo_full].labelNoticeMess[labels[i]])
				fmt.Printf("m[login].Repo[repo_full].labelIssueCount[labels[i]]: %v\n", m[login].Repo[repo_full].labelIssueCount[labels[i]])
			} else {
				util.Info(`This issue time is not expired, so skip it`)
			}
			fmt.Println()
		}

		//sunday check total
		if time.Now().In(time.FixedZone(`UTC`, 0)).Weekday() == time.Sunday {
			for j := 0; j < len(issues); j++ {
				login := issues[j].Assignee.Login
				issues[j].StartTime, err = github.GetProjectTime(repo[0], repo[1], issues[j].Number, `Start Time`)
				if err != nil {
					util.Error(err.Error())
				}
				issues[j].EndTime, err = github.GetProjectTime(repo[0], repo[1], issues[j].Number, `End Time`)
				if err != nil {
					util.Error(err.Error())
				}
				issues[j].Status = github.GetProjectStatus(repo[0], repo[1], issues[j].Number)

				if _, ok := tm[team[login].Team]; !ok {
					tm[team[login].Team] = &TeamIssue{
						cor_issue: make(map[string][]UserIssue, 10),
					}
				}
				if _, ok := tm[team[login].Team].cor_issue[login]; !ok {
					tm[team[login].Team].cor_issue[login] = make([]UserIssue, 0, 10)
				}
				tm[team[login].Team].cor_issue[login] = append(tm[team[login].Team].cor_issue[login], UserIssue{
					Number:    issues[j].Number,
					Title:     issues[j].Title,
					Url:       issues[j].URL,
					StartTime: issues[j].StartTime,
					EndTime:   issues[j].EndTime,
					Status:    issues[j].Status,
				})
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

func genDetailTable(user User) string {
	return ``
}
