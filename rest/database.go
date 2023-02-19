package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bdragon300/terminusgo/schema"
	"github.com/bdragon300/terminusgo/srverror"
	"github.com/bdragon300/terminusgo/woql/bare"
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

type DatabaseIntroducer struct {
	BaseIntroducer
	ctx context.Context
}

func (di *DatabaseIntroducer) OnOrganization(path OrganizationPath) *DatabaseRequester {
	return &DatabaseRequester{BaseRequester: BaseRequester{Client: di.client, path: path}}
}

func (di *DatabaseIntroducer) WithContext(ctx context.Context) *DatabaseIntroducer {
	r := *di
	r.ctx = ctx
	return &r
}

func (di *DatabaseIntroducer) ListDatabaseInfo(buf *[]DatabaseInfo) (response TerminusResponse, err error) {
	query := map[string]any{"verbose": true, "branches": true}
	sl := di.client.C.QueryStruct(&query).Get("db")
	return doRequest(di.ctx, sl, buf)
}

type DatabaseRequester struct {
	BaseRequester
	dataVersion string
}

func (dr *DatabaseRequester) WithContext(ctx context.Context) *DatabaseRequester {
	r := *dr
	r.ctx = ctx
	return &r
}

func (dr *DatabaseRequester) WithDataVersion(dataVersion string) *DatabaseRequester {
	dr.dataVersion = dataVersion
	return dr
}

func (dr *DatabaseRequester) ListAll(buf *[]Database, userName string) (response TerminusResponse, err error) {
	URL := "/" // Current user databases by default
	if userName != "" {
		URL = fmt.Sprintf(
			"organizations/%s/users/%s/databases",
			url.QueryEscape(dr.path.(OrganizationPath).Organization), url.QueryEscape(userName),
		)
	}
	sl := dr.Client.C.Get(URL)
	return doRequest(dr.ctx, sl, buf)
}

func (dr *DatabaseRequester) Get(buf *Database, name string) (response TerminusResponse, err error) {
	options := map[string]any{"verbose": true, "branches": true}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getURL(name, "db"))
	return doRequest(dr.ctx, sl, buf)
}

type DatabaseCreateOptions struct {
	Schema   bool    `json:"schema" default:"true"`
	Public   bool    `json:"public,omitempty"`
	Comment  string  `json:"comment,omitempty"`
	Prefixes *Prefix `json:"prefixes,omitempty"`
}

func (dr *DatabaseRequester) Create(name, label string, options *DatabaseCreateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		DatabaseCreateOptions
		Label string `json:"label"`
	}{*options, label}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getURL(name, "db"))
	return doRequest(dr.ctx, sl, nil)
}

type DatabaseDeleteOptions struct {
	Force bool `url:"force,omitempty"`
}

func (dr *DatabaseRequester) Delete(name string, options *DatabaseDeleteOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Delete(dr.getURL(name, "db"))
	return doRequest(dr.ctx, sl, nil)
}

func (dr *DatabaseRequester) IsExists(name string) (exists bool, response TerminusResponse, err error) {
	sl := dr.Client.C.Head(dr.getURL(name, "db"))
	response, err = doRequest(dr.ctx, sl, nil)
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

func (dr *DatabaseRequester) Update(name string, options *DatabaseUpdateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.BodyJSON(options).Put(dr.getURL(name, "db"))
	return doRequest(dr.ctx, sl, nil)
}

type DatabaseWOQLOptions struct {
	CommitAuthor  string
	CommitMessage string
	AllWitnesses  bool
}

// Query with database context
func (dr *DatabaseRequester) WOQL(buf *srverror.WOQLResponse, name string, query bare.RawQuery, options *DatabaseWOQLOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	type commitInfo struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}
	body := struct {
		AllWitnesses bool          `json:"all_witnesses,omitempty"`
		CommitInfo   commitInfo    `json:"commit_info"`
		Query        bare.RawQuery `json:"query"`
	}{options.AllWitnesses, commitInfo{options.CommitAuthor, options.CommitMessage}, query}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getURL(name, "woql"))
	return doRequest(dr.ctx, sl, buf)
}

type DatabaseCloneOptions struct {
	Public    bool   `json:"public"`
	RemoteURL string `json:"remote_url" default:"http://example.com/user/test_db"`
	Comment   string `json:"comment" default:"Default comment"`
}

func (dr *DatabaseRequester) Clone(newName, newLabel string, options *DatabaseCloneOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		DatabaseCloneOptions
		Label string `json:"label"`
	}{*options, newLabel}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getURL(newName, "clone"))
	return doRequest(dr.ctx, sl, nil)
}

// TODO: figure out what prefixes are
func (dr *DatabaseRequester) Prefixes(buf *Prefix, dbName string) (response TerminusResponse, err error) {
	sl := dr.Client.C.Get(dr.getURL(dbName, "prefixes"))
	return doRequest(dr.ctx, sl, buf)
}

type DatabaseCommitLogOptions struct {
	Count int `url:"count" default:"-1"`
	Start int `url:"start,omitempty" default:"0"`
}

func (dr *DatabaseRequester) CommitLog(buf *[]Commit, name string, options *DatabaseCommitLogOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getURL(name, "log"))
	return doRequest(dr.ctx, sl, buf)
}

func (dr *DatabaseRequester) Optimize(dbName string) (response TerminusResponse, err error) {
	sl := dr.Client.C.Post(dr.getURL(dbName, "optimize"))
	return doRequest(dr.ctx, sl, nil)
}

type DatabaseSchemaFrameOptions struct {
	CompressIDs    bool `json:"compress_ids" default:"true"`
	ExpandAbstract bool `json:"expand_abstract" default:"true"`
}

func (dr *DatabaseRequester) SchemaFrameAll(buf *[]schema.RawSchemaItem, name string, options *DatabaseSchemaFrameOptions) (response TerminusResponse, err error) {
	var resp map[string]map[string]any
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getURL(name, "schema"))
	response, err = doRequest(dr.ctx, sl, &resp)
	if err != nil {
		return
	}

	for k, v := range resp {
		v["@id"] = k
		*buf = append(*buf, v)
	}
	return
}

func (dr *DatabaseRequester) SchemaFrameType(buf *schema.RawSchemaItem, name, docType string, options *DatabaseSchemaFrameOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	params := struct {
		DatabaseSchemaFrameOptions
		Type string `json:"type"`
	}{*options, docType}
	sl := dr.Client.C.QueryStruct(params).Get(dr.getURL(name, "schema"))
	return doRequest(dr.ctx, sl, buf)
}

type DatabasePackOptions struct {
	RepositoryHead string `json:"repository_head,omitempty"`
}

func (dr *DatabaseRequester) Pack(buf io.Writer, name string, options *DatabasePackOptions) (writtenBytes int64, response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.BodyJSON(options).Post(dr.getURL(name, "pack"))

	var httpReq *http.Request
	if httpReq, err = sl.Request(); err != nil {
		return
	}
	if dr.ctx != nil {
		httpReq = httpReq.WithContext(dr.ctx)
	}
	var httpResp *http.Response
	if httpResp, err = dr.Client.implClient.Do(httpReq); err != nil {
		return
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode >= 300 {
		response = &srverror.TerminusErrorResponse{}
		err = json.NewDecoder(httpResp.Body).Decode(response)
		return
	}
	response = &srverror.TerminusOkResponse{Response: httpResp}
	writtenBytes, err = io.Copy(buf, httpResp.Body)
	return
}

func (dr *DatabaseRequester) UnpackUpload(dbName string, data io.Reader) (readBytes int64, response TerminusResponse, err error) {
	sl := dr.Client.C.Body(data).Post(dr.getURL(dbName, "unpack"))
	response, err = doRequest(dr.ctx, sl, nil)
	return
}

func (dr *DatabaseRequester) UnpackTusResource(dbName, tusLocation string) (readBytes int64, response TerminusResponse, err error) {
	body := map[string]string{"resource_uri": tusLocation}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getURL(dbName, "unpack"))
	response, err = doRequest(dr.ctx, sl, nil)
	return
}

func (dr *DatabaseRequester) getURL(dbName, action string) string {
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
