package objects

import (
	"context"
	"fmt"
	"net/url"
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

type BranchListOptions struct {
	LastDataVersion string // TODO: figure out with this
}

func (br *BranchRequester) ListAll(ctx context.Context, buf *[]Branch, _ *BranchListOptions) error {
	di := DocumentIntroducer[Branch]{client: br.Client}
	path := br.path.(RepoPath)
	err := di.OnBranch(BranchPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         path.Repo,
		Branch:       BranchCommits,
	}).ListAll(ctx, buf, &DocumentListOptions{Type: "Branch", GraphType: GraphTypeInstance, Prefixed: true})
	if err != nil {
		return err
	}

	return nil
}

type BranchCreateOptions struct {
	Origin string `json:"origin,omitempty"`
}

func (br *BranchRequester) Create(ctx context.Context, branchID string, options *BranchCreateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := br.Client.C.BodyJSON(options).Post(br.getURL(branchID, "branch"))
	_, err = doRequest(ctx, sl, nil)
	return
}

func (br *BranchRequester) Delete(ctx context.Context, branchID string) error {
	sl := br.Client.C.Delete(br.getURL(branchID, "branch"))
	_, err := doRequest(ctx, sl, nil)
	return err
}

type BranchPushOptions struct {
	PushPrefixes bool   `json:"push_prefixes"` // FIXME: this is not in python client
	Remote       string `json:"remote" validate:"required" default:"defaultRemote"`
	RemoteBranch string `json:"remote_branch" validate:"required" default:"origin"` // FIXME: is such default ok?
	Author       string `json:"author" default:"defaultAuthor"`                     // FIXME: figure out if this field is required and default author is ok
	Message      string `json:"message" default:"Default commit message"`
}

func (br *BranchRequester) Push(ctx context.Context, branchID string, options *BranchPushOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := br.Client.C.BodyJSON(options).Post(br.getURL(branchID, "push"))
	_, err = doRequest(ctx, sl, nil) // TODO: There is ok response also
	return
}

type BranchPullOptions struct {
	Remote       string `json:"remote" validate:"required" default:"defaultRemote"`
	RemoteBranch string `json:"remote_branch" validate:"required" default:"origin"` // FIXME: is such default ok?
	Author       string `json:"author" default:"defaultAuthor"`                     // FIXME: figure out if this field is required and default author is ok
	Message      string `json:"message" default:"Default commit message"`           // FIXME: author/message passibly is CommitInfo or smth like this
}

func (br *BranchRequester) Pull(ctx context.Context, branchID string, options *BranchPullOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := br.Client.C.BodyJSON(options).Post(br.getURL(branchID, "pull"))
	_, err = doRequest(ctx, sl, nil) // TODO: There is ok response also
	return
}

type BranchSquashOptions struct {
	Author  string `json:"author" default:"defaultAuthor"`           // FIXME: figure out if this field is required and default author is ok
	Message string `json:"message" default:"Default commit message"` // FIXME: author/message passibly is CommitInfo or smth like this
}

func (br *BranchRequester) Squash(ctx context.Context, branchID string, options *BranchSquashOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	body := struct {
		CommitInfo any `json:"commit_info"`
	}{*options}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "squash"))
	_, err = doRequest(ctx, sl, nil) // TODO: There is ok response also
	return err
}

type BranchResetOptions struct {
	UsePath bool
}

func (br *BranchRequester) Reset(ctx context.Context, branchID, commit string, options *BranchResetOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	if !options.UsePath {
		path := br.path.(RepoPath)
		commit = CommitPath{
			Organization: path.Organization,
			Database:     path.Database,
			Repo:         path.Repo,
			Branch:       branchID,
			Commit:       commit,
		}.GetPath("reset")
	}
	body := struct {
		Commit string `json:"commit_descriptor"`
	}{commit}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "reset"))
	_, err = doRequest(ctx, sl, nil) // TODO: There is ok response also
	return
}

type BranchApplyOptions struct {
	BeforeCommit    string `json:"before_commit" validate:"required"`
	AfterCommit     string `json:"after_commit" validate:"required"`
	Message         string
	Author          string
	Keep            map[string]string `json:"keep"` // FIXME: figure out correct type
	MatchFinalState bool              `json:"match_final_state"`
}

func (br *BranchRequester) Apply(ctx context.Context, branchID string, options *BranchApplyOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	type commitInfo struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}
	body := struct {
		BranchApplyOptions
		CommitInfo commitInfo `json:"commit_info"`
	}{*options, commitInfo{options.Author, options.Message}}
	sl := br.Client.C.BodyJSON(body).Post(br.getURL(branchID, "apply"))
	_, err = doRequest(ctx, sl, nil) // TODO: There is ok response also
	return
}

type BranchRebaseOptions struct {
	RebaseSource string `json:"rebase_source"` // FIXME: can be commit id or branch id, consider to use several methods such as RebaseToCommit, RebaseToBranch
	Author       string `json:"author" validate:"required" default:"Default author"`
	Message      string `json:"message" validate:"required" default:"Default message"`
}

func (br *BranchRequester) Rebase(ctx context.Context, branchID string, options *BranchRebaseOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := br.Client.C.BodyJSON(options).Post(br.getURL(branchID, "rebase"))
	_, err = doRequest(ctx, sl, nil) // TODO: There is ok response also
	return
}

type BranchCommitLogOptions struct {
	Count int `url:"count" default:"-1"`
	Start int `url:"start,omitempty" default:"0"`
}

// FIXME: check if options are actually used everywhere
func (br *BranchRequester) CommitLog(ctx context.Context, buf *[]Commit, branchID string, options *BranchCommitLogOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := br.Client.C.QueryStruct(options).Get(br.getURL(branchID, "log"))
	_, err = doRequest(ctx, sl, buf)
	return
}

func (br *BranchRequester) Optimize(ctx context.Context, branchID string) error {
	sl := br.Client.C.Post(br.getURL(branchID, "optimize"))
	if _, err := doRequest(ctx, sl, nil); err != nil { // TODO: There is ok response also
		return err
	}

	return nil
}

func (br *BranchRequester) getURL(branchID, action string) string {
	path := br.path.(RepoPath)
	return BranchPath{
		Organization: path.Organization,
		Database:     path.Database,
		Repo:         path.Repo,
		Branch:       branchID,
	}.GetPath(action)
}

type BranchPath struct {
	Organization, Database, Repo, Branch string
}

func (bp BranchPath) GetPath(action string) string {
	suburl := fmt.Sprintf(
		"%s/%s/%s",
		action,
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
