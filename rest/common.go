package rest

import (
	"context"
	"fmt"
	"net/url"
	"reflect"

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
	path   TerminusObjectPath
	ctx    context.Context
}

type TerminusResponse interface {
	IsOK() bool
}

type TerminusObjectPath interface {
	fmt.Stringer
	GetURL(action string) string
}

func doRequest(ctx context.Context, sling *sling.Sling, okResponse any) (TerminusResponse, error) {
	req, err := sling.Request()
	if err != nil {
		return nil, err
	}
	errResp := &srverror.TerminusError{}
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
		return errResp, errResp
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

func getDatabasePath(organization, database string) string {
	if database == DatabaseSystem {
		return DatabaseSystem
	}
	if organization == "" {
		organization = "NoOrganization"
	}
	return fmt.Sprintf("%s/%s", url.PathEscape(organization), url.PathEscape(database))
}

func fillUnescapedStringFields(vals []string, buf any) {
	buff := reflect.ValueOf(buf).Elem()
	typ := buff.Type()
	for i := 0; i < typ.NumField() && len(vals) > 0; i++ {
		fld := typ.Field(i)
		if !fld.IsExported() || fld.Type.Kind() != reflect.String {
			continue
		}
		s := vals[0]
		if us, err := url.PathUnescape(vals[0]); err == nil {
			s = us
		}
		buff.Field(i).SetString(s)
		vals = vals[1:]
	}
}
