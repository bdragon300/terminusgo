package objects

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// TODO: document can be also a list, a schema for instance
// https://terminusdb.com/docs/guides/reference-guides/json-diff-and-patch#patch-examples-using-curl
type GenericDocument map[string]any

type GraphTypes string

const (
	GraphTypeInstance GraphTypes = "instance"
	GraphTypeSchema   GraphTypes = "schema"
)

type DocumentIntroducer[DocumentT any] BaseIntroducer

func (di *DocumentIntroducer[DocumentT]) OnDatabase(path DatabasePath) *DocumentRequester[DocumentT] {
	return &DocumentRequester[DocumentT]{Client: di.client, path: path}
}

func (di *DocumentIntroducer[DocumentT]) OnBranch(path BranchPath) *DocumentRequester[DocumentT] {
	return &DocumentRequester[DocumentT]{Client: di.client, path: path}
}

func (di *DocumentIntroducer[DocumentT]) OnCommit(path CommitPath) *DocumentRequester[DocumentT] {
	return &DocumentRequester[DocumentT]{Client: di.client, path: path}
}

type DocumentRequester[DocumentT any] BaseRequester

type DocumentListOptions struct {
	CompressIDs bool       `url:"compress_ids,omitempty"`
	Type        string     `url:"type,omitempty"`
	Unfold      bool       `url:"unfold,omitempty"`
	Count       int        `url:"count,omitempty"`
	Skip        int        `url:"skip" default:"0,omitempty"`
	GraphType   GraphTypes `url:"graph_type" default:"instance"`
	Prefixed    bool       `url:"prefixed" default:"true"` // FIXME: figure out for what this param is
}

func (dr *DocumentRequester[DocumentT]) ListAll(ctx context.Context, buf *[]DocumentT, options *DocumentListOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	extraParams := struct {
		AsList    bool `url:"as_list"`
		Minimized bool `url:"minimized"` // Compress json on the server side
	}{true, true}
	sl := dr.Client.C.QueryStruct(options).QueryStruct(&extraParams).Get(dr.path.GetPath("document"))
	if _, err = doRequest(ctx, sl, buf); err != nil {
		return err
	}

	return
}

func (dr *DocumentRequester[DocumentT]) ListAllIterator(ctx context.Context, items chan<- DocumentT, options *DocumentListOptions) error {
	var err error
	if items == nil {
		panic("items channel cannot be nil")
	}
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	extraParams := struct {
		AsList    bool `url:"as_list"`
		Minimized bool `url:"minimized"` // Compress json on the server side
	}{false, true}
	sl := dr.Client.C.QueryStruct(options).QueryStruct(&extraParams).Get(dr.path.GetPath("document"))
	req, err := sl.Request()
	if err != nil {
		return err
	}
	// Making a request using implClient since sl drains and closes resp.Body after the request has been made
	resp, err := dr.Client.implClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	go func() {
		defer func() { _ = resp.Body.Close() }()
		defer func() { _, _ = io.Copy(io.Discard, resp.Body) }()
		defer func() { _ = recover() }() // FIXME: write to log
		decoder := json.NewDecoder(resp.Body)
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

	return nil
}

type DocumentGetOptions struct {
	CompressIDs bool       `url:"compress_ids,omitempty"`
	Type        string     `url:"type,omitempty"`
	Unfold      bool       `url:"unfold,omitempty"`
	GraphType   GraphTypes `url:"graph_type" default:"instance"`
	Prefixed    bool       `url:"prefixed"`
}

func (dr *DocumentRequester[DocumentT]) Get(ctx context.Context, buf *DocumentT, docID string, options *DocumentGetOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	extraParams := struct {
		Minimized bool   `url:"minimized"`
		ID        string `url:"id"`
	}{true, docID}
	sl := dr.Client.C.QueryStruct(options).QueryStruct(&extraParams).Get(dr.path.GetPath("document"))
	if _, err = doRequest(ctx, sl, buf); err != nil {
		return err
	}

	return
}

type DocumentCreateOptions struct {
	RawJSON         bool       `url:"raw_json"`
	FullReplace     bool       `url:"full_replace"`
	GraphType       GraphTypes `url:"graph_type" default:"instance"`
	Message         string     `url:"message" default:"Default message"`
	Author          string     `url:"author" default:"Default author"`
	LastDataVersion string     // TODO: figure out with this
}

func (dr *DocumentRequester[DocumentT]) Create(ctx context.Context, doc DocumentT, options *DocumentCreateOptions) error {
	var docSlice []DocumentT
	docSlice = append(docSlice, doc)
	if _, err := dr.CreateBulk(ctx, docSlice, options); err != nil {
		return err
	}
	return nil
}

func (dr *DocumentRequester[DocumentT]) CreateBulk(ctx context.Context, docs []DocumentT, options *DocumentCreateOptions) (insertedIDs []string, err error) {
	if options, err = prepareOptions(options); err != nil {
		return nil, err
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := dr.Client.C.QueryStruct(options).BodyJSON(docs).Post(dr.path.GetPath("document"))
	// FIXME: check actual response (insertedIDs)
	if _, err = doRequest(ctx, sl, &insertedIDs); err != nil {
		return nil, err
	}
	return
}

type DocumentUpdateOptions struct {
	RawJSON         bool       `url:"raw_json"`
	Create          bool       `url:"create"`
	GraphType       GraphTypes `url:"graph_type" default:"instance"` //  FIXME: check all params in options everywhere if they have to enums
	Message         string     `url:"message" default:"Default message"`
	Author          string     `url:"author" default:"Default author"`
	LastDataVersion string
}

func (dr *DocumentRequester[DocumentT]) UpdateBulk(ctx context.Context, docs []DocumentT, options *DocumentUpdateOptions) (updatedIDs []string, err error) {
	if options, err = prepareOptions(options); err != nil {
		return nil, err
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := dr.Client.C.QueryStruct(options).BodyJSON(docs).Put(dr.path.GetPath("document"))
	if options.LastDataVersion != "" {
		sl.Set("TerminusDB-Data-Version", options.LastDataVersion)
	}
	// FIXME: check actual response (updatedIDs)
	if _, err = doRequest(ctx, sl, &updatedIDs); err != nil {
		return nil, err
	}
	return
}

func (dr *DocumentRequester[DocumentT]) Update(ctx context.Context, doc DocumentT, options *DocumentUpdateOptions) error {
	var docSlice []DocumentT
	docSlice = append(docSlice, doc)
	if _, err := dr.UpdateBulk(ctx, docSlice, options); err != nil {
		return err
	}
	return nil
}

type DocumentDeleteOptions struct {
	LastDataVersion string
}

func (dr *DocumentRequester[DocumentT]) DeleteBulk(ctx context.Context, docIDs []string, options *DocumentDeleteOptions) (deletedIDs []string, err error) {
	if options, err = prepareOptions(options); err != nil {
		return nil, err
	}
	// TODO: maybe need to implement _convert_document function and use here
	sl := dr.Client.C.QueryStruct(options).BodyJSON(docIDs).Delete(dr.path.GetPath("document"))
	if options.LastDataVersion != "" {
		sl.Set("TerminusDB-Data-Version", options.LastDataVersion)
	}
	// FIXME: check actual response (deletedIDs)
	if _, err = doRequest(ctx, sl, &deletedIDs); err != nil {
		return nil, err
	}

	return deletedIDs, nil
}

func (dr *DocumentRequester[DocumentT]) Delete(ctx context.Context, docID string, options *DocumentDeleteOptions) error {
	var docSlice []string
	docSlice = append(docSlice, docID)
	if _, err := dr.DeleteBulk(ctx, docSlice, options); err != nil {
		return err
	}
	return nil
}
