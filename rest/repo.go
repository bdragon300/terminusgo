package rest

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bdragon300/terminusgo/schema"

	"github.com/bdragon300/terminusgo/srverror"
)

type RepoIntroducer BaseIntroducer

func (ri *RepoIntroducer) OnDatabase(path DatabasePath) *RepoRequester {
	return &RepoRequester{BaseRequester: BaseRequester{Client: ri.client, path: path}}
}

type RepoRequester struct {
	BaseRequester
}

func (rr *RepoRequester) WithContext(ctx context.Context) *RepoRequester {
	r := *rr
	r.ctx = ctx
	return &r
}

func (rr *RepoRequester) Fetch(repoID, remoteAuthorization string) (response TerminusResponse, err error) {
	// Implementation in db: src/core/api/db_fetch.pl:remote_fetch(). Quite awkward IMHO
	path := rr.path.(DatabasePath)
	URL := BranchPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         repoID,
		Branch:       BranchCommits,
	}.GetURL("fetch")
	sl := rr.Client.C.Post(URL).Set(srverror.RemoteAuthorizationHeader, remoteAuthorization)
	return doRequest(rr.ctx, sl, nil)
}

func (rr *RepoRequester) Optimize(repoID string) (response TerminusResponse, err error) {
	sl := rr.Client.C.Post(rr.getURL(repoID, "optimize"))
	return doRequest(rr.ctx, sl, nil)
}

type RepoSchemaFrameOptions struct {
	CompressIDs    bool `url:"compress_ids" default:"true"`
	ExpandAbstract bool `url:"expand_abstract" default:"true"`
}

func (rr *RepoRequester) SchemaFrameAll(name string, buf *[]schema.RawSchemaItem, options *RepoSchemaFrameOptions) (response TerminusResponse, err error) {
	var resp map[string]map[string]any
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := rr.Client.C.QueryStruct(options).Get(rr.getURL(name, "schema"))
	response, err = doRequest(rr.ctx, sl, &resp)
	if err != nil {
		return
	}

	for k, v := range resp {
		v["@id"] = k
		*buf = append(*buf, v)
	}
	return
}

func (rr *RepoRequester) SchemaFrameType(name, docType string, buf *schema.RawSchemaItem, options *RepoSchemaFrameOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	params := struct {
		RepoSchemaFrameOptions
		Type string `url:"type"`
	}{*options, docType}
	sl := rr.Client.C.QueryStruct(params).Get(rr.getURL(name, "schema"))
	return doRequest(rr.ctx, sl, buf)
}

func (rr *RepoRequester) getURL(repoID, action string) string {
	path := rr.path.(DatabasePath)
	return RepoPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         repoID,
	}.GetURL(action)
}

type RepoPath struct {
	Organization, Database, Repo string
}

func (rp RepoPath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, rp.String())
}

func (rp RepoPath) String() string {
	return fmt.Sprintf("%s/%s", getDatabasePath(rp.Organization, rp.Database), url.PathEscape(rp.Repo))
}

func (rp RepoPath) FromString(s string) RepoPath {
	res := RepoPath{}
	parts := strings.SplitN(s, "/", 3)
	if parts[0] == DatabaseSystem {
		parts = append(parts[:1], parts[0:]...) // Insert empty Organization part
		parts[0] = ""
	}
	if len(parts) < 3 {
		panic(fmt.Sprintf("too short path %q", s))
	}
	fillUnescapedStringFields(parts, &res)
	return res
}
