package rest

import (
	"context"
	"fmt"
	"net/url"
)

type UserCapability struct {
	ID    string `json:"@id"`
	Role  []Role `json:"role"`
	Scope string `json:"scope"`
}

type User struct {
	ID         string           `json:"@id"`
	Capability []UserCapability `json:"capability"`
	Name       string           `json:"name"`
}

type UserIntroducer BaseIntroducer

func (ui *UserIntroducer) OnOrganization(path OrganizationPath) *UserRequester {
	return &UserRequester{Client: ui.client, path: path}
}

func (ui *UserIntroducer) OnServer() *UserRequester {
	return &UserRequester{Client: ui.client, path: nil}
}

type UserRequester BaseRequester

func (ur *UserRequester) WithContext(ctx context.Context) *UserRequester {
	r := *ur
	r.ctx = ctx
	return &r
}

type UserListAllOptions struct {
	Capability bool `json:"capability" default:"false"`
}

func (ur *UserRequester) ListAll(buf *[]User, options *UserListAllOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := ur.Client.C.QueryStruct(options).Get(ur.getURL(""))
	return doRequest(ur.ctx, sl, buf)
}

type UserGetOptions struct {
	Capability bool `json:"capability" default:"false"`
}

func (ur *UserRequester) Get(buf *User, name string, options *UserListAllOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := ur.Client.C.QueryStruct(options).Get(ur.getURL(name))
	return doRequest(ur.ctx, sl, buf)
}

func (ur *UserRequester) Create(name, password string) (response TerminusResponse, err error) {
	body := struct {
		Name     string `json:"name"`
		Password string `json:"password,omitempty"`
	}{name, password}
	sl := ur.Client.C.BodyJSON(body).Post("users")
	return doRequest(ur.ctx, sl, nil)
}

func (ur *UserRequester) UpdatePassword(name, password string) (response TerminusResponse, err error) {
	body := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{name, password}
	sl := ur.Client.C.BodyJSON(body).Put("users")
	return doRequest(ur.ctx, sl, nil)
}

type UserUpdateCapabilitiesOptions struct {
	Scope     string   `json:"scope"`
	ScopeType string   `json:"scope_type,omitempty" default:"database"` // Only one possible value: "database" or empty
	Roles     []string `json:"roles"`                                   // Role IDs
}

type UserCapabilitiesOperation string

const (
	UserGrantCapabilities  UserCapabilitiesOperation = "grant"
	UserRevokeCapabilities UserCapabilitiesOperation = "revoke"
)

func (ur *UserRequester) UpdateCapabilities(name string, operation UserCapabilitiesOperation, options *UserUpdateCapabilitiesOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		UserUpdateCapabilitiesOptions
		User      string `json:"user"`
		Operation string `json:"operation"`
	}{*options, name, string(operation)}
	sl := ur.Client.C.BodyJSON(body).Post("capabilities")
	return doRequest(ur.ctx, sl, nil)
}

func (ur *UserRequester) getURL(objectID string) string {
	org := ""
	action := "users"
	if ur.path != nil {
		path := ur.path.(OrganizationPath)
		org = path.Organization
		action = "organizations"
	}
	return UserPath{
		Organization: org,
		User:         objectID,
	}.GetURL(action)
}

type UserPath struct {
	Organization, User string
}

func (up UserPath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, up.GetPath())
}

func (up UserPath) GetPath() string {
	if up.Organization == "" {
		return url.QueryEscape(up.User)
	}
	return fmt.Sprintf("%s/users/%s", url.QueryEscape(up.Organization), url.QueryEscape(up.User))
}
