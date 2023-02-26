package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bdragon300/terminusgo/srverror"
)

type TripleIntroducer BaseIntroducer

func (di *TripleIntroducer) OnDatabase(path DatabasePath) *TripleRequester {
	return &TripleRequester{Client: di.client, path: path}
}

func (di *TripleIntroducer) OnBranch(path BranchPath) *TripleRequester {
	return &TripleRequester{Client: di.client, path: path}
}

func (di *TripleIntroducer) OnRepo(path RepoPath) *TripleRequester {
	return &TripleRequester{Client: di.client, path: path}
}

func (di *TripleIntroducer) OnCommit(path CommitPath) *TripleRequester {
	return &TripleRequester{Client: di.client, path: path}
}

type TripleRequester BaseRequester

func (tr *TripleRequester) WithContext(ctx context.Context) *TripleRequester {
	r := *tr
	r.ctx = ctx
	return &r
}

type TripleDumpOptions struct {
	GraphType GraphTypes `url:"-" default:"instance"`
	Format    string     `url:"format" default:"turtle"`
}

func (tr *TripleRequester) DumpAsStream(w io.Writer, options *TripleDumpOptions) (writtenBytes int64, response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	acceptHeader := "application/json"
	if options.Format == "turtle" {
		acceptHeader = "text/turtle"
	}
	sl := tr.Client.C.QueryStruct(options).Add("Accept", acceptHeader).Get(tr.getURL("pack", options.GraphType))

	var httpReq *http.Request
	if httpReq, err = sl.Request(); err != nil {
		return
	}
	if tr.ctx != nil {
		httpReq = httpReq.WithContext(tr.ctx)
	}

	var httpResp *http.Response
	if httpResp, err = tr.Client.implClient.Do(httpReq); err != nil {
		return
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode >= 300 {
		res := srverror.TerminusError{Response: httpResp}
		if err = json.NewDecoder(httpResp.Body).Decode(&res); err != nil {
			return
		}
		return writtenBytes, res, res
	}
	response = &srverror.TerminusOkResponse{Response: httpResp}
	writtenBytes, err = io.Copy(w, httpResp.Body)
	return
}

func (tr *TripleRequester) DumpAsString(buf *string, options *TripleDumpOptions) (response TerminusResponse, err error) {
	if buf == nil {
		panic("buf is empty")
	}
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := tr.Client.C.QueryStruct(options).Add("Accept", "application/json").Get(tr.getURL("pack", options.GraphType))
	return doRequest(tr.ctx, sl, buf)
}

type TripleUpdateOptions struct {
	GraphType GraphTypes `json:"-" default:"instance"`
	Author    string     `json:"author" default:"defaultAuthor"`
	Message   string     `json:"message" default:"Default commit message"`
}

func (tr *TripleRequester) Update(data *string, options *TripleUpdateOptions) (response TerminusResponse, err error) {
	if data == nil {
		panic("data is empty")
	}
	if options, err = prepareOptions(options); err != nil {
		return
	}
	commitInfo := struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}{Author: options.Author, Message: options.Message}
	body := struct {
		CommitInfo any     `json:"commit_info"`
		Turtle     *string `json:"turtle"`
	}{CommitInfo: commitInfo, Turtle: data}
	sl := tr.Client.C.BodyJSON(body).Post(tr.getURL("pack", options.GraphType))
	return doRequest(tr.ctx, sl, nil)
}

type TripleInsertOptions struct {
	GraphType GraphTypes `json:"-" default:"instance"`
	Author    string     `json:"author" default:"defaultAuthor"`
	Message   string     `json:"message" default:"Default commit message"`
}

func (tr *TripleRequester) Insert(data *string, options *TripleInsertOptions) (response TerminusResponse, err error) {
	if data == nil {
		panic("data is empty")
	}
	if options, err = prepareOptions(options); err != nil {
		return
	}
	commitInfo := struct {
		Author  string `json:"author"`
		Message string `json:"message"`
	}{Author: options.Author, Message: options.Message}
	body := struct {
		CommitInfo any     `json:"commit_info"`
		Turtle     *string `json:"turtle"`
	}{CommitInfo: commitInfo, Turtle: data}
	sl := tr.Client.C.BodyJSON(body).Post(tr.getURL("pack", options.GraphType))
	return doRequest(tr.ctx, sl, nil)
}

func (tr *TripleRequester) getURL(action string, graphType GraphTypes) string {
	return fmt.Sprintf("%s/%s/%s", action, tr.path.String(), graphType)
}
