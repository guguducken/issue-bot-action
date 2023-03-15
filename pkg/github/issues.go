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

func GetLastUpdateTime(owner, repo string, number int) (last *time.Time, err error) {
	prLast, err := GetRelatePRLastUpdateTime(owner, repo, number)
	if err != nil {
		return
	}
	commitLast, err := GetRelatedCommitLastUpdateTime(owner, repo, number)
	if commitLast.After(prLast) {
		last = &commitLast
		return
	}
	last = &prLast
	return
}

func GetRelatedCommitLastUpdateTime(owner, repo string, number int) (la time.Time, err error) {
	type Temp struct {
		Data struct {
			Repository struct {
				Issue struct {
					TimelineItems struct {
						Edges []struct {
							Cursor string `json:"cursor"`
							Node   *struct {
								ID        string    `json:"id"`
								CreatedAt time.Time `json:"createdAt"`
								UpdatedAt time.Time `json:"updatedAt"`
							} `json:"node,omitempty"`
						} `json:"edges"`
					} `json:"timelineItems"`
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
				} `json:"issue"`
			} `json:"repository"`
		} `json:"data"`
	}
	cursor := `null`
	per_page := 20
	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(last: ` + strconv.Itoa(per_page) + `) { edges { node { ... on IssueComment { id createdAt updatedAt } } cursor } } createdAt updatedAt } }}"}`
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
			if edge[i].Node == nil || edge[i].Node.ID == `` {
				continue
			}
			if edge[i].Node.CreatedAt.After(edge[i].Node.UpdatedAt) {
				edge[i].Node.UpdatedAt = edge[i].Node.CreatedAt
			}
			if edge[i].Node.UpdatedAt.After(la) {
				la = edge[i].Node.UpdatedAt
			}
		}
		query = `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(last: ` + strconv.Itoa(per_page) + `, before: \"` + cursor + `\") { edges { node { ... on IssueComment { id createdAt updatedAt } } cursor } } createdAt updatedAt } }}"}`
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
	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(last: ` + strconv.Itoa(per_page) + `) { edges { node { ... on CrossReferencedEvent { source { ... on PullRequest { createdAt updatedAt number } } } } cursor } } } }}"}`
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
			if edge[i].Node.Source == nil || edge[i].Cursor == `` {
				continue
			}
			if edge[i].Node.Source.CreatedAt.After(edge[i].Node.Source.UpdatedAt) {
				edge[i].Node.Source.UpdatedAt = edge[i].Node.Source.CreatedAt
			}
			if edge[i].Node.Source.UpdatedAt.After(la) {
				la = edge[i].Node.Source.UpdatedAt
			}
		}
		query = `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { timelineItems(last: ` + strconv.Itoa(per_page) + `, before: \"` + cursor + `\") { edges { node { ... on CrossReferencedEvent { source { ... on PullRequest { createdAt updatedAt number } } } } cursor } } } }}"}`
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

func GetProjectTime(owner, repo string, number int, timeChose string) (tp TimeProject, err error) {
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
	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { projectItems(first: 5) { nodes { fieldValueByName(name: \"` + timeChose + `\") { ... on ProjectV2ItemFieldDateValue { id updatedAt createdAt date } } } } } }}"}`
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
		tp = *t.Data.Repository.Issue.ProjectItems.Nodes[i].FieldValueByName
	}
	util.Info(`Get issue ` + strconv.Itoa(number) + ` ` + timeChose + `: ` + tp.Date)
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

	query := `{"query":"{ repository(name: \"` + repo + `\", owner: \"` + owner + `\") { issue(number: ` + strconv.Itoa(number) + `) { number projectItems(first: 5) { nodes { fieldValueByName(name: \"Status\") { ... on ProjectV2ItemFieldSingleSelectValue { id name updatedAt } } } } } }}"}`
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
	if str == `` {
		util.Error(`Fail to get Status, please check again`)
	}
	util.Info(`Get issue ` + strconv.Itoa(number) + ` Status: ` + str)
	return str
}
