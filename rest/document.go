package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bdragon300/terminusgo/srverror"
)

// TODO: document can be also a list, a schema for instance
// https://terminusdb.com/docs/guides/reference-guides/json-diff-and-patch#patch-examples-using-curl
type GenericDocument map[string]any

type GraphTypes string

const (
	GraphTypeInstance GraphTypes = "instance"
	GraphTypeSchema   GraphTypes = "schema"
)

// TODO: Mention that DocumentT must have proper json tags to be unmarshalled from response
type DocumentIntroducer[DocumentT any] BaseIntroducer

func (di *DocumentIntroducer[DocumentT]) OnDatabase(path DatabasePath) *DocumentRequester[DocumentT] {
	return &DocumentRequester[DocumentT]{BaseRequester: BaseRequester{Client: di.client, path: path}}
}

func (di *DocumentIntroducer[DocumentT]) OnBranch(path BranchPath) *DocumentRequester[DocumentT] {
	return &DocumentRequester[DocumentT]{BaseRequester: BaseRequester{Client: di.client, path: path}}
}

func (di *DocumentIntroducer[DocumentT]) OnCommit(path CommitPath) *DocumentRequester[DocumentT] {
	return &DocumentRequester[DocumentT]{BaseRequester: BaseRequester{Client: di.client, path: path}}
}

type DocumentRequester[DocumentT any] struct {
	BaseRequester
	dataVersion string
}

func (dr *DocumentRequester[DocumentT]) WithContext(ctx context.Context) *DocumentRequester[DocumentT] {
	r := *dr
	r.ctx = ctx
	return &r
}

func (dr *DocumentRequester[DocumentT]) WithDataVersion(dataVersion string) *DocumentRequester[DocumentT] {
	dr.dataVersion = dataVersion
	return dr
}

type DocumentIsExistsOptions struct {
	GraphType GraphTypes `url:"graph_type" default:"instance"`
}

func (dr *DocumentRequester[DocumentT]) IsExists(options *DocumentIsExistsOptions) (exists bool, response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.QueryStruct(options).Head(dr.path.GetURL("document"))
	if dr.dataVersion != "" {
		sl = sl.Set(srverror.DataVersionHeader, dr.dataVersion)
	}
	response, err = doRequest(dr.ctx, sl, nil)
	if err != nil {
		return
	}
	exists = response.IsOK()
	return
}

type DocumentListOptions struct {
	CompressIDs bool       `url:"compress_ids" default:"true"`
	Type        string     `url:"type,omitempty"`
	Unfold      bool       `url:"unfold" default:"true"`
	Count       int        `url:"count,omitempty"`
	Skip        int        `url:"skip" default:"0"`
	GraphType   GraphTypes `url:"graph_type" default:"instance"`
	Prefixed    bool       `url:"prefixed" default:"true"`
}

func (dr *DocumentRequester[DocumentT]) ListAll(buf *[]DocumentT, options *DocumentListOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	extraParams := struct {
		DocumentListOptions
		AsList    bool `url:"as_list"`   // Feed the response in JSON instead of Concatenated JSON
		Minimized bool `url:"minimized"` // Compress json on the server side
	}{*options, true, true}
	sl := dr.Client.C.QueryStruct(extraParams).Get(dr.path.GetURL("document"))
	if dr.dataVersion != "" {
		sl = sl.Set(srverror.DataVersionHeader, dr.dataVersion)
	}
	return doRequest(dr.ctx, sl, buf)
}

func (dr *DocumentRequester[DocumentT]) ListAllIterator(items chan<- DocumentT, options *DocumentListOptions) (response *http.Response, err error) {
	if items == nil {
		panic("items channel cannot be nil")
	}
	if options, err = prepareOptions(options); err != nil {
		return
	}
	extraParams := struct {
		AsList    bool `url:"as_list"`   // Feed the response in JSON instead of Concatenated JSON
		Minimized bool `url:"minimized"` // Compress json on the server side
	}{false, true}
	sl := dr.Client.C.QueryStruct(options).QueryStruct(&extraParams).Get(dr.path.GetURL("document"))
	if dr.dataVersion != "" {
		sl = sl.Set(srverror.DataVersionHeader, dr.dataVersion)
	}
	req, err := sl.Request()
	if err != nil {
		return nil, err
	}
	if dr.ctx != nil {
		req = req.WithContext(dr.ctx)
	}
	// Making a request using implClient since sl drains and closes resp.Body after the request has been made
	if response, err = dr.Client.implClient.Do(req); err != nil {
		return
	}

	go func() {
		defer func() { _ = response.Body.Close() }()
		defer func() { _, _ = io.Copy(io.Discard, response.Body) }()
		defer func() { _ = recover() }() // FIXME: write to log
		decoder := json.NewDecoder(response.Body)
		for {
			var doc DocumentT
			if err := decoder.Decode(&doc); err == io.EOF {
				close(items)
				break
			} else if err != nil {
				close(items)
				panic(fmt.Sprintf("Error while decoding json response: %s", err))
			}
			items <- doc // TODO: fix externally closed channel
		}
	}()

	return
}

type DocumentGetOptions struct {
	CompressIDs bool       `url:"compress_ids" default:"true"`
	Type        string     `url:"type,omitempty"`
	Unfold      bool       `url:"unfold" default:"true"`
	GraphType   GraphTypes `url:"graph_type" default:"instance"`
	Prefixed    bool       `url:"prefixed" default:"true"`
}

func (dr *DocumentRequester[DocumentT]) Get(docID string, buf *interface{}, options *DocumentGetOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	extraParams := struct {
		DocumentGetOptions
		AsList    bool   `url:"as_list"`
		Minimized bool   `url:"minimized"`
		ID        string `url:"id"`
	}{*options, false, true, docID}
	sl := dr.Client.C.QueryStruct(&extraParams).Get(dr.path.GetURL("document"))
	if dr.dataVersion != "" {
		sl = sl.Set(srverror.DataVersionHeader, dr.dataVersion)
	}
	return doRequest(dr.ctx, sl, buf)
}

type DocumentCreateOptions struct {
	GraphType   GraphTypes `url:"graph_type" default:"instance"`
	Message     string     `url:"message" default:"Default message"`
	Author      string     `url:"author" default:"Default author"`
	RawJSON     bool       `url:"raw_json,omitempty"`
	FullReplace bool       `url:"full_replace,omitempty"`
}

func (dr *DocumentRequester[DocumentT]) Create(doc DocumentT, options *DocumentCreateOptions) (response TerminusResponse, err error) {
	var docSlice []DocumentT
	docSlice = append(docSlice, doc)
	_, response, err = dr.CreateBulk(docSlice, options)
	return
}

func (dr *DocumentRequester[DocumentT]) CreateBulk(docs []DocumentT, options *DocumentCreateOptions) (insertedIDs []string, response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := dr.Client.C.QueryStruct(options).BodyJSON(docs).Post(dr.path.GetURL("document"))
	if dr.dataVersion != "" {
		sl = sl.Set(srverror.DataVersionHeader, dr.dataVersion)
	}
	response, err = doRequest(dr.ctx, sl, &insertedIDs) // FIXME: figure out actual response schema
	return
}

type DocumentUpdateOptions struct {
	GraphType GraphTypes `url:"graph_type" default:"instance"` //  FIXME: check all params in options everywhere if they have to enums
	Message   string     `url:"message" default:"Default message"`
	Author    string     `url:"author" default:"Default author"`
	RawJSON   bool       `url:"raw_json,omitempty"`
	Create    bool       `url:"create,omitempty"`
}

func (dr *DocumentRequester[DocumentT]) UpdateBulk(docs []DocumentT, options *DocumentUpdateOptions) (updatedIDs []string, response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := dr.Client.C.QueryStruct(options).BodyJSON(docs).Put(dr.path.GetURL("document"))
	if dr.dataVersion != "" {
		sl = sl.Set(srverror.DataVersionHeader, dr.dataVersion)
	}
	response, err = doRequest(dr.ctx, sl, &updatedIDs) // FIXME: figure out actual response schema
	return
}

func (dr *DocumentRequester[DocumentT]) Update(doc DocumentT, options *DocumentUpdateOptions) (response TerminusResponse, err error) {
	var docSlice []DocumentT
	docSlice = append(docSlice, doc)
	_, response, err = dr.UpdateBulk(docSlice, options)
	return
}

type DocumentDeleteOptions struct {
	GraphType GraphTypes `url:"graph_type" default:"instance"` //  FIXME: check all params in options everywhere if they have to enums
	Message   string     `url:"message" default:"Default message"`
	Author    string     `url:"author" default:"Default author"`
	Nuke      bool       `url:"nuke,omitempty"`
	ID        string     `url:"id,omitempty"`
}

func (dr *DocumentRequester[DocumentT]) DeleteBulk(docIDs []string, options *DocumentDeleteOptions) (deletedIDs []string, response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := dr.Client.C.QueryStruct(options).BodyJSON(docIDs).Delete(dr.path.GetURL("document"))
	if dr.dataVersion != "" {
		sl = sl.Set(srverror.DataVersionHeader, dr.dataVersion)
	}
	response, err = doRequest(dr.ctx, sl, &deletedIDs) // FIXME: figure out actual response schema
	return
}

func (dr *DocumentRequester[DocumentT]) Delete(docID string, options *DocumentDeleteOptions) (response TerminusResponse, err error) {
	var docSlice []string
	docSlice = append(docSlice, docID)
	_, response, err = dr.DeleteBulk(docSlice, options)
	return
}
