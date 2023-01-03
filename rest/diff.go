package rest

import (
	"context"
)

type DiffOp string

const (
	DiffOpSwapValue   DiffOp = "SwapValue"
	DiffOpForceValue  DiffOp = "ForceValue"
	DiffOpCopyList    DiffOp = "CopyList"
	DiffOpSwapList    DiffOp = "SwapList"
	DiffOpPatchList   DiffOp = "PatchList"
	DiffOpKeepList    DiffOp = "KeepList"
	DiffOpModifyTable DiffOp = "ModifyTable"
)

type Diff map[string]any

type DiffRequester BaseRequester

type DiffShortOptions struct {
	Keep      map[string]bool `json:"keep,omitempty"`
	CopyValue bool            `json:"copy_value,omitempty"`
}

type DiffOptions struct {
	DiffShortOptions
	Before            any    `json:"before,omitempty"`
	After             any    `json:"after,omitempty"`
	BeforeDataVersion string `json:"before_data_version,omitempty"`
	AfterDataVersion  string `json:"after_data_version,omitempty"`
	DocumentID        string `json:"document_id,omitempty"`
}

func (dr *DiffRequester) Diff(ctx context.Context, buf *Diff, options *DiffOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.BodyJSON(options).Post("diff")
	return doRequest(ctx, sl, buf)
}

func (dr *DiffRequester) Patch(ctx context.Context, buf, before any, diff *Diff) (response TerminusResponse, err error) {
	body := struct {
		Before any   `json:"before"`
		Patch  *Diff `json:"patch"`
	}{before, diff}
	sl := dr.Client.C.BodyJSON(body).Post("patch")
	return doRequest(ctx, sl, buf)
}

func (dr *DiffRequester) DiffObjs(ctx context.Context, buf *Diff, objBefore, objAfter any, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(ctx, buf, &DiffOptions{
		DiffShortOptions: *options,
		Before:           objBefore,
		After:            objAfter,
	})
}

func (dr *DiffRequester) DiffObjAndDocRevision(ctx context.Context, buf *Diff, docRevision string, obj any, docID string, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(ctx, buf, &DiffOptions{
		DiffShortOptions:  *options,
		After:             obj,
		BeforeDataVersion: docRevision,
		DocumentID:        docID,
	})
}

func (dr *DiffRequester) DiffDocRevisions(ctx context.Context, buf *Diff, revisionBefore, revisionAfter string, docID string, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(ctx, buf, &DiffOptions{
		DiffShortOptions:  *options,
		BeforeDataVersion: revisionBefore,
		AfterDataVersion:  revisionAfter,
		DocumentID:        docID,
	})
}

func (dr *DiffRequester) DiffAllDocsRevisions(ctx context.Context, buf *Diff, revisionBefore, revisionAfter string, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(ctx, buf, &DiffOptions{
		DiffShortOptions:  *options,
		BeforeDataVersion: revisionBefore,
		AfterDataVersion:  revisionAfter,
	})
}
