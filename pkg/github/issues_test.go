package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/guguducken/issue-bot/pkg/util"
)

func Test_getIssue(t *testing.T) {
	q_issue_list := NewIssuesQuery(`matrixorigin`, `matrixone`, ``, `open`, ``, ``, ``, ``, ``, nil, nil)
	issue, err := q_issue_list.GetAllIssues()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for i := 0; i < len(issue); i++ {
		if issue[i].PullRequest == nil {
			fmt.Printf("issue.Title: %v\n", issue[i].Title)
			fmt.Printf("issue.CreatedAt: %v\n", issue[i].CreatedAt)
			fmt.Printf("issue.UpdatedAt: %v\n", issue[i].UpdatedAt)
			fmt.Printf("issue.CommentsURL: %v\n", issue[i].CommentsURL)
			t, _ := json.Marshal(issue[i])
			fmt.Printf("issue[0].Assignee.Login: %v\n", string(t))
			break
		} else {
			fmt.Println(`Skip ` + strconv.Itoa(issue[i].Number))
		}
	}
}

func Test_expired(t *testing.T) {
	q_issue_list := NewIssuesQuery(`matrixorigin`, `matrixone`, ``, `open`, ``, ``, ``, ``, ``, nil, []string{"kind/bug", "severity/s1"})
	issue, err := q_issue_list.GetAllIssues()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for i := 0; i < 10; i++ {
		work, holiday, err := util.GetPassedTimeWithoutWeekend(*issue[i].CreatedAt)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		loc := time.FixedZone(`UTC`, 8*3600)
		fmt.Printf("issue[i].CreatedAt.In(loc): %v\n", issue[i].CreatedAt.In(loc))
		fmt.Printf("time.Now(): %v\n", time.Now())
		fmt.Printf("issue[%d] work: %v\n", issue[i].Number, work)
		fmt.Printf("issue[%d] holiday: %v\n", issue[i].Number, holiday)
		fmt.Println("-----------------------------------------")
	}
}

func Test_graphql(t *testing.T) {
	url := `https://api.github.com/graphql`

	body := `{"query": "{ repository(name: \"matrixone\", owner: \"matrixorigin\") { issue(number: 8314) { projectItems(first: 10) { nodes { fieldValueByName(name: \"End Time\") { ... on ProjectV2ItemFieldDateValue { id updatedAt createdAt date } } } } } }}"}`
	req, _ := http.NewRequest(`POST`, url, strings.NewReader(body))
	req.Header.Set(`Authorization`, `Bearer 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`)
	resp, _ := http.DefaultClient.Do(req)
	ans, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("string(ans): %v\n", string(ans))
}

// func Test_GetLastUpdateTime(t *testing.T) {
// 	last, err := GetLastUpdateTime(`matrixorigin`, `matrixone`)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("last: %v\n", last)
// }
// func Test_GetRelatePRLastUpdateTime(t *testing.T) {
// 	last, err := GetRelatePRLastUpdateTime(`matrixorigin`, `matrixone`, 8420)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("last: %v\n", last)
// }
// func Test_GetRelatedCommitLastUpdateTime(t *testing.T) {
// 	last, err := GetRelatedCommitLastUpdateTime(`matrixorigin`, `matrixone`, 8420)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("last: %v\n", last)
// }

func Test_GetProjectTime(t *testing.T) {
	fmt.Println(GetProjectTime(`matrixorigin`, `matrixone`, 8216, `End Time`))
}

func Test_GetStatus(t *testing.T) {
	fmt.Printf("GetProjectStatus(`matrixorigin`, `matrixone`, 8216): %v\n", GetProjectStatus(`matrixorigin`, `matrixone`, 8216))
}

func Test_GetUnassignLogins(t *testing.T) {
	fmt.Printf("GetUnassignLogins(`matrixorigin`, `matrixone`, 7278): %v\n", GetUnassignLogins(`matrixorigin`, `matrixone`, 7278))
}

func Test_GetAssignDetail(t *testing.T) {
	events, err := GetAssignDetail(`matrixorigin`, `matrixone`, `https://api.github.com/repos/matrixorigin/matrixone/issues/7278/events`)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("len(events): %v\n", len(events))
	for i := len(events) - 1; i >= 0; i-- {
		fmt.Printf("events[%d].Event: %v\n", i, events[i].Event)
		if events[i].Event == `assigned` {
			fmt.Println(`Assigner: ` + events[i].Assigner.Login)
			fmt.Println(`Assignee: ` + events[i].Assignee.Login)
		}
	}
}
