package rest

import (
	"context"
	"fmt"
	"net/url"
)

type RoleAction string

const (
	RoleActionCreateDatabase      RoleAction = "create_database"
	RoleActionDeleteDatabase      RoleAction = "delete_database"
	RoleActionClassFrame          RoleAction = "class_frame"
	RoleActionClone               RoleAction = "clone"
	RoleActionFetch               RoleAction = "fetch"
	RoleActionPush                RoleAction = "push"
	RoleActionBranch              RoleAction = "branch"
	RoleActionRebase              RoleAction = "rebase"
	RoleActionInstanceReadAccess  RoleAction = "instance_read_access"
	RoleActionInstanceWriteAccess RoleAction = "instance_write_access"
	RoleActionSchemaReadAccess    RoleAction = "schema_read_access"
	RoleActionSchemaWriteAccess   RoleAction = "schema_write_access"
	RoleActionMetaReadAccess      RoleAction = "meta_read_access"
	RoleActionMetaWriteAccess     RoleAction = "meta_write_access"
	RoleActionCommitReadAccess    RoleAction = "commit_read_access"
	RoleActionCommitWriteAccess   RoleAction = "commit_write_access"
	RoleActionManageCapabilities  RoleAction = "manage_capabilities"
)

type Role struct {
	ID     string       `json:"@id"`
	Action []RoleAction `json:"action"`
	Name   string       `json:"name"`
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
	Action []RoleAction `json:"action"`
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
	Action []RoleAction `json:"action"`
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
