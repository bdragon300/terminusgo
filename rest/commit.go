package rest

import (
	"fmt"
	"net/url"
	"strings"
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

type CommitPath struct {
	Organization, Database, Repo, Branch, Commit string
}

func (cp CommitPath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, cp.String())
}

func (cp CommitPath) String() string {
	return fmt.Sprintf(
		"%s/%s/commit/%s",
		getDatabasePath(cp.Organization, cp.Database),
		url.PathEscape(cp.Repo),
		url.PathEscape(cp.Commit),
	)
	// FIXME: cp.Branch is not used (+ see FromString if fix is needed)
}

func (cp CommitPath) FromString(s string) CommitPath {
	res := CommitPath{}
	parts := strings.SplitN(s, "/", 5)
	if parts[0] == DatabaseSystem {
		parts = append(parts[:1], parts[0:]...) // Insert empty Organization part
		parts[0] = ""
	}
	if len(parts) < 5 {
		panic(fmt.Sprintf("too short path %q", s))
	}
	parts = append(parts[:3], parts[4:]...) // Cut "commit" part
	fillUnescapedStringFields(parts, &res)
	return res
}
