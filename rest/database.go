package rest

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/bdragon300/terminusgo/srverror"
)

type Database struct {
	ID           string    `json:"@id"`
	Path         string    `json:"path"`
	Type         string    `json:"@type"`
	Name         string    `json:"name"`
	Comment      string    `json:"comment"`
	CreationDate time.Time `json:"creation_date"`
	Label        string    `json:"label"`
	State        string    `json:"state"` // FIXME: maybe its needed enum?
}

type Prefix struct {
	Base   string `json:"@base"`
	Schema string `json:"@schema"`
	Type   string `json:"@type"`
}

type DatabaseIntroducer BaseIntroducer

func (di *DatabaseIntroducer) OnOrganization(path OrganizationPath) *DatabaseRequester {
	return &DatabaseRequester{Client: di.client, path: path}
}

func (di *DatabaseIntroducer) OnUser(path UserPath) *DatabaseRequester {
	return &DatabaseRequester{Client: di.client, path: path}
}

func (di *DatabaseIntroducer) OnServer() *DatabaseRequester {
	return &DatabaseRequester{Client: di.client, path: nil}
}

type DatabaseRequester BaseRequester

// FIXME: test on localhost
func (dr *DatabaseRequester) ListAll(ctx context.Context, buf *[]Database) error {
	var URL string
	switch path := dr.path.(type) {
	case UserPath:
		URL = fmt.Sprintf(
			"organizations/%s/users/%s/databases",
			url.QueryEscape(path.Organization), url.QueryEscape(path.User),
		)
	case OrganizationPath:
		// FIXME: figure out the difference from "/db" endpoint. "/db" returns only "path"
		URL = "db"
	case nil:
		URL = "" // Request on "/"
	default:
		panic("Unknown Path type")
	}
	sl := dr.Client.C.Get(URL)
	if _, err := doRequest(ctx, sl, buf); err != nil {
		return err
	}
	return nil
}

type DatabaseGetOptions struct {
	Verbose  bool `url:"verbose,omitempty"`
	Branches bool `url:"branches,omitempty"`
}

// FIXME: test additionally on localhost
func (dr *DatabaseRequester) Get(ctx context.Context, buf *Database, name string, options *DatabaseGetOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getOrganizationDBURL(name, "db"))
	if _, err = doRequest(ctx, sl, buf); err != nil {
		return err
	}

	return
}

type DatabasePrefixes struct {
	Base   string `json:"@base,omitempty" default:"terminusdb:///data"`
	Schema string `json:"@schema,omitempty" default:"terminusdb:///schema"`
}

type DatabaseCreateOptions struct {
	Schema   bool              `json:"schema,omitempty"`
	Public   bool              `json:"public,omitempty"`
	Label    string            `json:"label" validate:"required" default:"Default label"` // FIXME: check if correct validators specified everywhere
	Comment  string            `json:"comment" validate:"required" default:"Default comment"`
	Prefixes *DatabasePrefixes `json:"prefixes,omitempty"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Create(ctx context.Context, db Database, options *DatabaseCreateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := dr.Client.C.BodyJSON(options).Post(dr.getOrganizationDBURL(db.Name, "db"))
	_, err = doRequest(ctx, sl, nil)
	return
}

type DatabaseDeleteOptions struct {
	Force bool `url:"force,omitempty"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Delete(ctx context.Context, name string, options *DatabaseDeleteOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := dr.Client.C.QueryStruct(options).Delete(dr.getOrganizationDBURL(name, "db"))
	_, err = doRequest(ctx, sl, nil)
	return
}

// FIXME: test on localhost
func (dr *DatabaseRequester) IsExists(ctx context.Context, name string) (bool, error) {
	var res Database
	sl := dr.Client.C.Head(dr.getOrganizationDBURL(name, "db"))
	if _, err := doRequest(ctx, sl, &res); err != nil {
		if errors.Is(err, srverror.TerminusError{}) && err.(srverror.TerminusError).HTTPCode == 404 {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// FIXME: check if "omitempty" is required somewhere in other options structs
type DatabaseUpdateOptions struct {
	Schema   bool              `json:"schema,omitempty"`
	Public   bool              `json:"public,omitempty"`
	Label    string            `json:"label" validate:"required" default:"Default label"`
	Comment  string            `json:"comment" validate:"required" default:"Default comment"`
	Prefixes *DatabasePrefixes `json:"prefixes,omitempty"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Update(ctx context.Context, db Database, options *DatabaseUpdateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := dr.Client.C.BodyJSON(options).Put(dr.getOrganizationDBURL(db.Name, "db"))
	_, err = doRequest(ctx, sl, nil)
	return
}

type DatabaseCloneOptions struct {
	Public    bool   `json:"public"`
	RemoteURL string `json:"remote_url" validate:"required,url" default:"http://example.com/user/test_db"`
	Label     string `json:"label" validate:"required" default:"Default label"` // FIXME: check if such default is correct (and everywhere)
	Comment   string `json:"comment" validate:"required" default:"Default comment"`
}

// FIXME: test on localhost
func (dr *DatabaseRequester) Clone(ctx context.Context, newName string, options *DatabaseCloneOptions) (err error) {
	// TODO: requires to execute on an organization instead of on a user, implement such mechanism of separation
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := dr.Client.C.BodyJSON(options).Put(dr.getOrganizationDBURL(newName, "clone"))
	_, err = doRequest(ctx, sl, nil)
	return
}

// FIXME: additionally test on localhost, figure out what prefixes are
func (dr *DatabaseRequester) Prefixes(ctx context.Context, buf *Prefix, dbName string) error {
	sl := dr.Client.C.Get(dr.getOrganizationDBURL(dbName, "prefixes"))
	_, err := doRequest(ctx, sl, buf)
	return err
}

type DatabaseCommitLogOptions struct {
	Count int `url:"count" default:"-1"`
	Start int `url:"start,omitempty" default:"0"`
}

func (dr *DatabaseRequester) CommitLog(ctx context.Context, buf *[]Commit, name string, options *DatabaseCommitLogOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getOrganizationDBURL(name, "log"))
	_, err = doRequest(ctx, sl, buf)
	return
}

func (dr *DatabaseRequester) Optimize(ctx context.Context, dbName string) error {
	sl := dr.Client.C.Post(dr.getOrganizationDBURL(dbName, "optimize"))
	if _, err := doRequest(ctx, sl, nil); err != nil { // TODO: There is ok response also
		return err
	}

	return nil
}

func (dr *DatabaseRequester) getOrganizationDBURL(dbName, action string) string {
	switch v := dr.path.(type) {
	case OrganizationPath:
		return DatabasePath{
			Organization: v.Organization,
			Database:     dbName,
		}.GetPath(action)
	case UserPath:
		panic("Should not happen") // FIXME: use more nice approach
	case nil:
		panic("Should not happen") // FIXME: use more nice approach
	}
	return ""
}

type DatabasePath struct {
	Organization, Database string
}

func (dp DatabasePath) GetPath(action string) string {
	return fmt.Sprintf("%s/%s", action, getDBBase(dp.Database, dp.Organization)) // FIXME: if org is empty
}
