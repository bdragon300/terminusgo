package srverror

import (
	"fmt"
	"net/http"
)

const (
	DataVersionHeader         = "Terminusdb-Data-Version"
	RemoteAuthorizationHeader = "Authorization-Remote"
)

type TerminusErrorResponse struct {
	Type            string           `json:"@type"`
	Auth            string           `json:"auth"`
	Scope           string           `json:"scope"`
	APIError        string           `json:"api:error"` // FIXME: string or dict
	APIPath         string           `json:"api:path"`
	APIMessage      string           `json:"api:message"`
	APIStatus       string           `json:"api:status"`
	SystemWitnesses []map[string]any `json:"system:witnesses"`

	Response *http.Response
}

func (te TerminusErrorResponse) String() string {
	return fmt.Sprintf("terminus db server returned HTTP code %d: %s", te.Response.StatusCode, te.APIMessage)
}

func (te TerminusErrorResponse) IsOK() bool {
	return false
}

func (te TerminusErrorResponse) DataVersion() string {
	return te.Response.Header.Get(DataVersionHeader)
}

type TerminusOkResponse struct {
	Type      string `json:"@type"`
	APIStatus string `json:"api:success"`
	// Push
	APIRepoHeadUpdated bool   `json:"api:repo_head_updated"`
	APIRepoHead        string `json:"api:repo_head"`
	// Squash
	APICommit      string `json:"api:commit"`
	APIOldCommit   string `json:"api:old_commit"`
	APIEmptyCommit bool   `json:"api:empty_commit"`
	// Rebase
	APIForwardedCommits string `json:"api:forwarded_commits"`
	APIRebaseReport     string `json:"api:rebase_report"`
	APICommonCommitID   string `json:"api:common_commit_id"`
	// DB
	APIDatabaseURI string `json:"api:database_uri"`

	Response *http.Response
}

func (or TerminusOkResponse) IsOK() bool {
	return true
}

func (or TerminusOkResponse) String() string {
	return "Successful request"
}
