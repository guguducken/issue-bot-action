//go:build issue_status
// +build issue_status

package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/guguducken/issue-bot/pkg/github"
	"github.com/guguducken/issue-bot/pkg/util"
	"github.com/guguducken/issue-bot/pkg/wecom"
)

var (
	loc                *time.Location = time.FixedZone(`UTC`, 8*3600)
	locUTC             *time.Location = time.FixedZone(`UTC`, 0)
	now                time.Time      = time.Now().In(loc)
	oneWeekagoStart    time.Time
	nowUTC             time.Time
	oneWeekagoStartUTC time.Time
)

func init() {
	now = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 23, 59, 59, 0, loc)
	oneWeekagoStart = time.Date(now.Year(), now.Month(), now.Day()-6, 0, 0, 0, 0, loc)

	nowUTC = now.In(locUTC)
	oneWeekagoStartUTC = oneWeekagoStart.In(locUTC)
}

type people struct {
	login       string
	todo        []github.Issue
	in_progress []github.Issue
	holding     []github.Issue
	unplanned   []github.Issue
	new_open    []github.Issue
	finished    []github.Issue
	delayed     []github.Issue
	adjusted    []github.Issue
	no_endtime  []github.Issue
}

type TeamIssue struct {
	cor_issue map[string]*people
}

func main() {
	//generate all array
	arr_repos := strings.Split(repos, `,`)
	arr_milestones := strings.Split(cor_milestones, `,`)

	//针对每个team进行issue统计
	tm := make(map[string]*TeamIssue, 10)

	//获取每个issue的状态
	for i := 0; i < len(arr_repos); i++ {
		repo := strings.Split(arr_repos[i], `/`)
		checkIssueWithRepo(repo[0], repo[1], arr_milestones[i], tm)
	}

	totals := make([]int, 9)
	people_toal := 0
	totalMessage := "> 当前统计范围为: V0.8.0(matrixone repo)和MO-V0.8.0(MO-Cloud repo)两个milestone下的所有issue。\n\n> 固定更新周期，每周一早上8点统计输出\n\n| Team | People | Todo | In Progress | Holding | Unplanned | New Open | Finished | Delayed | Adjusted | NO EndTime |\n| :----: | :------: | :----: | :-----------: | :-------: | :---------: | :--------: | :--------: | :-------: | :--------: | :----------: |\n"
	detailMessage := "| Team | People | Progress | Number | Name |\n| :----: | :------: | :--------: | :------: | :----: |\n"
	noEndTimeMessage := ``
	for key, val := range tm {
		for k, v := range val.cor_issue {
			totalMessage += "|" + key + "|" + k + "|" + strconv.Itoa(len(v.todo)) + "|" + strconv.Itoa(len(v.in_progress)) + "|" + strconv.Itoa(len(v.holding)) + "|" + strconv.Itoa(len(v.unplanned)) + "|" + strconv.Itoa(len(v.new_open)) + "|" + strconv.Itoa(len(v.finished)) + "|" + strconv.Itoa(len(v.delayed)) + "|" + strconv.Itoa(len(v.adjusted)) + "|" + strconv.Itoa(len(v.no_endtime)) + "|\n"
			people_toal++
			totals[0] += len(v.todo)
			totals[1] += len(v.in_progress)
			totals[2] += len(v.holding)
			totals[3] += len(v.unplanned)
			totals[4] += len(v.new_open)
			totals[5] += len(v.finished)
			totals[6] += len(v.delayed)
			totals[7] += len(v.adjusted)
			totals[8] += len(v.no_endtime)
			// for i := 0; i < len(v.todo); i++ {
			// 	detailMessage += "|" + key + "|" + k + "|todo|" + strconv.Itoa(v.todo[i].Number) + "|[" + v.todo[i].Title + "](" + v.todo[i].URL + ")|\n"
			// }
			// for i := 0; i < len(v.in_progress); i++ {
			// 	detailMessage += "|" + key + "|" + k + "|in_progress|" + strconv.Itoa(v.in_progress[i].Number) + "|[" + v.in_progress[i].Title + "](" + v.in_progress[i].HTMLURL + ")|\n"
			// }
			// for i := 0; i < len(v.holding); i++ {
			// 	detailMessage += "|" + key + "|" + k + "|holding|" + strconv.Itoa(v.holding[i].Number) + "|[" + v.holding[i].Title + "](" + v.holding[i].HTMLURL + ")|\n"
			// }
			// for i := 0; i < len(v.unplanned); i++ {
			// 	detailMessage += "|" + key + "|" + k + "|unplanned|" + strconv.Itoa(v.unplanned[i].Number) + "|[" + v.unplanned[i].Title + "](" + v.unplanned[i].HTMLURL + ")|\n"
			// }
			for i := 0; i < len(v.new_open); i++ {
				detailMessage += "|" + key + "|" + k + "|new_open|" + strconv.Itoa(v.new_open[i].Number) + "|[" + v.new_open[i].Title + "](" + v.new_open[i].HTMLURL + ")|\n"
			}
			for i := 0; i < len(v.finished); i++ {
				detailMessage += "|" + key + "|" + k + "|finished|" + strconv.Itoa(v.finished[i].Number) + "|[" + v.finished[i].Title + "](" + v.finished[i].HTMLURL + ")|\n"
			}
			for i := 0; i < len(v.delayed); i++ {
				detailMessage += "|" + key + "|" + k + "|delayed|" + strconv.Itoa(v.delayed[i].Number) + "|[" + v.delayed[i].Title + "](" + v.delayed[i].HTMLURL + ")|\n"
			}
			for i := 0; i < len(v.adjusted); i++ {
				detailMessage += "|" + key + "|" + k + "|adjusted|" + strconv.Itoa(v.adjusted[i].Number) + "|[" + v.adjusted[i].Title + "](" + v.adjusted[i].HTMLURL + ")|\n"
			}
			for i := 0; i < len(v.no_endtime); i++ {
				noEndTimeMessage += "|" + key + "|" + k + "|no_endtime|" + strconv.Itoa(v.no_endtime[i].Number) + "|[" + v.no_endtime[i].Title + "](" + v.no_endtime[i].HTMLURL + ")|\n"
			}
		}
	}
	totalMessage += "|Total|" + strconv.Itoa(people_toal) + "|" + strconv.Itoa(totals[0]) + "|" + strconv.Itoa(totals[1]) + "|" + strconv.Itoa(totals[2]) + "|" + strconv.Itoa(totals[3]) + "|" + strconv.Itoa(totals[4]) + "|" + strconv.Itoa(totals[5]) + "|" + strconv.Itoa(totals[6]) + "|" + strconv.Itoa(totals[7]) + "|" + strconv.Itoa(totals[8]) + "|\n"
	total_send, err := wecom.GenFileMessage(totalMessage, oneWeekagoStart.Format(time.RFC3339)+"__"+now.Format(time.RFC3339)+"_total_check.md")
	if err != nil {
		util.Error(err.Error())
		util.Warning(totalMessage)
	} else if total_send.SendWecomMessage() != nil {
		util.Error(err.Error())
		util.Warning(totalMessage)
	}

	detailMessage += noEndTimeMessage
	detail_send, err := wecom.GenFileMessage(detailMessage, oneWeekagoStart.Format(time.RFC3339)+"__"+now.Format(time.RFC3339)+"_detail_check.md")
	if err != nil {
		util.Error(err.Error())
		util.Warning(detailMessage)
	} else if detail_send.SendWecomMessage() != nil {
		util.Error(err.Error())
		util.Warning(detailMessage)
	}
}

func checkIssueWithRepo(owner, repo string, milestone string, tm map[string]*TeamIssue) {
	q_issue := github.NewIssuesQuery(owner, repo, milestone, `all`, ``, ``, ``, ``, `created`, nil, nil)
	issues, err := q_issue.GetAllIssues()
	if err != nil {
		return
	}

	for j := 0; j < len(issues); j++ {
		if issues[j].PullRequest != nil {
			util.Warning(`this issue ` + strconv.Itoa(issues[j].Number) + ` -- title: ` + issues[j].Title + ` is pull request, so skip it`)
			continue
		}
		if issues[j].Assignee == nil {
			util.Warning(`this issue ` + strconv.Itoa(issues[j].Number) + ` -- title: ` + issues[j].Title + ` don't have an assignee, so check it's creator`)
			issues[j].Assignee = issues[j].User
		}
		login := issues[j].Assignee.Login
		issues[j].EndTime, err = github.GetProjectTime(owner, repo, issues[j].Number, `End Time`)
		if err != nil {
			util.Error(err.Error())
		}
		issues[j].Status = github.GetProjectStatus(owner, repo, issues[j].Number)
		if err != nil {
			util.Error(err.Error())
		}

		teamInit(tm, login)
		totalCheck(issues[j], tm[teams[login].Team].cor_issue[login])
		if issues[j].CreatedAt.After(oneWeekagoStartUTC) && issues[j].CreatedAt.Before(nowUTC) {
			detailCheck(owner, repo, issues[j], tm)
		}
	}
	return

}

func totalCheck(issue github.Issue, p *people) {
	switch issue.Status {
	case `Todo`:
		if len(issue.EndTime.Date) == 0 {
			p.unplanned = append(p.unplanned, issue) //unplanned
			return
		}
		p.todo = append(p.todo, issue) //todo
		return
	case `In Progress`:
		// if len(issue.EndTime.Date) == 0 {
		// }
		p.in_progress = append(p.in_progress, issue) //in_progress
		return
	case `Holding`:
		p.holding = append(p.holding, issue) //holding
		return
	}
}

func detailCheck(owner, repo string, issue github.Issue, tm map[string]*TeamIssue) {
	login := issue.Assignee.Login

	if issue.State == `open` {
		if issue.CreatedAt.After(oneWeekagoStartUTC) && issue.User.Login != login { //新创建且创建者不为现在的assignee的即为当前用户的new_open
			tm[teams[login].Team].cor_issue[login].new_open = append(tm[teams[login].Team].cor_issue[login].new_open, issue)
		}
	}

	//对finished的情况进行检查
	finished := finishedCheck(owner, repo, issue, tm)

	if !finished {
		if issue.EndTime.Date != `` { //只有在endtime有数据且没有finished时候进行delay和adjust检查
			da := strings.Split(issue.EndTime.Date, `-`)
			year, _ := strconv.Atoi(da[0])
			month, _ := strconv.Atoi(da[1])
			day, _ := strconv.Atoi(da[2])
			time_da := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
			if time_da.Before(now) { //设定的结束时间在现在之前且没有finished
				tm[teams[login].Team].cor_issue[login].delayed = append(tm[teams[login].Team].cor_issue[login].delayed, issue)
			} else if issue.EndTime.UpdatedAt.After(oneWeekagoStartUTC) && issue.EndTime.UpdatedAt.Before(nowUTC) { //上周期间发生过update
				tm[teams[login].Team].cor_issue[login].adjusted = append(tm[teams[login].Team].cor_issue[login].adjusted, issue)
			}
		} else {
			tm[teams[login].Team].cor_issue[login].no_endtime = append(tm[teams[login].Team].cor_issue[login].no_endtime, issue)
		}
	}
}

func finishedCheck(owner, repo string, issue github.Issue, tm map[string]*TeamIssue) bool {
	//返还给创建者并且由创建者关闭的issue可标记为finished
	if issue.User.Login == issue.Assignee.Login && issue.State == `close` {
		assig, err := github.GetAssignDetail(owner, repo, issue.EventsURL)
		if err != nil {
			util.Error(err.Error())
			return false
		}
		for i := 0; i < len(assig); i++ {
			if assig[i].Event == `assigned` && assig[i].Assignee.Login == issue.User.Login {
				teamInit(tm, assig[i].Assigner.Login)
				tm[teams[assig[i].Assigner.Login].Team].cor_issue[assig[i].Assigner.Login].finished = append(tm[teams[assig[i].Assigner.Login].Team].cor_issue[assig[i].Assigner.Login].finished, issue)
				return true
			}
		}
	}

	//检查此issue状态，满足状态则将其加入到现在assign的finished数组中，并返回此issue已经完成，不能进行后续delay或adjust检查
	login := issue.Assignee.Login
	if issue.Status == `Done` && login != issue.User.Login { //该用户名(不为测试)下issue 状态为Done的即为finished
		tm[teams[login].Team].cor_issue[login].finished = append(tm[teams[login].Team].cor_issue[login].finished, issue)
		return true
	}

	return false
}

func teamInit(tm map[string]*TeamIssue, login string) {
	//检查对应的team是否已经初始化
	if _, ok := tm[teams[login].Team]; !ok {
		tm[teams[login].Team] = &TeamIssue{
			cor_issue: make(map[string]*people, 15),
		}
	}

	//检查对应的people是否已经初始化
	if _, ok := tm[teams[login].Team].cor_issue[login]; !ok {
		tm[teams[login].Team].cor_issue[login] = &people{
			login:       login,
			todo:        make([]github.Issue, 0, 5),
			in_progress: make([]github.Issue, 0, 5),
			holding:     make([]github.Issue, 0, 5),
			unplanned:   make([]github.Issue, 0, 5),
			new_open:    make([]github.Issue, 0, 5),
			finished:    make([]github.Issue, 0, 5),
			delayed:     make([]github.Issue, 0, 5),
			adjusted:    make([]github.Issue, 0, 5),
		}
	}
}
