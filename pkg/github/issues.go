package github

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/guguducken/issue-bot/pkg/util"
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
func (q *Q_Issus) SubPage() {
	q.page--
}

func (q *Q_Issus) GetIssuesByPage() ([]Issue, error) {
	url := githubRestAPI + `/repos/` + q.owner + `/` + q.repo + `/issues`
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
	if q.labels != nil && len(q.labels) != 0 {
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
	try_number := 1
	for {
		issues, err := q.GetIssuesByPage()
		if err == util.TimeoutErr {
			if try_number > 10 {
				return nil, util.NertworkErr
			}
			util.Warning(`Get issues timeout, and we will retry, the retry number is ` + strconv.Itoa(try_number))
			q.SubPage()
			try_number++
			continue
		}
		if err != nil {
			return nil, err
		}
		try_number = 0
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

func GetLastUpdateTime(owner, repo string, issue Issue) time.Time {
	prLast, err := GetRelatePRLastUpdateTime(owner, repo, issue.Number)
	if err != nil {
		util.Error(err.Error())
		return time.Date(1, 1, 1, 0, 0, 0, 0, time.FixedZone(`UTC`, 0))
	}
	commitLast, err := GetRelatedCommitLastUpdateTime(owner, repo, issue.CommentsURL, issue.Assignee.Login)
	if err != nil {
		util.Error(err.Error())
		return time.Date(1, 1, 1, 0, 0, 0, 0, time.FixedZone(`UTC`, 0))
	}

	if commitLast.After(prLast) {
		return commitLast
	}
	return prLast
}

func GetRelatedCommitLastUpdateTime(owner, repo string, comments_url string, assignee string) (la time.Time, err error) {
	la = time.Date(1, 1, 1, 1, 0, 0, 0, time.FixedZone(`UTC`, 0))
	resp, err := http.Get(comments_url)
	if err != nil {
		return
	}
	if resp.Header.Get(`x-ratelimit-remaining`) == `0` {
		util.Error(`The github resource have been consumed, the reset UTC time is: ` + resp.Header.Get(`x-ratelimit-reset`))
		panic(`The github resource have been consumed, the reset UTC time is: ` + resp.Header.Get(`x-ratelimit-reset`))
	}
	reply, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	comments := make([]Comment, 0, 10)
	err = json.Unmarshal(reply, &comments)
	if err != nil {
		return
	}
	for i := len(comments) - 1; i >= 0; i-- {
		if comments[i].User.Login == assignee {
			return *comments[i].UpdatedAt, nil
		}
	}
	return
}

func GetRelatePRLastUpdateTime(owner, repo string, number int) (la time.Time, err error) {
	type Temp struct {
		Data struct {
			Repository struct {
				Issue struct {
					TimelineItems struct {
						Edges []struct {
							Cursor string `json:"cursor"`
							Node   struct {
								Source *struct {
									CreatedAt time.Time `json:"createdAt"`
									UpdatedAt time.Time `json:"updatedAt"`
									Number    int       `json:"number"`
								} `json:"source"`
							} `json:"node,omitempty"`
						} `json:"edges"`
					} `json:"timelineItems"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}

	cursor := `null`
	per_page := 20
	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(last: ` + strconv.Itoa(per_page) + `, itemTypes: CROSS_REFERENCED_EVENT) { edges { node { ... on CrossReferencedEvent { source { ... on PullRequest { createdAt updatedAt number } } } } cursor } } } }}"}`
	reply, err := post(githubGraphqlAPI, token_github, []byte(query))
	if err != nil {
		return
	}
	t := Temp{}
	err = json.Unmarshal(reply, &t)
	if err != nil {
		return
	}

	la = time.Date(1, 1, 1, 0, 0, 0, 0, time.FixedZone(`UTC`, 0))
	edge := t.Data.Repository.Issue.TimelineItems.Edges
	for len(edge) != 0 && cursor != edge[0].Cursor {
		cursor = edge[0].Cursor
		for i := 0; i < len(edge); i++ {
			if edge[i].Node.Source == nil {
				continue
			}
			if edge[i].Node.Source.UpdatedAt.After(la) {
				la = edge[i].Node.Source.UpdatedAt
			}
		}
		query = `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(last: ` + strconv.Itoa(per_page) + `, before: \"` + cursor + `\", itemTypes: CROSS_REFERENCED_EVENT) { edges { node { ... on CrossReferencedEvent { source { ... on PullRequest { createdAt updatedAt number } } } } cursor } } } }}"}`
		reply, err = post(githubGraphqlAPI, token_github, []byte(query))
		if err != nil {
			return
		}
		err = json.Unmarshal(reply, &t)
		if err != nil {
			return
		}
		edge = t.Data.Repository.Issue.TimelineItems.Edges
	}
	return
}

func GetProjectTime(owner, repo string, number int, timeChose string) (tp *TimeProject, err error) {
	type Temp struct {
		Data struct {
			Repository struct {
				Issue struct {
					ProjectItems struct {
						Nodes []struct {
							FieldValueByName *TimeProject `json:"fieldValueByName"`
						} `json:"nodes,omitempty"`
					} `json:"projectItems"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}
	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { projectItems(includeArchived: false, first: 5) { nodes { fieldValueByName(name: \"` + timeChose + `\") { ... on ProjectV2ItemFieldDateValue { id updatedAt createdAt date } } } } } }}"}`
	req, err := http.NewRequest(`POST`, githubGraphqlAPI, strings.NewReader(query))
	if err != nil {
		return
	}
	req.Header.Set(`Authorization`, `Bearer `+token_github)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	reply, _ := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	t := Temp{}
	err = json.Unmarshal(reply, &t)
	if err != nil {
		return
	}
	for i := 0; i < len(t.Data.Repository.Issue.ProjectItems.Nodes); i++ {
		if t.Data.Repository.Issue.ProjectItems.Nodes[i].FieldValueByName == nil {
			continue
		}
		tp = t.Data.Repository.Issue.ProjectItems.Nodes[i].FieldValueByName
	}
	if tp != nil {
		util.Info(`Get issue ` + strconv.Itoa(number) + ` ` + timeChose + `: ` + tp.Date)
	} else {
		util.Warning(`This issue ` + strconv.Itoa(number) + ` are not added to a project or set ` + timeChose + `, so set it to default black time`)
		return &TimeProject{
			Date:      "",
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}, nil
	}
	return
}

func GetProjectStatus(owner, repo string, number int) (str string) {
	type Temp struct {
		Data struct {
			Repository struct {
				Issue struct {
					Number       int `json:"number"`
					ProjectItems struct {
						Nodes []struct {
							FieldValueByName *struct {
								ID        string    `json:"id"`
								Name      string    `json:"name"`
								UpdatedAt time.Time `json:"updatedAt"`
							} `json:"fieldValueByName,omitempty"`
						} `json:"nodes"`
					} `json:"projectItems"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}

	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { number projectItems(includeArchived: false, first: 5) { nodes { fieldValueByName(name: \"Status\") { ... on ProjectV2ItemFieldSingleSelectValue { id name updatedAt } } } } } }}"}`
	reply, err := post(githubGraphqlAPI, token_github, []byte(query))
	if err != nil {
		return
	}
	t := Temp{}
	err = json.Unmarshal(reply, &t)
	if err != nil {
		return
	}
	str = ``
	la := time.Date(1, 1, 1, 0, 0, 0, 0, time.FixedZone(`UTC`, 0))
	for i := 0; i < len(t.Data.Repository.Issue.ProjectItems.Nodes); i++ {
		if t.Data.Repository.Issue.ProjectItems.Nodes[i].FieldValueByName != nil {
			if t.Data.Repository.Issue.ProjectItems.Nodes[i].FieldValueByName.UpdatedAt.After(la) {
				la = t.Data.Repository.Issue.ProjectItems.Nodes[i].FieldValueByName.UpdatedAt
				str = t.Data.Repository.Issue.ProjectItems.Nodes[i].FieldValueByName.Name
			}
		}
	}
	if str == "" {
		util.Warning(`this issue do not set progress or add it to a project, so set it to default Todo`)
		str = "Todo"
	}
	util.Info(`Get issue ` + strconv.Itoa(number) + ` Status: ` + str)
	return str
}

func GetUnassignLogins(owner, repo string, number int) (logins []string) {
	per_page := 10
	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(itemTypes: UNASSIGNED_EVENT, first: ` + strconv.Itoa(per_page) + `) { edges { node { ... on UnassignedEvent { createdAt user { login } } } } filteredCount } } }}"}`

	type Temp struct {
		Data struct {
			Repository struct {
				Issue struct {
					TimelineItems struct {
						Edges []struct {
							Node struct {
								User struct {
									Login string `json:"login"`
								} `json:"user"`
							} `json:"node"`
						} `json:"edges"`
						FilteredCount int `json:"filteredCount"`
						PageInfo      struct {
							EndCursor string `json:"endCursor"`
						} `json:"pageInfo"`
					} `json:"timelineItems"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}

	logins = make([]string, 0, 10)
	m_logins := make(map[string]struct{}, 5)
	reply, err := post(githubGraphqlAPI, token_github, []byte(query))
	if err != nil {
		return
	}
	t := Temp{}
	err = json.Unmarshal(reply, &t)
	if err != nil {
		return
	}
	cursor := ``
	for {
		for i := 0; i < len(t.Data.Repository.Issue.TimelineItems.Edges); i++ {
			if _, ok := m_logins[t.Data.Repository.Issue.TimelineItems.Edges[i].Node.User.Login]; !ok {
				logins = append(logins, t.Data.Repository.Issue.TimelineItems.Edges[i].Node.User.Login)
				m_logins[t.Data.Repository.Issue.TimelineItems.Edges[i].Node.User.Login] = struct{}{}
			}
		}
		if t.Data.Repository.Issue.TimelineItems.FilteredCount < per_page {
			break
		}
		cursor = t.Data.Repository.Issue.TimelineItems.PageInfo.EndCursor
		query = `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(itemTypes: UNASSIGNED_EVENT, first: ` + strconv.Itoa(per_page) + `, after: \"` + cursor + `\") { edges { node { ... on UnassignedEvent { createdAt user { login } } } } filteredCount } } }}"}`
		reply, err = post(githubGraphqlAPI, token_github, []byte(query))
		if err != nil {
			return
		}
		err = json.Unmarshal(reply, &t)
		if err != nil {
			return
		}
	}
	return logins
}

func GetAssignDetail(owner, repo string, eventURL string) ([]Event, error) {
	now := 100
	per_page, page := 100, 1
	events_all := make([]Event, 0, 100)
	url := eventURL + `?per_page=` + strconv.Itoa(per_page) + `&page=` + strconv.Itoa(page)
	reply, err := get(url, token_github)
	if err != nil {
		return events_all, err
	}
	err = json.Unmarshal(reply, &events_all)
	if err != nil {
		return events_all, err
	}
	now = len(events_all)

	for per_page == now {
		page++
		url = eventURL + `?per_page=` + strconv.Itoa(per_page) + `&page=` + strconv.Itoa(page)
		reply, err = get(url, token_github)
		if err != nil {
			return events_all, err
		}
		events := make([]Event, 0, 100)
		err = json.Unmarshal(reply, &events)
		if err != nil {
			break
		}
		events_all = append(events_all, events...)
		now = len(events)
	}
	return events_all, nil
}
