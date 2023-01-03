package rest

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type Database struct {
	ID           string    `json:"@id"`
	Type         string    `json:"@type"`
	Name         string    `json:"name"`
	Comment      string    `json:"comment"`
	CreationDate time.Time `json:"creation_date"`
	Label        string    `json:"label"`
	State        string    `json:"state"`
	Path         string    `json:"path"`
	Branches     []string  `json:"branches"`
}

type Prefix struct {
	Base   string `json:"@base"`
	Schema string `json:"@schema"`
	Type   string `json:"@type"`
}

type DatabaseInfo struct {
	ID           string    `json:"@id"`
	Type         string    `json:"@type"`
	Name         string    `json:"name"`
	Comment      string    `json:"comment"`
	CreationDate time.Time `json:"creation_date"`
	Label        string    `json:"label"`
	State        string    `json:"state"`
}

type DatabaseIntroducer BaseIntroducer

func (di *DatabaseIntroducer) OnOrganization(path OrganizationPath) *DatabaseRequester {
	return &DatabaseRequester{Client: di.client, path: path}
}

func (di *DatabaseIntroducer) ListDatabaseInfo(ctx context.Context, buf *[]DatabaseInfo) (response TerminusResponse, err error) {
	query := map[string]any{"verbose": true, "branches": true}
	sl := di.client.C.QueryStruct(&query).Get("db")
	return doRequest(ctx, sl, buf)
}

type DatabaseRequester BaseRequester

// FIXME: test on localhost
func (dr *DatabaseRequester) ListAll(ctx context.Context, buf *[]Database, userName string) (response TerminusResponse, err error) {
	URL := "/" // Current user databases by default
	if userName != "" {
		URL = fmt.Sprintf(
			"organizations/%s/users/%s/databases",
			url.QueryEscape(dr.path.(OrganizationPath).Organization), url.QueryEscape(userName),
		)
	}
	sl := dr.Client.C.Get(URL)
	return doRequest(ctx, sl, buf)
}

// FIXME: test additionally on localhost
func (dr *DatabaseRequester) Get(ctx context.Context, buf *Database, name string) (response TerminusResponse, err error) {
	options := map[string]any{"verbose": true, "branches": true}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getOrganizationDBURL(name, "db"))
	return doRequest(ctx, sl, buf)
}

type DatabaseCreateOptions struct {
	Schema   bool    `json:"schema" default:"true"`
	Public   bool    `json:"public,omitempty"`
	Comment  string  `json:"comment,omitempty"`
	Prefixes *Prefix `json:"prefixes,omitempty"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Create(ctx context.Context, db Database, label string, options *DatabaseCreateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		DatabaseCreateOptions
		Label string `json:"label"`
	}{*options, label}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getOrganizationDBURL(db.Name, "db"))
	return doRequest(ctx, sl, nil)
}

type DatabaseDeleteOptions struct {
	Force bool `url:"force,omitempty"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Delete(ctx context.Context, name string, options *DatabaseDeleteOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Delete(dr.getOrganizationDBURL(name, "db"))
	return doRequest(ctx, sl, nil)
}

// FIXME: test on localhost
func (dr *DatabaseRequester) IsExists(ctx context.Context, name string) (exists bool, response TerminusResponse, err error) {
	sl := dr.Client.C.Head(dr.getOrganizationDBURL(name, "db"))
	response, err = doRequest(ctx, sl, nil)
	if err != nil {
		return
	}
	exists = response.IsOK()
	return
}

type DatabaseUpdateOptions struct {
	Schema   bool    `json:"schema,omitempty"`
	Public   bool    `json:"public,omitempty"`
	Label    string  `json:"label,omitempty"`
	Comment  string  `json:"comment,omitempty"`
	Prefixes *Prefix `json:"prefixes,omitempty"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Update(ctx context.Context, db Database, options *DatabaseUpdateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.BodyJSON(options).Put(dr.getOrganizationDBURL(db.Name, "db"))
	return doRequest(ctx, sl, nil)
}

type DatabaseCloneOptions struct {
	Public    bool   `json:"public"`
	RemoteURL string `json:"remote_url" default:"http://example.com/user/test_db"`
	Comment   string `json:"comment" default:"Default comment"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Clone(ctx context.Context, newName, newLabel string, options *DatabaseCloneOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		DatabaseCloneOptions
		Label string `json:"label"`
	}{*options, newLabel}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getOrganizationDBURL(newName, "clone"))
	return doRequest(ctx, sl, nil)
}

// FIXME: additionally test on localhost, figure out what prefixes are
func (dr *DatabaseRequester) Prefixes(ctx context.Context, buf *Prefix, dbName string) (response TerminusResponse, err error) {
	sl := dr.Client.C.Get(dr.getOrganizationDBURL(dbName, "prefixes"))
	return doRequest(ctx, sl, buf)
}

type DatabaseCommitLogOptions struct {
	Count int `url:"count" default:"-1"`
	Start int `url:"start,omitempty" default:"0"`
}

func (dr *DatabaseRequester) CommitLog(ctx context.Context, buf *[]Commit, name string, options *DatabaseCommitLogOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getOrganizationDBURL(name, "log"))
	return doRequest(ctx, sl, buf)
}

func (dr *DatabaseRequester) Optimize(ctx context.Context, dbName string) (response TerminusResponse, err error) {
	sl := dr.Client.C.Post(dr.getOrganizationDBURL(dbName, "optimize"))
	return doRequest(ctx, sl, nil)
}

func (dr *DatabaseRequester) getOrganizationDBURL(dbName, action string) string {
	return DatabasePath{
		Organization: dr.path.(OrganizationPath).Organization,
		Database:     dbName,
	}.GetURL(action)
}

type DatabasePath struct {
	Organization, Database string
}

func (dp DatabasePath) GetURL(action string) string {
	return fmt.Sprintf("%s/%s", action, dp.GetPath())
}

func (dp DatabasePath) GetPath() string {
	return getDBBase(dp.Database, dp.Organization)
}
