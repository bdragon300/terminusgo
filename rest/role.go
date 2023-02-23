package rest

import (
	"context"
	"fmt"
	"net/url"
)

type Role struct {
	ID     string   `json:"@id"`
	Action []string `json:"action"`
	Name   string   `json:"name"`
}

type RoleRequester BaseRequester

func (rr *RoleRequester) WithContext(ctx context.Context) *RoleRequester {
	r := *rr
	r.ctx = ctx
	return &r
}

func (rr *RoleRequester) ListAll(buf *[]Role) (response TerminusResponse, err error) {
	sl := rr.Client.C.Get("roles")
	return doRequest(rr.ctx, sl, buf)
}

func (rr *RoleRequester) Get(name string, buf *Role) (response TerminusResponse, err error) {
	sl := rr.Client.C.Get(rr.getURL(name))
	return doRequest(rr.ctx, sl, buf)
}

type RoleCreateOptions struct {
	Action []string `json:"action"`
}

func (rr *RoleRequester) Create(name string, options *RoleCreateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		RoleCreateOptions
		Name string `json:"name"`
	}{*options, name}
	sl := rr.Client.C.BodyJSON(body).Post("roles")
	return doRequest(rr.ctx, sl, nil)
}

type RoleUpdateOptions struct {
	Action []string `json:"action"`
}

func (rr *RoleRequester) Update(name string, options *RoleUpdateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		RoleUpdateOptions
		Name string `json:"name"`
	}{*options, name}
	sl := rr.Client.C.BodyJSON(body).Put("roles")
	return doRequest(rr.ctx, sl, nil)
}

func (rr *RoleRequester) Delete(name string) (response TerminusResponse, err error) {
	sl := rr.Client.C.Delete(rr.getURL(name))
	return doRequest(rr.ctx, sl, nil)
}

func (rr *RoleRequester) getURL(objectID string) string {
	return fmt.Sprintf("roles/%s", url.QueryEscape(objectID))
}
