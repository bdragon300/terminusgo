package objects

import (
	"context"
	"fmt"
	"net/url"
)

type Organization struct {
	ID       string   `json:"@id"`
	Type     string   `json:"@type"`
	Database []string `json:"database"` // TODO: Figure out why it's a list of databases
	Name     string   `json:"name"`
}

// TODO: test on local instance
type OrganizationRequester BaseRequester

func (or *OrganizationRequester) ListAll(ctx context.Context, buf *[]Organization) error {
	sl := or.Client.C.Get("organizations")
	if _, err := doRequest(ctx, sl, buf); err != nil {
		return err
	}
	return nil
}

func (or *OrganizationRequester) Get(ctx context.Context, buf *Organization, name string) error {
	sl := or.Client.C.Get(or.getURL(name))
	if _, err := doRequest(ctx, sl, buf); err != nil {
		return err
	}
	return nil
}

func (or *OrganizationRequester) Create(ctx context.Context, name string) error {
	sl := or.Client.C.Post(or.getURL(name))
	_, err := doRequest(ctx, sl, nil)
	return err
}

func (or *OrganizationRequester) Delete(ctx context.Context, name string) error {
	sl := or.Client.C.Delete(or.getURL(name))
	_, err := doRequest(ctx, sl, nil)
	return err
}

func (or *OrganizationRequester) getURL(objectID string) string {
	return fmt.Sprintf("organizations/%s", objectID)
}

type OrganizationPath struct {
	Organization string
}

func (op OrganizationPath) GetPath(action string) string {
	return fmt.Sprintf("%s/%s", action, url.QueryEscape(op.Organization))
}
