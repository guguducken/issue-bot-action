package github

import (
	"encoding/json"
	"strconv"

	"github.com/guguducken/issue-bot/pkg/util"
)

type Q_Pulls struct {
	owner     string
	repo      string
	state     string
	head      string
	base      string
	sort      string
	direction string
	per_page  int
	page      int
}

func NewPullsQuery(owner, repo, state, head, base, sort, direction string) Q_Pulls {
	q := Q_Pulls{
		owner:     owner,
		repo:      repo,
		state:     state,
		head:      head,
		base:      base,
		sort:      sort,
		direction: direction,
		per_page:  100,
		page:      1,
	}
	return q
}

func (q *Q_Pulls) AddPage() {
	q.page++
}

func (q *Q_Pulls) GetPullsByPage() (pulls []PullRequest, err error) {
	url := githubRestAPI + `/repos/` + q.owner + `/` + q.repo + `/pulls`

	path := `per_page=` + strconv.Itoa(q.per_page)
	if q.state != "" {
		path += `&state=` + q.state
	}
	if q.head != "" {
		path += `&head=` + q.head
	}
	if q.base != "" {
		path += `&base=` + q.base
	}
	if q.sort != "" {
		path += `&sort=` + q.sort
	}
	if q.direction != "" {
		path += `&direction=` + q.direction
	}
	path += `&page=` + strconv.Itoa(q.page)

	url += `?` + path
	if q.page == 1 {
		util.Info(`The URL for get pulls is: ` + url)
	}

	resp, err := get(url, token_github)
	if err != nil {
		return nil, err
	}

	pulls = make([]PullRequest, 0, q.per_page)
	err = json.Unmarshal(resp, &pulls)
	if err != nil {
		return nil, err
	}
	q.page++
	return pulls, nil
}

func (q *Q_Pulls) GetAllPulls() ([]PullRequest, error) {
	pulls_all := make([]PullRequest, 0, 300)
	for {
		pulls, err := q.GetPullsByPage()
		if err != nil {
			return nil, err
		}
		pulls_all = append(pulls_all, pulls...)
		q.getPullInfo(len(pulls), len(pulls_all))
		if len(pulls) < q.per_page {
			break
		}
	}
	return pulls_all, nil
}

func (q *Q_Pulls) getPullInfo(num, total int) {
	str := `get pulls from ` + q.owner + `/` + q.repo

	if q.state != "" {
		str += ` with state: ` + q.state
	}
	if q.head != "" {
		str += ` with head: ` + q.head
	}
	if q.base != "" {
		str += ` with base: ` + q.base
	}

	util.Info(str + ` round ` + strconv.Itoa(q.page-1) + ` finished. round number: ` + strconv.Itoa(num) + ` total number: ` + strconv.Itoa(total))
}

func (q *Q_Pulls) UpdatePull(number int, head_sha string) error {
	url := githubRestAPI + `/repos/` + q.owner + `/` + q.repo + `/pulls/` + strconv.Itoa(number) + `/update-branch`
	body, err := json.Marshal(struct {
		HeadSha string `json:"expected_head_sha"`
	}{
		HeadSha: head_sha,
	})
	if err != nil {
		return nil
	}
	_, err = put(url, token_github, string(body))
	if err == nil {
		util.Info(`update pull request:` + strconv.Itoa(number) + ` success`)
	}
	return err
}
