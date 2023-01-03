package rest

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bdragon300/terminusgo/srverror"
)

type RepoIntroducer BaseIntroducer

func (ri *RepoIntroducer) OnDatabase(path DatabasePath) *RepoRequester {
	return &RepoRequester{BaseRequester: BaseRequester{Client: ri.client, path: path}}
}

type RepoRequester struct {
	BaseRequester
	remoteAuthorization string
}

// Contents will be put to 'Authorization' header to a remote POST request, see src/core/api/db_fetch.pl:authorized_fetch() in TerminusDB sources
func (rr *RepoRequester) WithRemoteAuth(contents string) *RepoRequester {
	rr.remoteAuthorization = contents
	return rr
}

func (rr *RepoRequester) Fetch(ctx context.Context, repoID string) (response TerminusResponse, err error) {
	// Implementation in db: src/core/api/db_fetch.pl:remote_fetch(). Quite awkward IMHO
	path := rr.path.(DatabasePath)
	URL := BranchPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         repoID,
		Branch:       BranchCommits,
	}.GetURL("fetch")
	sl := rr.Client.C.Post(URL)
	if rr.remoteAuthorization != "" {
		sl = sl.Set(srverror.RemoteAuthorizationHeader, rr.remoteAuthorization)
	}
	return doRequest(ctx, sl, nil)
}

func (rr *RepoRequester) Optimize(ctx context.Context, repoID string) (response TerminusResponse, err error) {
	path := rr.path.(DatabasePath)
	URL := RepoPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         repoID,
	}.GetURL("optimize")
	sl := rr.Client.C.Post(URL)
	return doRequest(ctx, sl, nil)
}

type RepoPath struct {
	Organization, Database, Repo string
}

func (rp RepoPath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, rp.GetPath())
}

func (rp RepoPath) GetPath() string {
	return fmt.Sprintf("%s/%s", getDBBase(rp.Database, rp.Organization), url.QueryEscape(rp.Repo))
}
