package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/guguducken/issue-bot/pkg/util"
)

func Test_getIssue(t *testing.T) {
	q_issue_list := NewIssuesQuery(`matrixorigin`, `matrixone`, ``, `open`, ``, ``, ``, ``, ``, nil, nil)
	issue, err := q_issue_list.GetAllIssues()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for i := 0; i < len(issue); i++ {
		if issue[i].PullRequest == nil {
			fmt.Printf("issue.Title: %v\n", issue[i].Title)
			fmt.Printf("issue.CreatedAt: %v\n", issue[i].CreatedAt)
			fmt.Printf("issue.UpdatedAt: %v\n", issue[i].UpdatedAt)
			fmt.Printf("issue.CommentsURL: %v\n", issue[i].CommentsURL)
			t, _ := json.Marshal(issue[i])
			fmt.Printf("issue[0].Assignee.Login: %v\n", string(t))
			break
		} else {
			fmt.Println(`Skip ` + strconv.Itoa(issue[i].Number))
		}
	}
}

func Test_expired(t *testing.T) {
	q_issue_list := NewIssuesQuery(`matrixorigin`, `matrixone`, ``, `open`, ``, ``, ``, ``, ``, nil, []string{"kind/bug", "severity/s1"})
	issue, err := q_issue_list.GetAllIssues()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for i := 0; i < 10; i++ {
		work, holiday, err := util.GetPassedTimeWithoutWeekend(*issue[i].CreatedAt)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		loc := time.FixedZone(`UTC`, 8*3600)
		fmt.Printf("issue[i].CreatedAt.In(loc): %v\n", issue[i].CreatedAt.In(loc))
		fmt.Printf("time.Now(): %v\n", time.Now())
		fmt.Printf("issue[%d] work: %v\n", issue[i].Number, work/3600000)
		fmt.Printf("issue[%d] holiday: %v\n", issue[i].Number, holiday/3600000)
		fmt.Println("-----------------------------------------")
	}
}

func Test_graphql(t *testing.T) {
	url := `https://api.github.com/graphql`
	// body := `{"query":"{ repository(name: \"matrixone\", owner: \"matrixorigin\") { projectV2(number: 13) { items(first: 10) { edges { node { fieldValues(first: 10) { nodes { ... on ProjectV2ItemFieldSingleSelectValue { id name optionId } } } content { ... on Issue { number title repository { name } } } } } } } }}"}`
	// body := `{"query":"{ repository(name: \"matrixone\", owner: \"matrixorigin\") { projectV2(number: 13) { items(first: 10) { nodes { content { ... on Issue { id state number url title } } fieldValueByName(name: \"Status\") { ... on ProjectV2ItemFieldSingleSelectValue { id name item { fieldValueByName(name: \"Status\") { ... on ProjectV2ItemFieldSingleSelectValue { id name } } } } } } } } }}"}`
	body := `{"query":"query{repository(name: \"matrixone\", owner: \"matrixorigin\") { issue(number: 10) { updatedAt createdAt state timelineItems(last: 30) { nodes { ... on IssueComment { id createdAt updatedAt } ... on CrossReferencedEvent { id source { ... on PullRequest { id createdAt updatedAt }}}}}}}}"}`
	req, _ := http.NewRequest(`POST`, url, strings.NewReader(body))
	req.Header.Set(`Authorization`, `Bearer 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`)
	resp, _ := http.DefaultClient.Do(req)
	ans, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	project := ProjectV2{}
	json.Unmarshal(ans, &project)
	fmt.Printf("string(ans): %v\n", string(ans))
}

func Test_GetLastUpdateTime(t *testing.T) {
	GetLastUpdateTime(`matrixorigin/matrixone`, 8440)
}

//   {
// 	repository(name: "matrixone", owner: "matrixorigin") {
// 	  projectV2(number: 13) {
// 		items(first: 10) {
// 		  nodes {
// 			content {
// 			  ... on Issue {
// 				id
// 				state
// 				title
// 				resourcePath
// 				url
// 				createdAt
// 				labels(first: 10) {
// 				  nodes {
// 					color
// 					createdAt
// 					description
// 					id
// 					isDefault
// 					name
// 					resourcePath
// 					updatedAt
// 					url
// 				  }
// 				}
// 				locked
// 				number
// 				lastEditedAt
// 				milestone {
// 				  number
// 				  title
// 				  state
// 				  id
// 				  url
// 				}
// 				publishedAt
// 				stateReason
// 				timelineItems(last: 10) {
// 				  nodes {
// 					... on IssueComment {
// 					  id
// 					  updatedAt
// 					  createdAt
// 					}
// 				  }
// 				}
// 				participants(first: 50) {
// 				  nodes {
// 					bio
// 					bioHTML
// 					company
// 					companyHTML
// 					createdAt
// 					email
// 					id
// 					isDeveloperProgramMember
// 					location
// 					login
// 					name
// 					url
// 				  }
// 				  pageInfo {
// 					hasNextPage
// 				  }
// 				}
// 			  }
// 			}
// 			fieldValueByName(name: "Status") {
// 			  ... on ProjectV2ItemFieldSingleSelectValue {
// 				id
// 				name
// 			  }
// 			}
// 		  }
// 		}
// 	  }
// 	}
//   }
