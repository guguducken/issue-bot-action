package github

import "time"

type Issue struct {
	ID                int         `json:"id"`
	NodeID            string      `json:"node_id"`
	URL               string      `json:"url"`
	RepositoryURL     string      `json:"repository_url"`
	LabelsURL         string      `json:"labels_url"`
	CommentsURL       string      `json:"comments_url"`
	EventsURL         string      `json:"events_url"`
	HTMLURL           string      `json:"html_url"`
	Number            int         `json:"number"`
	State             string      `json:"state"`
	Title             string      `json:"title"`
	Body              string      `json:"body"`
	User              User        `json:"user"`
	Labels            []Labels    `json:"labels"`
	Assignee          User        `json:"assignee"`
	Assignees         []User      `json:"assignees"`
	Milestone         Milestone   `json:"milestone"`
	Locked            bool        `json:"locked"`
	ActiveLockReason  string      `json:"active_lock_reason"`
	Comments          int         `json:"comments"`
	PullRequest       PullRequest `json:"pull_request"`
	ClosedAt          time.Time   `json:"closed_at"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	ClosedBy          User        `json:"closed_by"`
	AuthorAssociation string      `json:"author_association"`
	StateReason       string      `json:"state_reason"`
}
type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
type Labels struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type Milestone struct {
	URL          string    `json:"url"`
	HTMLURL      string    `json:"html_url"`
	LabelsURL    string    `json:"labels_url"`
	ID           int       `json:"id"`
	NodeID       string    `json:"node_id"`
	Number       int       `json:"number"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Creator      User      `json:"creator"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	DueOn        time.Time `json:"due_on"`
}
type PullRequest struct {
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
	DiffURL  string `json:"diff_url"`
	PatchURL string `json:"patch_url"`
}

type Comment struct {
	URL                   string      `json:"url"`
	HTMLURL               string      `json:"html_url"`
	IssueURL              string      `json:"issue_url"`
	ID                    int         `json:"id"`
	NodeID                string      `json:"node_id"`
	User                  User        `json:"user"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	AuthorAssociation     string      `json:"author_association"`
	Body                  string      `json:"body"`
	Reactions             Reactions   `json:"reactions"`
	PerformedViaGithubApp interface{} `json:"performed_via_github_app"`
}

type Reactions struct {
	URL        string `json:"url"`
	TotalCount int    `json:"total_count"`
	Addone     int    `json:"+1"`
	Subone     int    `json:"-1"`
	Laugh      int    `json:"laugh"`
	Hooray     int    `json:"hooray"`
	Confused   int    `json:"confused"`
	Heart      int    `json:"heart"`
	Rocket     int    `json:"rocket"`
	Eyes       int    `json:"eyes"`
}
