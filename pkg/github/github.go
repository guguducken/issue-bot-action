package github

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/guguducken/auto-release/pkg/util"
)

var (
	githubAPI    string
	token_github string
)

func init() {
	githubAPI = os.Getenv(`GITHUB_API_URL`)
	token_github = os.Getenv(`INPUT_TOKEN_ACTION`)
}

type Q_Issue_List struct {
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
	since     string
	per_page  int
	page      int
}

func NewIssueListQuery(owner, repo, milestone, state, assignee, creator, mentoined, sort, direction, since string, labels []string) Q_Issue_List {
	q := Q_Issue_List{
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

func (q Q_Issue_List) AddPage() {
	q.page++
}

func (q Q_Issue_List) GetIssuesByPage() ([]Issue, error) {
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
	if q.since != "" {
		path += `&since=` + q.since
	}
	if q.per_page != 0 {
		path += `&per_page=` + strconv.Itoa(q.per_page)
	}
	if q.page != 0 {
		path += `&page=` + strconv.Itoa(q.page)
	}
	path = util.URLValid(path)
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
	return issues, nil
}

func (q Q_Issue_List) GetAllIssues() ([]Issue, error) {
	issues_all := make([]Issue, 0, 300)
	for {
		issues, err := q.GetIssuesByPage()
		if err != nil {
			return nil, err
		}
		issues_all = append(issues_all, issues...)
		util.Info(`get issues round ` + strconv.Itoa(q.page) + ` finished. round number: ` + strconv.Itoa(len(issues)) + ` total number: ` + strconv.Itoa(len(issues_all)))
		if len(issues) < q.per_page {
			break
		}
		q.AddPage()
	}
	return issues_all, nil
}
