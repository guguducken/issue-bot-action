package github

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Q_Commits struct {
	owner     string
	repo      string
	sha       string
	path      string
	author    string
	committer string
	since     *time.Time
	until     *time.Time
	per_page  int
	page      int
}

func (q *Q_Pulls) GetBranchFromCommit() (branch string) {
	return ""
}

func NewCommitQuery(owner, repo, sha, path, author, committer string, since, until *time.Time, per_page int) Q_Commits {
	q := Q_Commits{
		owner:     owner,
		repo:      repo,
		sha:       sha,
		path:      path,
		author:    author,
		committer: committer,
		since:     since,
		until:     until,
		page:      1,
	}
	if per_page > 0 {
		q.per_page = per_page
	} else {
		q.per_page = 30
	}
	return q
}

func (q *Q_Commits) GetCommits() ([]Commit, error) {
	url := githubAPI + `/repos/` + q.owner + `/` + q.repo + `/commits`
	path := `per_page=` + strconv.Itoa(q.per_page)
	if q.sha != "" {
		path += `&sha=` + q.sha
	}
	if q.author != "" {
		path += `&author=` + q.author
	}
	if q.committer != "" {
		path += `&committer=` + q.committer
	}
	path += `&page=` + strconv.Itoa(q.page)
	url += `?` + path

	commits := make([]Commit, 0, q.per_page)
	resp, err := http.Get(url)
	if err != nil {
		return commits, err
	}
	reply, err := io.ReadAll(resp.Body)
	if err != nil {
		return commits, err
	}
	err = json.Unmarshal(reply, &commits)
	if err == nil {
		q.page++
	}
	return commits, err
}
