package rest

import "context"

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

func (dr *DiffRequester) WithContext(ctx context.Context) *DiffRequester {
	r := *dr
	r.ctx = ctx
	return &r
}

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

func (dr *DiffRequester) Diff(buf *Diff, options *DiffOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	sl := dr.Client.C.BodyJSON(options).Post("diff")
	return doRequest(dr.ctx, sl, buf)
}

func (dr *DiffRequester) Patch(before any, diff *Diff, buf any) (response TerminusResponse, err error) {
	body := struct {
		Before any   `json:"before"`
		Patch  *Diff `json:"patch"`
	}{before, diff}
	sl := dr.Client.C.BodyJSON(body).Post("patch")
	return doRequest(dr.ctx, sl, buf)
}

func (dr *DiffRequester) DiffObjs(objBefore, objAfter any, buf *Diff, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(buf, &DiffOptions{
		DiffShortOptions: *options,
		Before:           objBefore,
		After:            objAfter,
	})
}

func (dr *DiffRequester) DiffObjAndDocRevision(docRevision string, obj any, buf *Diff, docID string, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(buf, &DiffOptions{
		DiffShortOptions:  *options,
		After:             obj,
		BeforeDataVersion: docRevision,
		DocumentID:        docID,
	})
}

func (dr *DiffRequester) DiffDocRevisions(revisionBefore, revisionAfter string, buf *Diff, docID string, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(buf, &DiffOptions{
		DiffShortOptions:  *options,
		BeforeDataVersion: revisionBefore,
		AfterDataVersion:  revisionAfter,
		DocumentID:        docID,
	})
}

func (dr *DiffRequester) DiffAllDocsRevisions(revisionBefore, revisionAfter string, buf *Diff, options *DiffShortOptions) (response TerminusResponse, err error) {
	if options, err = prepareOptions(options); err != nil {
		return
	}
	return dr.Diff(buf, &DiffOptions{
		DiffShortOptions:  *options,
		BeforeDataVersion: revisionBefore,
		AfterDataVersion:  revisionAfter,
	})
}
