package rest

import (
	"bytes"
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

type DatabaseIntroducer BaseIntroducer

func (di *DatabaseIntroducer) OnOrganization(path OrganizationPath) *DatabaseRequester {
	return &DatabaseRequester{BaseRequester: BaseRequester{Client: di.client, path: path}}
}

func (di *DatabaseIntroducer) ListDatabaseInfo(ctx context.Context, buf *[]DatabaseInfo) (response TerminusResponse, err error) {
	query := map[string]any{"verbose": true, "branches": true}
	sl := di.client.C.QueryStruct(&query).Get("db")
	return doRequest(ctx, sl, buf)
}

type DatabaseRequester struct {
	BaseRequester
	dataVersion string
}

func (dr *DatabaseRequester) WithDataVersion(dataVersion string) *DatabaseRequester {
	dr.dataVersion = dataVersion
	return dr
}

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

func (dr *DatabaseRequester) Get(ctx context.Context, buf *Database, name string) (response TerminusResponse, err error) {
	options := map[string]any{"verbose": true, "branches": true}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getURL(name, "db"))
	return doRequest(ctx, sl, buf)
}

type DatabaseCreateOptions struct {
	Schema   bool    `json:"schema" default:"true"`
	Public   bool    `json:"public,omitempty"`
	Comment  string  `json:"comment,omitempty"`
	Prefixes *Prefix `json:"prefixes,omitempty"`
}

func (dr *DatabaseRequester) Create(ctx context.Context, name, label string, options *DatabaseCreateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		DatabaseCreateOptions
		Label string `json:"label"`
	}{*options, label}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getURL(name, "db"))
	return doRequest(ctx, sl, nil)
}

type DatabaseDeleteOptions struct {
	Force bool `url:"force,omitempty"`
}

func (dr *DatabaseRequester) Delete(ctx context.Context, name string, options *DatabaseDeleteOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Delete(dr.getURL(name, "db"))
	return doRequest(ctx, sl, nil)
}

func (dr *DatabaseRequester) IsExists(ctx context.Context, name string) (exists bool, response TerminusResponse, err error) {
	sl := dr.Client.C.Head(dr.getURL(name, "db"))
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

func (dr *DatabaseRequester) Update(ctx context.Context, name string, options *DatabaseUpdateOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.BodyJSON(options).Put(dr.getURL(name, "db"))
	return doRequest(ctx, sl, nil)
}

type DatabaseWOQLOptions struct {
	CommitAuthor  string
	CommitMessage string
	AllWitnesses  bool
}

// Query with database context
func (dr *DatabaseRequester) WOQL(ctx context.Context, buf *srverror.WOQLResponse, name string, query bare.RawQuery, options *DatabaseWOQLOptions) (response TerminusResponse, err error) {
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
	return doRequest(ctx, sl, buf)
}

type DatabaseCloneOptions struct {
	Public    bool   `json:"public"`
	RemoteURL string `json:"remote_url" default:"http://example.com/user/test_db"`
	Comment   string `json:"comment" default:"Default comment"`
}

func (dr *DatabaseRequester) Clone(ctx context.Context, newName, newLabel string, options *DatabaseCloneOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	body := struct {
		DatabaseCloneOptions
		Label string `json:"label"`
	}{*options, newLabel}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getURL(newName, "clone"))
	return doRequest(ctx, sl, nil)
}

// TODO: figure out what prefixes are
func (dr *DatabaseRequester) Prefixes(ctx context.Context, buf *Prefix, dbName string) (response TerminusResponse, err error) {
	sl := dr.Client.C.Get(dr.getURL(dbName, "prefixes"))
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
	sl := dr.Client.C.QueryStruct(options).Get(dr.getURL(name, "log"))
	return doRequest(ctx, sl, buf)
}

func (dr *DatabaseRequester) Optimize(ctx context.Context, dbName string) (response TerminusResponse, err error) {
	sl := dr.Client.C.Post(dr.getURL(dbName, "optimize"))
	return doRequest(ctx, sl, nil)
}

type DatabaseSchemaFrameOptions struct {
	CompressIDs    bool `json:"compress_ids" default:"true"`
	ExpandAbstract bool `json:"expand_abstract" default:"true"`
}

func (dr *DatabaseRequester) SchemaFrameAll(ctx context.Context, buf *[]schema.RawSchemaItem, name string, options *DatabaseSchemaFrameOptions) (response TerminusResponse, err error) {
	var resp map[string]map[string]any
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Get(dr.getURL(name, "schema"))
	response, err = doRequest(ctx, sl, &resp)
	if err != nil {
		return
	}

	for k, v := range resp {
		v["@id"] = k
		*buf = append(*buf, v)
	}
	return
}

func (dr *DatabaseRequester) SchemaFrameType(ctx context.Context, buf *schema.RawSchemaItem, name, docType string, options *DatabaseSchemaFrameOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	params := struct {
		DatabaseSchemaFrameOptions
		Type string `json:"type"`
	}{*options, docType}
	sl := dr.Client.C.QueryStruct(params).Get(dr.getURL(name, "schema"))
	return doRequest(ctx, sl, buf)
}

type DatabasePackOptions struct {
	RepositoryHead string `json:"repository_head,omitempty"`
}

func (dr *DatabaseRequester) Pack(ctx context.Context, buf io.Writer, name string, options *DatabasePackOptions) (writtenBytes int64, response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.BodyJSON(options).Post(dr.getURL(name, "pack"))

	var httpReq *http.Request
	if httpReq, err = sl.Request(); err != nil {
		return
	}
	var httpResp *http.Response
	if httpResp, err = dr.Client.implClient.Do(httpReq.WithContext(ctx)); err != nil {
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

func (dr *DatabaseRequester) Unpack(ctx context.Context, data io.Reader, name string) (readBytes int64, response TerminusResponse, err error) {
	buf := bytes.NewBuffer(nil)
	if readBytes, err = buf.ReadFrom(data); err != nil {
		return
	}
	sl := dr.Client.C.Body(buf).Post(dr.getURL(name, "unpack"))
	response, err = doRequest(ctx, sl, nil)
	return
}

func (dr *DatabaseRequester) UnpackResourceURI(ctx context.Context, name, resourceURI string) (readBytes int64, response TerminusResponse, err error) {
	body := map[string]string{"resource_uri": resourceURI}
	sl := dr.Client.C.BodyJSON(body).Post(dr.getURL(name, "unpack"))
	response, err = doRequest(ctx, sl, nil)
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
