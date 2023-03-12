package github

import (
	"fmt"
	"testing"
)

func Test_GetCommits(t *testing.T) {
	q_commit := NewCommitQuery(`mergtil`, `test_repo`, ``, ``, ``, ``, nil, nil, 1)
	commits, err := q_commit.GetCommits()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("commits[0].Sha: %v\n", commits[0].Sha)
}
