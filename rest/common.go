package rest

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bdragon300/terminusgo/srverror"
	"github.com/creasty/defaults"
	"github.com/dghubble/sling"
)

const (
	DatabaseSystem = "_system"
	RepoMeta       = "_meta"
	RepoLocal      = "local"
	BranchCommits  = "_commits"
)

type BaseIntroducer struct {
	client *Client
}

type BaseRequester struct {
	Client *Client
	path   ObjectPathProvider
	ctx    context.Context
}

type TerminusResponse interface {
	IsOK() bool
}

type ObjectPathProvider interface {
	GetURL(action string) string
	GetPath() string
}

func doRequest(ctx context.Context, sling *sling.Sling, okResponse any) (TerminusResponse, error) {
	req, err := sling.Request()
	if err != nil {
		return nil, err
	}
	errResp := &srverror.TerminusErrorResponse{}
	okResp := &srverror.TerminusOkResponse{}
	if okResponse == nil {
		okResponse = okResp
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := sling.Do(req, okResponse, errResp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		errResp.Response = resp
		return errResp, nil
	}
	okResp.Response = resp
	return okResp, nil
}

func prepareOptions[T any](options *T) (*T, error) {
	if options == nil {
		options = new(T)
		defaults.MustSet(options)
	}
	return options, nil
}

func getDBBase(dbName, organization string) string {
	if dbName == DatabaseSystem {
		return DatabaseSystem
	}
	org := organization
	if org == "" {
		org = "NoOrganization"
	}
	return fmt.Sprintf("%s/%s", url.QueryEscape(org), url.QueryEscape(dbName))
}
