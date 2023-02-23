package rest

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bdragon300/terminusgo/schema"
)

type Branch struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
	Name string `json:"name"`
	Head string `json:"head"`
}

func (b Branch) GetPath() RepoPath {
	return RepoPath{} // TODO
}

type BranchIntroducer BaseIntroducer

func (bi *BranchIntroducer) OnRepo(path RepoPath) *BranchRequester {
	return &BranchRequester{Client: bi.client, path: path}
}

type BranchRequester BaseRequester

func (br *BranchRequester) WithContext(ctx context.Context) *BranchRequester {
	r := *br
	r.ctx = ctx
	return &r
}

func (br *BranchRequester) ListAll(buf *[]Branch) (response TerminusResponse, err error) {
	di := DocumentIntroducer[Branch]{client: br.Client}
	path := br.path.(RepoPath)
	return di.OnBranch(BranchPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         path.Repo,
		Branch:       BranchCommits,
	}).WithContext(br.ctx).ListAll(buf, &DocumentListOptions{Type: "Branch", GraphType: GraphTypeInstance, Prefixed: true})
}

type BranchCreateOptions struct {
	// Origin is the thing we wish to create a branch out of. it can be any kind of branch descriptor or commit descriptor.
	Origin   string  `json:"origin,omitempty"`
	Schema   bool    `json:"schema,omitempty"`
	Prefixes *Prefix `json:"prefixes,omitempty"`
}

func (br *BranchRequester) Create(branchID string, options *BranchCreateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := br.Client.C.BodyJSON(options).Post(br.getURL(branchID, "branch"))
	return doRequest(br.ctx, sl, nil)
}

func (br *BranchRequester) Delete(branchID string) (response TerminusResponse, err error) {
	sl := br.Client.C.Delete(br.getURL(branchID, "branch"))
	return doRequest(br.ctx, sl, nil)
}

type BranchPushOptions struct {
	PushPrefixes bool   `json:"push_prefixes" default:"true"`
	Author       string `json:"author" default:"defaultAuthor"` // FIXME: figure out if author, message are actually used in db (and required or not)
	Message      string `json:"message" default:"Default commit message"`
}

// error conditions:
// - branch to push does not exist
// - repository does not exist
// - we tried to push to a repository that is not a remote
// - tried to push without having fetched first. The repository exists as an entity in our metadata graph but it hasn't got an associated commit graph. We always need one.
// - remote diverged - someone else committed and pushed and we know about that
// - We try to push an empty branch, but we know that remote is non-empty
// - remote returns an error
// -- history diverged (we check locally, but there's a race)
// -- remote doesn't know what we're talking about
// -- remote authorization failed
// - communication error while talking to the remote
func (br *BranchRequester) Push(branchID, remote, remoteBranch string, options *BranchPushOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		BranchPushOptions
		Remote       string `json:"remote"`
		RemoteBranch string `json:"remote_branch"`
	}{*options, remote, remoteBranch}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "push"))
	return doRequest(br.ctx, sl, nil)
}

type BranchPullOptions struct {
	Author  string `json:"author" default:"defaultAuthor"` // FIXME: figure out if author, message are actually used in db (and required or not)
	Message string `json:"message" default:"Default commit message"`
}

func (br *BranchRequester) Pull(branchID, remote, remoteBranch string, options *BranchPullOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		BranchPullOptions
		Remote       string `json:"remote"`
		RemoteBranch string `json:"remote_branch"`
	}{*options, remote, remoteBranch}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "pull"))
	return doRequest(br.ctx, sl, nil) // TODO: There is ok response also
}

type BranchSquashOptions struct {
	Author  string `json:"author" default:"defaultAuthor"` // FIXME: figure out if this field is required and default author is ok
	Message string `json:"message" default:"Default commit message"`
}

func (br *BranchRequester) Squash(branchID string, options *BranchSquashOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	commitInfo := struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}{Author: options.Author, Message: options.Message}
	body := struct {
		CommitInfo any `json:"commit_info"`
	}{commitInfo}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "squash"))
	return doRequest(br.ctx, sl, nil)
}

type BranchResetOptions struct {
	UsePath bool
}

func (br *BranchRequester) Reset(branchID, commit string, options *BranchResetOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	if !options.UsePath {
		path := br.path.(RepoPath)
		commit = CommitPath{
			Organization: path.Organization,
			Database:     path.Database,
			Repo:         path.Repo,
			Branch:       branchID,
			Commit:       commit,
		}.GetPath()
	}
	body := struct {
		Commit string `json:"commit_descriptor"`
	}{commit}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "reset"))
	return doRequest(br.ctx, sl, nil)
}

type BranchApplyOptions struct {
	Message         string          `json:"-" default:"Default commit message"`
	Author          string          `json:"-" default:"defaultAuthor"`
	Keep            map[string]bool `json:"keep" default:"{\"@id\": true, \"@type\": true}"` // Fields to keep after apply
	MatchFinalState bool            `json:"match_final_state" default:"true"`
	Type            string          `json:"type,omitempty" default:"squash"`
}

func (br *BranchRequester) Apply(branchID, beforeCommit, afterCommit string, options *BranchApplyOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	commitInfo := struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}{Author: options.Author, Message: options.Message}
	body := struct {
		BranchApplyOptions
		CommitInfo   any    `json:"commit_info"`
		BeforeCommit string `json:"before_commit"`
		AfterCommit  string `json:"after_commit"`
	}{*options, commitInfo, beforeCommit, afterCommit}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "apply"))
	return doRequest(br.ctx, sl, nil)
}

type BranchRebaseOptions struct {
	Author string `json:"author" default:"Default author"`
}

func (br *BranchRequester) RebaseFromPath(branchID, rebaseFrom string, options *BranchRebaseOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		BranchRebaseOptions
		RebaseFrom string `json:"rebase_from"`
	}{*options, rebaseFrom}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "rebase"))
	return doRequest(br.ctx, sl, nil)
}

func (br *BranchRequester) Rebase(branchID string, rebaseFrom ObjectPathProvider, options *BranchRebaseOptions) (response TerminusResponse, err error) {
	if rebaseFrom == nil {
		panic("rebaseFrom is nil")
	}
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		BranchRebaseOptions
		RebaseFrom string `json:"rebase_from"`
	}{*options, rebaseFrom.GetPath()}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "rebase"))
	return doRequest(br.ctx, sl, nil)
}

type BranchCommitLogOptions struct {
	Count int `url:"count" default:"-1"`
	Start int `url:"start" default:"0"`
}

func (br *BranchRequester) CommitLog(branchID string, buf *[]Commit, options *BranchCommitLogOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := br.Client.C.QueryStruct(options).Get(br.getURL(branchID, "log"))
	return doRequest(br.ctx, sl, buf)
}

func (br *BranchRequester) Optimize(branchID string) (response TerminusResponse, err error) {
	sl := br.Client.C.Post(br.getURL(branchID, "optimize"))
	return doRequest(br.ctx, sl, nil)
}

type BranchSchemaFrameOptions struct {
	CompressIDs    bool `url:"compress_ids" default:"true"`
	ExpandAbstract bool `url:"expand_abstract" default:"true"`
}

func (br *BranchRequester) SchemaFrameAll(name string, buf *[]schema.RawSchemaItem, options *BranchSchemaFrameOptions) (response TerminusResponse, err error) {
	var resp map[string]map[string]any
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := br.Client.C.QueryStruct(options).Get(br.getURL(name, "schema"))
	response, err = doRequest(br.ctx, sl, &resp)
	if err != nil {
		return
	}

	for k, v := range resp {
		v["@id"] = k
		*buf = append(*buf, v)
	}
	return
}

func (br *BranchRequester) SchemaFrameType(name, docType string, buf *schema.RawSchemaItem, options *BranchSchemaFrameOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	params := struct {
		BranchSchemaFrameOptions
		Type string `url:"type"`
	}{*options, docType}
	sl := br.Client.C.QueryStruct(params).Get(br.getURL(name, "schema"))
	return doRequest(br.ctx, sl, buf)
}

func (br *BranchRequester) getURL(branchID, action string) string {
	path := br.path.(RepoPath)
	return BranchPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         path.Repo,
		Branch:       branchID,
	}.GetURL(action)
}

type BranchPath struct {
	Organization, Database, Repo, Branch string
}

func (bp BranchPath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, bp.GetPath())
}

func (bp BranchPath) GetPath() string {
	suburl := fmt.Sprintf(
		"%s/%s",
		getDBBase(bp.Database, bp.Organization),
		url.QueryEscape(bp.Repo),
	)
	if bp.Repo == RepoMeta {
		return suburl
	}
	if bp.Branch == BranchCommits {
		return fmt.Sprintf("%s/%s", suburl, bp.Branch)
	}
	return fmt.Sprintf("%s/branch/%s", suburl, bp.Branch)
}
