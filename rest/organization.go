package rest

import (
	"context"
	"fmt"
	"net/url"
)

// TODO: seems that these structures are items of system_schema items and should be aligned with it
//
//	(lack of `child` field below for example) (but no Remote in system_schema for example)
type Organization struct {
	ID       string   `json:"@id"`
	Type     string   `json:"@type"`
	Database []string `json:"database"`
	Name     string   `json:"name"`
}

// TODO: test on local instance
type OrganizationRequester BaseRequester

func (or *OrganizationRequester) WithContext(ctx context.Context) *OrganizationRequester {
	r := *or
	r.ctx = ctx
	return &r
}

func (or *OrganizationRequester) ListAll(buf *[]Organization) (response TerminusResponse, err error) {
	sl := or.Client.C.Get("organizations")
	return doRequest(or.ctx, sl, buf)
}

func (or *OrganizationRequester) Get(buf *Organization, name string) (response TerminusResponse, err error) {
	sl := or.Client.C.Get(or.getURL(name))
	return doRequest(or.ctx, sl, buf)
}

func (or *OrganizationRequester) Create(name string) (response TerminusResponse, err error) {
	sl := or.Client.C.Post(or.getURL(name))
	return doRequest(or.ctx, sl, nil)
}

func (or *OrganizationRequester) Delete(name string) (response TerminusResponse, err error) {
	sl := or.Client.C.Delete(or.getURL(name))
	return doRequest(or.ctx, sl, nil)
}

func (or *OrganizationRequester) getURL(objectID string) string {
	return fmt.Sprintf("organizations/%s", objectID)
}

type OrganizationPath struct {
	Organization string
}

func (op OrganizationPath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, op.GetPath())
}

func (op OrganizationPath) GetPath() string {
	return url.QueryEscape(op.Organization)
}
