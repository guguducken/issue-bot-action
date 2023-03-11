package github

import (
	"fmt"
	"testing"
)

func Test_getIssue(t *testing.T) {
	q_issue_list := NewIssueListQuery(`matrixorigin`, `matrixone`, ``, `open`, ``, ``, ``, ``, ``, ``, ``)
	var issue Issue
	for true {
		issues, _ := q_issue_list.GetIssuesByPage()
		if len(issues) == 0 {
			break
		}
		flag := false
		for i := 0; i < len(issues); i++ {
			if issues[i].Number == 7431 {
				issue = issues[i]
				flag = true
				break
			}
		}
		if flag {
			break
		}
		q_issue_list.AddPage()
	}
	fmt.Printf("issue.Title: %v\n", issue.Title)
	fmt.Printf("issue.CreatedAt: %v\n", issue.CreatedAt)
	fmt.Printf("issue.UpdatedAt: %v\n", issue.UpdatedAt)
	fmt.Printf("issue.CommentsURL: %v\n", issue.CommentsURL)
}
