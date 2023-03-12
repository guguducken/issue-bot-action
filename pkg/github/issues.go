package github

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/guguducken/auto-release/pkg/util"
)

type Q_Issus struct {
	owner     string
	repo      string
	milestone string
	state     string
	assignee  string
	creator   string
	mentioned string
	labels    []string
	sort      string
	direction string
	since     *time.Time
	per_page  int
	page      int
}

func NewIssuesQuery(owner, repo, milestone, state, assignee, creator, mentoined, sort, direction string, since *time.Time, labels []string) Q_Issus {
	q := Q_Issus{
		owner:     owner,
		repo:      repo,
		milestone: milestone,
		state:     state,
		assignee:  assignee,
		creator:   creator,
		mentioned: mentoined,
		labels:    labels,
		sort:      sort,
		direction: direction,
		since:     since,
		per_page:  100,
		page:      1,
	}
	if q.sort == "" {
		q.sort = `created`
	}
	return q
}

func (q *Q_Issus) AddPage() {
	q.page++
}

func (q *Q_Issus) GetIssuesByPage() ([]Issue, error) {
	url := githubAPI + `/repos/` + q.owner + `/` + q.repo + `/issues`
	path := `sort=` + q.sort
	if q.state != "" {
		path += `&state=` + q.state
	}
	if q.assignee != "" {
		path += `&assignee=` + q.assignee
	}
	if q.creator != "" {

		path += `&creator=` + q.creator
	}
	if q.mentioned != "" {
		path += `&mentioned=` + q.mentioned
	}
	if len(q.labels) != 0 {
		path += `&labels=` + strings.Join(q.labels, ",")
	}
	if q.milestone != "" {
		path += `&milestone=` + q.milestone
	}
	if q.direction != "" {
		path += `&direction=` + q.direction
	}
	if q.since != nil {
		path += `&since=` + q.since.Format(time.RFC3339)
	}
	if q.per_page != 0 {
		path += `&per_page=` + strconv.Itoa(q.per_page)
	}
	if q.page != 0 {
		path += `&page=` + strconv.Itoa(q.page)
	}
	// path = util.URLValid(path)
	url += `?` + path
	if q.page == 1 {
		util.Info(`The URL for get issues is: ` + url)
	}

	resp, err := get(url, token_github)
	if err != nil {
		return nil, err
	}

	issues := make([]Issue, 0, q.per_page)
	err = json.Unmarshal(resp, &issues)
	if err != nil {
		return nil, err
	}
	q.page++
	return issues, nil
}

func (q *Q_Issus) GetAllIssues() ([]Issue, error) {
	issues_all := make([]Issue, 0, 300)
	for {
		issues, err := q.GetIssuesByPage()
		if err != nil {
			return nil, err
		}
		issues_all = append(issues_all, issues...)
		q.getIssuesInfo(len(issues), len(issues_all))
		if len(issues) < q.per_page {
			break
		}
	}
	return issues_all, nil
}

func (q *Q_Issus) getIssuesInfo(num, total int) {
	str := `get issues from ` + q.owner + `/` + q.repo
	if len(q.labels) != 0 {
		str += ` with labels: ` + strings.Join(q.labels, `,`)
	}
	if q.milestone != "" {
		str += ` with milestone: ` + q.milestone
	}
	if q.assignee != "" {
		str += ` with assignee: ` + q.assignee
	}
	if q.creator != "" {
		str += ` with creator: ` + q.creator
	}
	if q.state != "" {
		str += ` with state: ` + q.state
	}

	util.Info(str + ` round ` + strconv.Itoa(q.page-1) + ` finished. round number: ` + strconv.Itoa(num) + ` total number: ` + strconv.Itoa(total))

}
