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
	m := make(map[string]*user, 100)
	//query for different repo
	for i := 0; i < len(arr_repos); i++ {
		util.Info(`Start to check repo ` + arr_repos[i] + ` with milestone ` + arr_milestones[i])
		processWithRepo(arr_repos[i], arr_milestones[i], arr_label_check, arr_time_check, arr_label_skip, arr_time_skip, arr_mentions, teams, m)
	}
	//get all issues
}

type user struct {
	login   string
	wecom   string
	content map[string]*struct {
		labelContent map[string]*string
	}
	total  int
	issues map[string]*struct {
		number      int
		title       string
		url         string
		workTime    util.Time_m
		holiadyTime util.Time_m
	}
	teamName string
	tables   struct {
	}
}

func processWithRepo(repo_full string, milestone string, labels []string, times []string, labels_skip, times_skip []string, arr_mentions []string, team map[string]Team, m map[string]*user) {
	repo := strings.Split(repo_full, `/`)
	for i := 0; i < len(labels); i++ {
		q_issue := github.NewIssuesQuery(repo[0], repo[1], milestone, `open`, ``, ``, ``, ``, ``, nil, []string{labels[i]})
		issues, err := q_issue.GetAllIssues()
		if err != nil {
			util.Error(err.Error())
			continue
		}
		for j := 0; j < len(issues); j++ {
			util.Info(`Start to check issue ->>>> number: ` + strconv.Itoa(issues[j].Number) + ` title: ` + issues[j].Title)
			la, err := github.GetLastUpdateTime(repo[0], repo[1], issues[j].Number)
			if la.After(*issues[j].UpdatedAt) {
				issues[j].UpdatedAt = la
			}
			if err != nil {
				util.Error(err.Error())
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
			util.Info(`Pass: work--> ` + strconv.FormatInt(work.Days, 10) + `d-` + strconv.FormatInt(work.Hours, 10) + `h:` + strconv.FormatInt(work.Minute, 10) + `m:` + strconv.FormatInt(work.Second, 10) + `s == ` + strconv.FormatInt(work.TotalMillSecond, 10) + `ms || holiday--> ` + strconv.FormatInt(holiday.Days, 10) + `d-` + strconv.FormatInt(holiday.Hours, 10) + `h:` + strconv.FormatInt(holiday.Minute, 10) + `m:` + strconv.FormatInt(holiday.Second, 10) + `s` + strconv.FormatInt(holiday.MillSecond, 10) + ` == ` + strconv.FormatInt(holiday.TotalMillSecond, 10) + `ms || target--> ` + strconv.FormatInt(expTime, 10) + `h`)

			if work.TotalMillSecond >= expTime*3600000 {
				util.Warning(`>>> ` + repo_full + ` Warning issue: ` + strconv.Itoa(issues[j].Number) + ` - ` + issues[j].Title + ` update time: ` + issues[j].UpdatedAt.Format(time.RFC3339))
				if _, ok := m[issues[j].Assignee.Login]; !ok {
					m[issues[j].Assignee.Login] = &user{
						login:    issues[j].Assignee.Login,
						wecom:    team[issues[j].Assignee.Login].Wecom,
						teamName: team[issues[j].Assignee.Login].Team,
						total:    0,
						content: make(map[string]*struct {
							labelContent map[string]*string
						}, 0),
						issues: make(map[string]*struct {
							number      int
							title       string
							url         string
							workTime    util.Time_m
							holiadyTime util.Time_m
						}),
						tables: struct {
						}{},
					}
				}
				m[issues[j].Assignee.Login].total++

				if _, ok := m[issues[j].Assignee.Login].content[repo[1]]; !ok {
					m[issues[j].Assignee.Login].content[repo[1]] = &struct{ labelContent map[string]*string }{
						labelContent: make(map[string]*string, 0),
					}
				}
				if m[issues[j].Assignee.Login].content[repo[1]] == nil {
					m[issues[j].Assignee.Login].content[repo[1]] = &struct {
						labelContent map[string]*string
					}{
						labelContent: make(map[string]*string, 10),
					}
				}
				if _, ok := m[issues[j].Assignee.Login].content[repo[1]].labelContent[labels[i]]; !ok {
					str := `**<font color=\"warning\">` + labels[i] + `</font>**\n`
					m[issues[j].Assignee.Login].content[repo[1]].labelContent[labels[i]] = &str
				}
				*m[issues[j].Assignee.Login].content[repo[1]].labelContent[labels[i]] += `-------------------------------------\n[` + issues[j].Title + `](` + issues[j].URL + `)\nUpdateAt: ` + issues[j].UpdatedAt.Format(time.RFC3339) + `\nWorked: ` + strconv.FormatInt(work.Days, 10) + `d-` + strconv.FormatInt(work.Hours, 10) + `h:` + strconv.FormatInt(work.Minute, 10) + `m:` + strconv.FormatInt(work.Second, 10) + `s\n`
				fmt.Printf("m[issues[j]].login: %v\n", m[issues[j].Assignee.Login].login)
			} else {
				util.Info(`This issue time is not expired, so skip it`)
			}
		}

		//sunday check total
		if time.Now().In(time.FixedZone(`UTC`, 0)).Weekday() == time.Sunday {
			for j := 0; j < len(issues); j++ {
				issues[j].StartTime, err = github.GetProjectTime(repo[0], repo[1], issues[j].Number, `Start Time`)
				if err != nil {
					util.Error(err.Error())
				}
				issues[j].EndTime, err = github.GetProjectTime(repo[0], repo[1], issues[j].Number, `End Time`)
				if err != nil {
					util.Error(err.Error())
				}
				issues[j].Status = github.GetProjectStatus(repo[0], repo[1], issues[j].Number)
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
