package github

import (
	"fmt"
	"testing"
)

func Test_GetAllPulls(t *testing.T) {
	q_pulls := NewPullsQuery(`matrixorigin`, `matrixone`, `open`, ``, ``, ``, ``)
	pulls, err := q_pulls.GetAllPulls()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("len(pulls): %v\n", len(pulls))
	fmt.Printf("pulls[0].Title: %v\n", pulls[0].Title)
	fmt.Printf("pulls[0].Number: %v\n", pulls[0].Number)
	fmt.Printf("pulls[0].HTMLURL: %v\n", pulls[0].HTMLURL)
}

func Test_UpdatePull(t *testing.T) {
	q_commit := NewCommitQuery(`mergtil`, `test_repo`, ``, ``, ``, ``, nil, nil, 1)
	commits, err := q_commit.GetCommits()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("commits[0].Sha: %v\n", commits[0].Sha)
	q_pull := NewPullsQuery(`mergtil`, `test_repo`, `open`, ``, ``, ``, ``)
	pulls, err := q_pull.GetAllPulls()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for i := 0; i < len(pulls); i++ {
		err = q_pull.UpdatePull(pulls[i].Number, pulls[i].Head.Sha)
	}
}
