package github

import (
	"fmt"
	"testing"
	"time"

	"github.com/guguducken/auto-release/pkg/util"
)

func Test_getIssue(t *testing.T) {
	q_issue_list := NewIssueListQuery(`matrixorigin`, `matrixone`, ``, `open`, ``, ``, ``, ``, ``, ``, []string{"kind/bug", "severity/s1"})
	issue, err := q_issue_list.GetAllIssues()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("issue.Title: %v\n", issue[0].Title)
	fmt.Printf("issue.CreatedAt: %v\n", issue[0].CreatedAt)
	fmt.Printf("issue.UpdatedAt: %v\n", issue[0].UpdatedAt)
	fmt.Printf("issue.CommentsURL: %v\n", issue[0].CommentsURL)
}

func Test_expired(t *testing.T) {
	q_issue_list := NewIssueListQuery(`matrixorigin`, `matrixone`, ``, `open`, ``, ``, ``, ``, ``, ``, []string{"kind/bug", "severity/s1"})
	issue, err := q_issue_list.GetAllIssues()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for i := 0; i < 10; i++ {
		work, holiday, err := util.GetPassedTimeWithoutWeekend(issue[i].CreatedAt)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		loc := time.FixedZone(`UTC`, 8*3600)
		fmt.Printf("issue[i].CreatedAt.In(loc): %v\n", issue[i].CreatedAt.In(loc))
		fmt.Printf("time.Now(): %v\n", time.Now())
		fmt.Printf("issue[%d] work: %v\n", issue[i].Number, work/3600000)
		fmt.Printf("issue[%d] holiday: %v\n", issue[i].Number, holiday/3600000)
		fmt.Println("-----------------------------------------")
	}
}
