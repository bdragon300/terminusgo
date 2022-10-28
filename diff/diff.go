package diff

import (
	"context"
	"strings"

	"github.com/bdragon300/terminusgo/objects"
	"github.com/bdragon300/terminusgo/srverror"
	"github.com/mitchellh/mapstructure"
)

type Operation string

const (
	SwapValue   Operation = "SwapValue"
	ForceValue  Operation = "ForceValue"
	CopyList    Operation = "CopyList"
	SwapList    Operation = "SwapList"
	PatchList   Operation = "PatchList"
	KeepList    Operation = "KeepList"
	ModifyTable Operation = "ModifyTable"
)

type FieldDiff struct {
	FieldName string    `json:"field_name"`
	Op        Operation `json:"@op" mapstructure:"@op"`
	Before    string    `json:"@before"  mapstructure:"@before"`
	After     string    `json:"@after"  mapstructure:"@after"`
	Patch     string    `json:"@patch"  mapstructure:"@patch"`
	Rest      string    `json:"@rest"  mapstructure:"@rest"`
	To        string    `json:"@to"  mapstructure:"@to"`
	// TODO: also figure out how to handle ModifyTable
	// TODO: also which fields should be here or in another kind of diff?
}

type Diff map[string]any

func (d Diff) ExtractFields() ([]FieldDiff, error) {
	var res []FieldDiff
	for k, v := range d {
		if strings.HasPrefix(k, "@") {
			continue
		}
		var item FieldDiff
		if err := mapstructure.Decode(v, &item); err != nil {
			return res, err
		}
		item.FieldName = k
		res = append(res, item)
	}

	return res, nil
}

type GetDiffOptions struct {
	Keep      map[string]any `json:"keep"` // FIXME: figure out correct type
	CopyValue bool           `json:"copy_value"`
}

// TODO: refactor the code below
// FIXME: figure out if GetDiff* methods must be moved to Branch
// https://terminusdb.com/docs/guides/reference-guides/json-diff-and-patch#diff
func GetDiffObjects(ctx context.Context, client *objects.Client, buf *Diff, before, after any, options *GetDiffOptions) error {
	if options == nil {
		options = new(GetDiffOptions)
	}
	return RequestDiff(ctx, client, buf, Request{
		Before:    before,
		After:     after,
		Keep:      options.Keep,
		CopyValue: options.CopyValue,
	})
}

func GetDiffObjectAndDocumentRevision(ctx context.Context, client *objects.Client, buf *Diff, beforeVersion string, after any, documentID string, options *GetDiffOptions) error {
	if options == nil {
		options = new(GetDiffOptions)
	}
	return RequestDiff(ctx, client, buf, Request{
		After:             after,
		BeforeDataVersion: beforeVersion,
		DocumentID:        documentID,
		Keep:              options.Keep,
		CopyValue:         options.CopyValue,
	})
}

func GetDiffDocumentRevisions(ctx context.Context, client *objects.Client, buf *Diff, beforeVersion, afterVersion string, documentID string, options *GetDiffOptions) error {
	if options == nil {
		options = new(GetDiffOptions)
	}
	return RequestDiff(ctx, client, buf, Request{
		BeforeDataVersion: beforeVersion,
		AfterDataVersion:  afterVersion,
		DocumentID:        documentID,
		Keep:              options.Keep,
		CopyValue:         options.CopyValue,
	})
}

func GetDiffAllDocumentsRevisions(ctx context.Context, client *objects.Client, buf *Diff, beforeVersion, afterVersion string, options *GetDiffOptions) error {
	if options == nil {
		options = new(GetDiffOptions)
	}
	return RequestDiff(ctx, client, buf, Request{
		BeforeDataVersion: beforeVersion,
		AfterDataVersion:  afterVersion,
		Keep:              options.Keep,
		CopyValue:         options.CopyValue,
	})
}

type Request struct {
	Before            any            `json:"before,omitempty"`
	After             any            `json:"after,omitempty"`
	BeforeDataVersion string         `json:"before_data_version,omitempty"`
	AfterDataVersion  string         `json:"after_data_version,omitempty"`
	DocumentID        string         `json:"document_id,omitempty"`
	Keep              map[string]any `json:"keep,omitempty"` // FIXME: figure out correct type
	CopyValue         bool           `json:"copy_value,omitempty"`
}

func RequestDiff(ctx context.Context, client *objects.Client, buf *Diff, body Request) error {
	// FIXME: figure out response schema
	sl := client.C.BodyJSON(body).Post("diff") // FIXME: maybe this can be a branch path?
	errTerminus := new(srverror.TerminusError)
	req, err := sl.Request()
	req = req.WithContext(ctx)
	if err != nil {
		return err
	}
	resp, err := sl.Do(req, buf, errTerminus)
	if err != nil {
		return err
	}
	if resp.StatusCode > 300 {
		errTerminus.HTTPCode = resp.StatusCode
		return errTerminus
	}

	return nil
}

func GetPatchedObject(ctx context.Context, client *objects.Client, buf any, before any, diff *Diff) error {
	body := struct {
		Before any   `json:"before"`
		Patch  *Diff `json:"patch"`
	}{before, diff} // FIXME: check if this is a Document actually
	sl := client.C.BodyJSON(body).Post("patch") // FIXME: maybe this can be a branch path?
	errTerminus := new(srverror.TerminusError)
	req, err := sl.Request()
	if err != nil {
		return err
	}
	resp, err := sl.Do(req.WithContext(ctx), buf, errTerminus)
	if err != nil {
		return err
	}
	if resp.StatusCode > 300 {
		errTerminus.HTTPCode = resp.StatusCode
		return errTerminus
	}

	return nil
}
