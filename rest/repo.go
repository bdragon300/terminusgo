package rest

import (
	"context"
	"fmt"
	"net/url"
)

type Repo struct {
	// TODO: local and remote repos
	Name string
}

type RepoIntroducer BaseIntroducer

func (ri *RepoIntroducer) OnDatabase(path DatabasePath) *RepoRequester {
	return &RepoRequester{Client: ri.client, path: path}
}

type RepoRequester BaseRequester

type RepoFetchOptions struct {
	RemoteURL           string `validate:"required,url" default:"http://example.com/user/test_db"` // FIXME: in python client it is called "remote id", check why
	RemoteAuthorization string `validate:"required" default:"TOKEN"`                               // FIXME: in python client it does not exist
}

// TODO: this relates either to repo or to branch, figure out
func (rr *RepoRequester) Fetch(ctx context.Context, repoID string, options *RepoFetchOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	URL := rr.getURL(repoID, "fetch")
	sl := rr.Client.C.Set("AUTHORIZATION_REMOTE", options.RemoteAuthorization).Post(URL)
	if _, err = doRequest(ctx, sl, nil); err != nil { // FIXME: there is ok response
		return err
	}

	return nil
}

func (rr *RepoRequester) Optimize(ctx context.Context, repoID string) error {
	sl := rr.Client.C.Post(rr.getURL(repoID, "optimize"))
	if _, err := doRequest(ctx, sl, nil); err != nil { // TODO: There is ok response also
		return err
	}

	return nil
}

func (rr *RepoRequester) getURL(objectID, action string) string {
	path := rr.path.(DatabasePath)
	return RepoPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         objectID,
	}.GetPath(action)
}

type RepoPath struct {
	Organization, Database, Repo string
}

func (rp RepoPath) GetPath(action string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		action,
		getDBBase(rp.Database, rp.Organization),
		url.QueryEscape(rp.Repo),
	)
}
