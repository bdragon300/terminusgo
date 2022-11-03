package rest

import (
	"context"
	"fmt"
	"net/url"
)

type UserCapability struct {
	ID    string `json:"@id"`
	Type  string `json:"type"`
	Role  []Role `json:"role"`
	Scope string `json:"scope"`
}

type User struct {
	ID         string           `json:"@id"`
	Type       string           `json:"type"` // FIXME: actually missing... check
	Capability []UserCapability `json:"capability"`
	Name       string           `json:"name"`
}

type UserIntroducer BaseIntroducer

func (ui *UserIntroducer) OnOrganization(path OrganizationPath) *UserRequester {
	return &UserRequester{Client: ui.client, path: path}
}

// TODO: test on localhost
func (ui *UserIntroducer) OnServer() *UserRequester {
	return &UserRequester{Client: ui.client, path: nil}
}

type UserRequester BaseRequester

// TODO: test on localhost
func (ur *UserRequester) ListAll(ctx context.Context, buf *[]User) error {
	sl := ur.Client.C.Get(ur.getURL("")) // FIXME: hack, make smth like getListUrl
	_, err := doRequest(ctx, sl, buf)
	return err
}

// TODO: test on localhost
func (ur *UserRequester) Get(ctx context.Context, buf *User, name string) error {
	sl := ur.Client.C.Get(ur.getURL(name))
	_, err := doRequest(ctx, sl, buf)
	return err
}

type UserCreateOptions struct {
	Password string `json:"password" validate:"required"`
}

// TODO: test on localhost
func (ur *UserRequester) Create(ctx context.Context, name string, options *UserCreateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	body := struct {
		UserCreateOptions
		Name string `json:"name"`
	}{*options, name}
	sl := ur.Client.C.BodyJSON(body).Post("users")
	if _, err = doRequest(ctx, sl, nil); err != nil { // TODO: there is ok response
		return err
	}
	return
}

// TODO: test on localhost
func (ur *UserRequester) UpdatePassword(ctx context.Context, name, password string) error {
	body := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{name, password}
	sl := ur.Client.C.BodyJSON(body).Put("users")
	if _, err := doRequest(ctx, sl, nil); err != nil { // TODO: there is ok response
		return err
	}

	return nil
}

type UserUpdateCapabilitiesOptions struct {
	Operation string   `json:"operation" default:"revoke"` // FIXME: make enum; ensure that revoke is ok as a default value
	Scope     string   `json:"scope"`                      // FIXME: figure out a default value for scope
	Roles     []string `json:"roles"`                      // FIXME: make Roles type
}

// TODO: test on localhost
func (ur *UserRequester) UpdateCapabilities(ctx context.Context, name string, options *UserUpdateCapabilitiesOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	body := struct {
		UserUpdateCapabilitiesOptions
		User string `json:"user"`
	}{*options, name}
	sl := ur.Client.C.BodyJSON(body).Post("capabilities")
	if _, err = doRequest(ctx, sl, nil); err != nil { // TODO: there is ok response
		return err
	}

	return
}

func (ur *UserRequester) getURL(objectID string) string {
	org := ""
	if ur.path != nil {
		path := ur.path.(OrganizationPath)
		org = path.Organization
	}
	return UserPath{
		Organization: org,
		User:         objectID,
	}.GetPath("organizations")
}

type UserPath struct {
	Organization, User string
}

func (up UserPath) GetPath(action string) string {
	if up.Organization == "" {
		return fmt.Sprintf("%s/%s", action, url.QueryEscape(up.User))
	}
	return fmt.Sprintf("%s/%s/users/%s", action, url.QueryEscape(up.Organization), url.QueryEscape(up.User))
}
