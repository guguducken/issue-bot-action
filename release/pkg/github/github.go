package github

import (
	"encoding/json"
	"os"
	"strconv"
)

var (
	githubAPI    string
	token_github string
	repo_full    string
)

func init() {
	githubAPI = os.Getenv(`GITHUB_API_URL`)
	token_github = os.Getenv(`INPUT_TOKEN-ACTION`)
}

type Q_Issue_List struct {
	owner     string
	repo      string
	milestone string
	state     string
	assignee  string
	creator   string
	mentioned string
	labels    string
	sort      string
	direction string
	since     string
	per_page  int
	page      int
}

func NewIssueListQuery(owner, repo, milestone, state, assignee, creator, mentoined, labels, sort, direction, since string) *Q_Issue_List {
	return &Q_Issue_List{
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
}

func (q *Q_Issue_List) AddPage() {
	q.page++
}

func (q *Q_Issue_List) GetIssuesByPage() ([]Issue, error) {
	url := githubAPI + `/repos/` + q.owner + `/` + q.repo + `/issues`
	flag := true
	if q.milestone != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `milestone=` + q.milestone
	}
	if q.state != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `state=` + q.state
	}
	if q.assignee != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `assignee=` + q.assignee
	}
	if q.creator != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `creator=` + q.creator
	}
	if q.mentioned != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `mentioned=` + q.mentioned
	}
	if q.labels != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `labels=` + q.labels
	}
	if q.sort != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `sort=` + q.sort
	}
	if q.direction != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `direction=` + q.direction
	}
	if q.since != "" {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `since=` + q.since
	}
	if q.per_page != 0 {
		if flag {
			url += `?`
			flag = false
		} else {
			url += `&`
		}
		url += `per_page=` + strconv.Itoa(q.per_page)
	}
	if q.page != 0 {
		if flag {
			url += `?`
		} else {
			url += `&`
		}
		url += `page=` + strconv.Itoa(q.page)
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
	return issues, nil
}

func (q *Q_Issue_List) GetAllIssues() ([]Issue, error) {
	issues_all := make([]Issue, 0, 300)
	for {
		issues, err := q.GetIssuesByPage()
		if err != nil {
			return nil, err
		}
		if len(issues) == 0 {
			break
		}
		q.AddPage()
		issues_all = append(issues_all, issues...)
	}
	return issues_all, nil
}
