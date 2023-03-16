package github

import "time"

type Issue struct {
	ID                int          `json:"id"`
	NodeID            string       `json:"node_id"`
	URL               string       `json:"url"`
	RepositoryURL     string       `json:"repository_url"`
	LabelsURL         string       `json:"labels_url"`
	CommentsURL       string       `json:"comments_url"`
	EventsURL         string       `json:"events_url"`
	HTMLURL           string       `json:"html_url"`
	Number            int          `json:"number"`
	State             string       `json:"state"`
	Title             string       `json:"title"`
	Body              string       `json:"body"`
	User              *User        `json:"user"`
	Labels            []Label      `json:"labels"`
	Assignee          *User        `json:"assignee"`
	Assignees         []User       `json:"assignees"`
	Milestone         *Milestone   `json:"milestone"`
	Locked            bool         `json:"locked"`
	ActiveLockReason  string       `json:"active_lock_reason"`
	Comments          int          `json:"comments"`
	PullRequest       *PullRequest `json:"pull_request"`
	ClosedAt          *time.Time   `json:"closed_at"`
	CreatedAt         *time.Time   `json:"created_at"`
	UpdatedAt         *time.Time   `json:"updated_at"`
	ClosedBy          *User        `json:"closed_by"`
	AuthorAssociation string       `json:"author_association"`
	StateReason       string       `json:"state_reason"`
	StartTime         *TimeProject
	EndTime           *TimeProject
	Status            string
}

type TimeProject struct {
	ID        string
	Date      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Event struct {
	ID        int64     `json:"id"`
	NodeID    string    `json:"node_id"`
	URL       string    `json:"url"`
	Actor     User      `json:"actor"`
	Event     string    `json:"event"`
	CommitID  string    `json:"commit_id"`
	CommitURL string    `json:"commit_url"`
	CreatedAt time.Time `json:"created_at"`
	Assignee  *User     `json:"assignee"`
	Assigner  *User     `json:"assigner"`
}

type User struct {
	Login             string     `json:"login"`
	ID                int        `json:"id"`
	NodeID            string     `json:"node_id"`
	AvatarURL         string     `json:"avatar_url"`
	GravatarID        string     `json:"gravatar_id"`
	URL               string     `json:"url"`
	HTMLURL           string     `json:"html_url"`
	FollowersURL      string     `json:"followers_url"`
	FollowingURL      string     `json:"following_url"`
	GistsURL          string     `json:"gists_url"`
	StarredURL        string     `json:"starred_url"`
	SubscriptionsURL  string     `json:"subscriptions_url"`
	OrganizationsURL  string     `json:"organizations_url"`
	ReposURL          string     `json:"repos_url"`
	EventsURL         string     `json:"events_url"`
	ReceivedEventsURL string     `json:"received_events_url"`
	Type              string     `json:"type"`
	SiteAdmin         bool       `json:"site_admin"`
	Name              string     `json:"name"`
	Company           string     `json:"company"`
	Blog              string     `json:"blog"`
	Location          string     `json:"location"`
	Email             string     `json:"email"`
	Hireable          bool       `json:"hireable"`
	Bio               string     `json:"bio"`
	TwitterUsername   string     `json:"twitter_username"`
	PublicRepos       int        `json:"public_repos"`
	PublicGists       int        `json:"public_gists"`
	Followers         int        `json:"followers"`
	Following         int        `json:"following"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}
type Label struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type Milestone struct {
	URL          string     `json:"url"`
	HTMLURL      string     `json:"html_url"`
	LabelsURL    string     `json:"labels_url"`
	ID           int        `json:"id"`
	NodeID       string     `json:"node_id"`
	Number       int        `json:"number"`
	State        string     `json:"state"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Creator      *User      `json:"creator"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	DueOn        *time.Time `json:"due_on"`
}

type Comment struct {
	URL                   string      `json:"url"`
	HTMLURL               string      `json:"html_url"`
	IssueURL              string      `json:"issue_url"`
	ID                    int         `json:"id"`
	NodeID                string      `json:"node_id"`
	User                  User        `json:"user"`
	CreatedAt             *time.Time  `json:"created_at"`
	UpdatedAt             *time.Time  `json:"updated_at"`
	AuthorAssociation     string      `json:"author_association"`
	Body                  string      `json:"body"`
	Reactions             *Reactions  `json:"reactions"`
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

type PullRequest struct {
	URL                 string     `json:"url"`
	ID                  int        `json:"id"`
	NodeID              string     `json:"node_id"`
	HTMLURL             string     `json:"html_url"`
	DiffURL             string     `json:"diff_url"`
	PatchURL            string     `json:"patch_url"`
	IssueURL            string     `json:"issue_url"`
	CommitsURL          string     `json:"commits_url"`
	ReviewCommentsURL   string     `json:"review_comments_url"`
	ReviewCommentURL    string     `json:"review_comment_url"`
	CommentsURL         string     `json:"comments_url"`
	StatusesURL         string     `json:"statuses_url"`
	Number              int        `json:"number"`
	State               string     `json:"state"`
	Locked              bool       `json:"locked"`
	Title               string     `json:"title"`
	User                *User      `json:"user"`
	Body                string     `json:"body"`
	Labels              []Label    `json:"labels"`
	Milestone           Milestone  `json:"milestone"`
	ActiveLockReason    string     `json:"active_lock_reason"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	ClosedAt            *time.Time `json:"closed_at"`
	MergedAt            *time.Time `json:"merged_at"`
	MergeCommitSha      string     `json:"merge_commit_sha"`
	Assignee            *User      `json:"assignee"`
	Assignees           []User     `json:"assignees"`
	RequestedReviewers  []User     `json:"requested_reviewers"`
	RequestedTeams      []Team     `json:"requested_teams"`
	Head                *Head      `json:"head"`
	Base                *Base      `json:"base"`
	Links               *Links     `json:"_links"`
	AuthorAssociation   string     `json:"author_association"`
	AutoMerge           *AutoMerge `json:"auto_merge"`
	Draft               bool       `json:"draft"`
	Merged              bool       `json:"merged"`
	Mergeable           bool       `json:"mergeable"`
	Rebaseable          bool       `json:"rebaseable"`
	MergeableState      string     `json:"mergeable_state"`
	MergedBy            *User      `json:"merged_by"`
	Comments            int        `json:"comments"`
	ReviewComments      int        `json:"review_comments"`
	MaintainerCanModify bool       `json:"maintainer_can_modify"`
	Commits             int        `json:"commits"`
	Additions           int        `json:"additions"`
	Deletions           int        `json:"deletions"`
	ChangedFiles        int        `json:"changed_files"`
}

type AutoMerge struct {
	EnableBy      *User  `json:"enabled_by"`
	MergeMethod   string `json:"merge_method"`
	CommitTitle   string `json:"commit_title"`
	CommitMessage string `json:"commit_message"`
}

type Team struct {
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	URL             string `json:"url"`
	HTMLURL         string `json:"html_url"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Description     string `json:"description"`
	Privacy         string `json:"privacy"`
	Permission      string `json:"permission"`
	MembersURL      string `json:"members_url"`
	RepositoriesURL string `json:"repositories_url"`
}
type Permissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}
type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	SpdxID string `json:"spdx_id"`
	NodeID string `json:"node_id"`
}
type Head struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  *User  `json:"user"`
	Repo  *Repo  `json:"repo"`
}
type Repo struct {
	ID               int          `json:"id"`
	NodeID           string       `json:"node_id"`
	Name             string       `json:"name"`
	FullName         string       `json:"full_name"`
	Owner            *User        `json:"owner"`
	Private          bool         `json:"private"`
	HTMLURL          string       `json:"html_url"`
	Description      string       `json:"description"`
	Fork             bool         `json:"fork"`
	URL              string       `json:"url"`
	ArchiveURL       string       `json:"archive_url"`
	AssigneesURL     string       `json:"assignees_url"`
	BlobsURL         string       `json:"blobs_url"`
	BranchesURL      string       `json:"branches_url"`
	CollaboratorsURL string       `json:"collaborators_url"`
	CommentsURL      string       `json:"comments_url"`
	CommitsURL       string       `json:"commits_url"`
	CompareURL       string       `json:"compare_url"`
	ContentsURL      string       `json:"contents_url"`
	ContributorsURL  string       `json:"contributors_url"`
	DeploymentsURL   string       `json:"deployments_url"`
	DownloadsURL     string       `json:"downloads_url"`
	EventsURL        string       `json:"events_url"`
	ForksURL         string       `json:"forks_url"`
	GitCommitsURL    string       `json:"git_commits_url"`
	GitRefsURL       string       `json:"git_refs_url"`
	GitTagsURL       string       `json:"git_tags_url"`
	GitURL           string       `json:"git_url"`
	IssueCommentURL  string       `json:"issue_comment_url"`
	IssueEventsURL   string       `json:"issue_events_url"`
	IssuesURL        string       `json:"issues_url"`
	KeysURL          string       `json:"keys_url"`
	LabelsURL        string       `json:"labels_url"`
	LanguagesURL     string       `json:"languages_url"`
	MergesURL        string       `json:"merges_url"`
	MilestonesURL    string       `json:"milestones_url"`
	NotificationsURL string       `json:"notifications_url"`
	PullsURL         string       `json:"pulls_url"`
	ReleasesURL      string       `json:"releases_url"`
	SSHURL           string       `json:"ssh_url"`
	StargazersURL    string       `json:"stargazers_url"`
	StatusesURL      string       `json:"statuses_url"`
	SubscribersURL   string       `json:"subscribers_url"`
	SubscriptionURL  string       `json:"subscription_url"`
	TagsURL          string       `json:"tags_url"`
	TeamsURL         string       `json:"teams_url"`
	TreesURL         string       `json:"trees_url"`
	CloneURL         string       `json:"clone_url"`
	MirrorURL        string       `json:"mirror_url"`
	HooksURL         string       `json:"hooks_url"`
	SvnURL           string       `json:"svn_url"`
	Homepage         string       `json:"homepage"`
	Language         interface{}  `json:"language"`
	ForksCount       int          `json:"forks_count"`
	StargazersCount  int          `json:"stargazers_count"`
	WatchersCount    int          `json:"watchers_count"`
	Size             int          `json:"size"`
	DefaultBranch    string       `json:"default_branch"`
	OpenIssuesCount  int          `json:"open_issues_count"`
	Topics           []string     `json:"topics"`
	HasIssues        bool         `json:"has_issues"`
	HasProjects      bool         `json:"has_projects"`
	HasWiki          bool         `json:"has_wiki"`
	HasPages         bool         `json:"has_pages"`
	HasDownloads     bool         `json:"has_downloads"`
	HasDiscussions   bool         `json:"has_discussions"`
	Archived         bool         `json:"archived"`
	Disabled         bool         `json:"disabled"`
	PushedAt         *time.Time   `json:"pushed_at"`
	CreatedAt        *time.Time   `json:"created_at"`
	UpdatedAt        *time.Time   `json:"updated_at"`
	Permissions      *Permissions `json:"permissions"`
	AllowRebaseMerge bool         `json:"allow_rebase_merge"`
	TempCloneToken   string       `json:"temp_clone_token"`
	AllowSquashMerge bool         `json:"allow_squash_merge"`
	AllowMergeCommit bool         `json:"allow_merge_commit"`
	Forks            int          `json:"forks"`
	OpenIssues       int          `json:"open_issues"`
	License          *License     `json:"license"`
	Watchers         int          `json:"watchers"`
}
type Base struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  *User  `json:"user"`
	Repo  *Repo  `json:"repo"`
}
type Link_Self struct {
	Href string `json:"href"`
}
type Link_HTML struct {
	Href string `json:"href"`
}
type Link_Issue struct {
	Href string `json:"href"`
}
type Link_Comments struct {
	Href string `json:"href"`
}
type Link_ReviewComments struct {
	Href string `json:"href"`
}
type Link_ReviewComment struct {
	Href string `json:"href"`
}
type Link_Commits struct {
	Href string `json:"href"`
}
type Link_Statuses struct {
	Href string `json:"href"`
}
type Links struct {
	Self           *Link_Self           `json:"self"`
	HTML           *Link_HTML           `json:"html"`
	Issue          *Link_Issue          `json:"issue"`
	Comments       *Link_Comments       `json:"comments"`
	ReviewComments *Link_ReviewComments `json:"review_comments"`
	ReviewComment  *Link_ReviewComment  `json:"review_comment"`
	Commits        *Link_Commits        `json:"commits"`
	Statuses       *Link_Statuses       `json:"statuses"`
}

type HeadCommit struct {
	Name   string `json:"name"`
	Commit struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	Protected bool `json:"protected"`
}

type Commit struct {
	URL         string `json:"url"`
	Sha         string `json:"sha"`
	NodeID      string `json:"node_id"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Commit      struct {
		URL    string `json:"url"`
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			URL string `json:"url"`
			Sha string `json:"sha"`
		} `json:"tree"`
		CommentCount int `json:"comment_count"`
		Verification struct {
			Verified  bool        `json:"verified"`
			Reason    string      `json:"reason"`
			Signature interface{} `json:"signature"`
			Payload   interface{} `json:"payload"`
		} `json:"verification"`
	} `json:"commit"`
	Author    User `json:"author"`
	Committer User `json:"committer"`
	Parents   []struct {
		URL string `json:"url"`
		Sha string `json:"sha"`
	} `json:"parents"`
}
