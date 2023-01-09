package srverror

import (
	"fmt"
	"net/http"
)

const (
	DataVersionHeader         = "Terminusdb-Data-Version"
	RemoteAuthorizationHeader = "Authorization-Remote"
	XTerminusDBApiBaseHeader  = "X-Terminusdb-Api-Base"
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

func (ter TerminusErrorResponse) String() string {
	return fmt.Sprintf("terminus db server returned HTTP code %d: %s", ter.Response.StatusCode, ter.APIMessage)
}

func (ter TerminusErrorResponse) IsOK() bool {
	return false
}

func (ter TerminusErrorResponse) DataVersion() string {
	return ter.Response.Header.Get(DataVersionHeader)
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

func (tor TerminusOkResponse) String() string {
	return "successful request"
}

func (tor TerminusOkResponse) IsOK() bool {
	return true
}

func (tor TerminusOkResponse) DataVersion() string {
	return tor.Response.Header.Get(DataVersionHeader)
}

type WOQLResponse struct {
	APIStatus             string         `json:"api:status"`
	APIVariableNames      string         `json:"api:variable_names"`
	Bindings              map[string]any `json:"bindings"` // TODO: implement smth to extract/working with this field (and the one above)
	Inserts               uint           `json:"inserts"`
	Deletes               uint           `json:"deletes"`
	TransactionRetryCount uint           `json:"transaction_retry_count"`
}
