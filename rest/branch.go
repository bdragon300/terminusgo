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

func (br *BranchRequester) ListAll(ctx context.Context, buf *[]Branch) (response TerminusResponse, err error) {
	di := DocumentIntroducer[Branch]{client: br.Client}
	path := br.path.(RepoPath)
	return di.OnBranch(BranchPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         path.Repo,
		Branch:       BranchCommits,
	}).ListAll(ctx, buf, &DocumentListOptions{Type: "Branch", GraphType: GraphTypeInstance, Prefixed: true})
}

type BranchCreateOptions struct {
	// Origin is the thing we wish to create a branch out of. it can be any kind of branch descriptor or commit descriptor.
	Origin   string  `json:"origin,omitempty"`
	Schema   bool    `json:"schema,omitempty"`
	Prefixes *Prefix `json:"prefixes,omitempty"`
}

func (br *BranchRequester) Create(ctx context.Context, branchID string, options *BranchCreateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := br.Client.C.BodyJSON(options).Post(br.getURL(branchID, "branch"))
	return doRequest(ctx, sl, nil)
}

func (br *BranchRequester) Delete(ctx context.Context, branchID string) (response TerminusResponse, err error) {
	sl := br.Client.C.Delete(br.getURL(branchID, "branch"))
	return doRequest(ctx, sl, nil)
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
func (br *BranchRequester) Push(ctx context.Context, branchID, remote, remoteBranch string, options *BranchPushOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		BranchPushOptions
		Remote       string `json:"remote"`
		RemoteBranch string `json:"remote_branch"`
	}{*options, remote, remoteBranch}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "push"))
	return doRequest(ctx, sl, nil)
}

type BranchPullOptions struct {
	Author  string `json:"author" default:"defaultAuthor"` // FIXME: figure out if author, message are actually used in db (and required or not)
	Message string `json:"message" default:"Default commit message"`
}

func (br *BranchRequester) Pull(ctx context.Context, branchID, remote, remoteBranch string, options *BranchPullOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		BranchPullOptions
		Remote       string `json:"remote"`
		RemoteBranch string `json:"remote_branch"`
	}{*options, remote, remoteBranch}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "pull"))
	return doRequest(ctx, sl, nil) // TODO: There is ok response also
}

type BranchSquashOptions struct {
	Author  string `json:"author" default:"defaultAuthor"` // FIXME: figure out if this field is required and default author is ok
	Message string `json:"message" default:"Default commit message"`
}

func (br *BranchRequester) Squash(ctx context.Context, branchID string, options *BranchSquashOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		CommitInfo any `json:"commit_info"`
	}{*options}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "squash"))
	return doRequest(ctx, sl, nil)
}

type BranchResetOptions struct {
	UsePath bool
}

func (br *BranchRequester) Reset(ctx context.Context, branchID, commit string, options *BranchResetOptions) (response TerminusResponse, err error) {
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
	return doRequest(ctx, sl, nil)
}

type BranchApplyOptions struct {
	Message         string          `json:"-" default:"Default commit message"`
	Author          string          `json:"-" default:"defaultAuthor"`
	Keep            map[string]bool `json:"keep" default:"{\"@id\": true, \"@type\": true}"` // Fields to keep after apply
	MatchFinalState bool            `json:"match_final_state" default:"true"`
	Type            string          `json:"type,omitempty" default:"squash"`
}

func (br *BranchRequester) Apply(ctx context.Context, branchID, beforeCommit, afterCommit string, options *BranchApplyOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	type commitInfo struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}
	body := struct {
		BranchApplyOptions
		CommitInfo   commitInfo `json:"commit_info"`
		BeforeCommit string     `json:"before_commit"`
		AfterCommit  string     `json:"after_commit"`
	}{*options, commitInfo{options.Author, options.Message}, beforeCommit, afterCommit}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "apply"))
	return doRequest(ctx, sl, nil)
}

type BranchRebaseOptions struct {
	Author string `json:"author" default:"Default author"`
}

func (br *BranchRequester) RebaseFromPath(ctx context.Context, branchID, rebaseFrom string, options *BranchRebaseOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		BranchRebaseOptions
		RebaseFrom string `json:"rebase_from"`
	}{*options, rebaseFrom}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "rebase"))
	return doRequest(ctx, sl, nil)
}

func (br *BranchRequester) Rebase(ctx context.Context, branchID string, rebaseFrom ObjectPathProvider, options *BranchRebaseOptions) (response TerminusResponse, err error) {
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
	return doRequest(ctx, sl, nil)
}

type BranchCommitLogOptions struct {
	Count int `url:"count" default:"-1"`
	Start int `url:"start" default:"0"`
}

func (br *BranchRequester) CommitLog(ctx context.Context, buf *[]Commit, branchID string, options *BranchCommitLogOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := br.Client.C.QueryStruct(options).Get(br.getURL(branchID, "log"))
	return doRequest(ctx, sl, buf)
}

func (br *BranchRequester) Optimize(ctx context.Context, branchID string) (response TerminusResponse, err error) {
	sl := br.Client.C.Post(br.getURL(branchID, "optimize"))
	return doRequest(ctx, sl, nil)
}

type BranchSchemaFrameOptions struct {
	CompressIDs    bool `json:"compress_ids" default:"true"`
	ExpandAbstract bool `json:"expand_abstract" default:"true"`
}

func (br *BranchRequester) SchemaFrameAll(ctx context.Context, buf *[]schema.RawSchemaItem, name string, options *BranchSchemaFrameOptions) (response TerminusResponse, err error) {
	var resp map[string]map[string]any
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := br.Client.C.QueryStruct(options).Get(br.getURL(name, "schema"))
	response, err = doRequest(ctx, sl, &resp)
	if err != nil {
		return
	}

	for k, v := range resp {
		v["@id"] = k
		*buf = append(*buf, v)
	}
	return
}

func (br *BranchRequester) SchemaFrameType(ctx context.Context, buf *schema.RawSchemaItem, name, docType string, options *BranchSchemaFrameOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	params := struct {
		BranchSchemaFrameOptions
		Type string `json:"type"`
	}{*options, docType}
	sl := br.Client.C.QueryStruct(params).Get(br.getURL(name, "schema"))
	return doRequest(ctx, sl, buf)
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
