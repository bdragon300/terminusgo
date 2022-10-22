package objects

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bdragon300/terminusgo/srverror"
	"github.com/creasty/defaults"
	"github.com/dghubble/sling"
	"github.com/go-playground/validator/v10"
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
}

type Validated interface {
	Validate() error
}

type ObjectPathProvider interface {
	GetPath(action string) string
}

func doRequest(ctx context.Context, sling *sling.Sling, okResponse any) (*http.Response, error) {
	req, err := sling.Request()
	if err != nil {
		return nil, err
	}
	errTerminus := new(srverror.TerminusError)
	resp, err := sling.Do(req.WithContext(ctx), okResponse, errTerminus)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 300 {
		errTerminus.HTTPCode = resp.StatusCode
		return resp, errTerminus
	}
	return resp, nil
}

func prepareOptions[T any](options *T) (*T, error) {
	if options == nil {
		options = new(T)
		defaults.MustSet(options)
	}
	validate := validator.New()
	return options, validate.Struct(options)
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
