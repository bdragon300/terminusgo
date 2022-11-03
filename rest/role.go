package rest

import (
	"context"
	"fmt"
	"net/url"
)

type Role struct {
	ID     string   `json:"@id"`
	Type   string   `json:"type"`
	Action []string `json:"action"`
	Name   string   `json:"name"`
}

type RoleRequester BaseRequester

func (rr *RoleRequester) ListAll(ctx context.Context, buf *[]Role) error {
	sl := rr.Client.C.Get("roles")
	if _, err := doRequest(ctx, sl, buf); err != nil {
		return err
	}
	return nil
}

func (rr *RoleRequester) Get(ctx context.Context, buf *Role, name string) error {
	sl := rr.Client.C.Get(rr.getURL(name))
	if _, err := doRequest(ctx, sl, buf); err != nil {
		return err
	}
	return nil
}

type RoleCreateOptions struct {
	Action []string `json:"action"`
}

func (rr *RoleRequester) Create(ctx context.Context, name string, options *RoleCreateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	body := struct {
		RoleCreateOptions
		Name string `json:"name"`
	}{*options, name}
	sl := rr.Client.C.BodyJSON(body).Post("roles")
	if _, err = doRequest(ctx, sl, nil); err != nil { // TODO: there is ok response
		return err
	}

	return nil
}

type RoleUpdateOptions struct {
	Action []string `json:"action"`
}

func (rr *RoleRequester) Update(ctx context.Context, name string, options *RoleUpdateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	body := struct {
		RoleUpdateOptions
		Name string `json:"name"`
	}{*options, name}
	sl := rr.Client.C.BodyJSON(body).Put("roles")
	if _, err = doRequest(ctx, sl, nil); err != nil { // TODO: there is ok response
		return err
	}

	return
}

func (rr *RoleRequester) Delete(ctx context.Context, name string) error {
	sl := rr.Client.C.Delete(rr.getURL(name))
	if _, err := doRequest(ctx, sl, nil); err != nil {
		return err
	}
	return nil
}

func (rr *RoleRequester) getURL(objectID string) string {
	return fmt.Sprintf("roles/%s", url.QueryEscape(objectID))
}
