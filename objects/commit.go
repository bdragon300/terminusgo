package objects

import (
	"fmt"
	"net/url"
)

type Commit struct {
	ID         string  `json:"@id"`
	Type       string  `json:"@type"`
	Author     string  `json:"author"`
	Identifier string  `json:"identifier"`
	Instance   string  `json:"instance"`
	Message    string  `json:"message"`
	Parent     string  `json:"parent"`
	Schema     string  `json:"schema"`
	Timestamp  float64 `json:"timestamp"`
	// FIXME: commit can be identified by id or a path (db, branch, etc.)
}

type CommitIntroducer BaseIntroducer

func (ci *CommitIntroducer) OnBranch(path BranchPath) *CommitRequester {
	return &CommitRequester{Client: ci.client, path: path}
}

type CommitRequester BaseRequester

// TODO: implement smth or remove

type CommitPath struct {
	Organization, Database, Repo, Branch, Commit string
}

func (cp CommitPath) GetPath(action string) string {
	return fmt.Sprintf(
		"%s/%s/%s/commit/%s",
		action,
		getDBBase(cp.Database, cp.Organization),
		url.QueryEscape(cp.Repo),
		cp.Commit,
	)
}
