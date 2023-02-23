package srverror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const (
	DataVersionHeader         = "Terminusdb-Data-Version"
	RemoteAuthorizationHeader = "Authorization-Remote"
	XTerminusDBApiBaseHeader  = "X-Terminusdb-Api-Base" // TODO: related to resumable_endpoint option in tus prolog client...
)

type TerminusError struct {
	Type            string           `json:"@type"`
	Auth            string           `json:"auth"`
	Scope           string           `json:"scope"`
	APIError        string           `json:"api:error"` // FIXME: string or dict
	APIPath         string           `json:"api:path"`
	APIMessage      string           `json:"api:message"`
	APIStatus       string           `json:"api:status"`
	SystemWitnesses []map[string]any `json:"system:witnesses"`

	Response *http.Response `json:"-"`
}

func (ter TerminusError) String() string {
	return fmt.Sprintf("terminus db server returned HTTP code %d: %s", ter.Response.StatusCode, ter.APIMessage)
}

func (ter TerminusError) IsOK() bool {
	return false
}

func (ter TerminusError) DataVersion() string {
	return ter.Response.Header.Get(DataVersionHeader)
}

func (ter TerminusError) Error() string {
	return ter.String()
}

func (ter TerminusError) Is(e error) bool {
	if v, ok := e.(TerminusError); ok {
		ter.Response = v.Response // Exclude Response pointer from comparing
		return reflect.DeepEqual(ter, v)
	}
	return false
}

type TerminusOkResponse struct {
	Type      string
	APIFields map[string]any
	Rest      map[string]any

	Response *http.Response
}

func (tor *TerminusOkResponse) UnmarshalJSON(b []byte) error {
	m := make(map[string]any)
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	if v, ok := m["@type"]; ok {
		tor.Type = v.(string)
		delete(m, "@type")
	}

	for k, v := range m {
		if !strings.HasPrefix(k, "api:") {
			if tor.Rest == nil {
				tor.Rest = make(map[string]any)
			}
			tor.Rest[k] = v
			delete(m, k)
		}
	}

	tor.APIFields = m
	return nil
}

func (tor TerminusOkResponse) String() string {
	return "successful request"
}

func (tor TerminusOkResponse) IsOK() bool {
	return true
}

func (tor TerminusOkResponse) DataVersion() string {
	return tor.Response.Header.Get(DataVersionHeader)
}

type WOQLResponse struct {
	APIStatus             string         `json:"api:status"`
	APIVariableNames      string         `json:"api:variable_names"`
	Bindings              map[string]any `json:"bindings"` // TODO: implement smth to extract/working with this field (and the one above)
	Inserts               uint           `json:"inserts"`
	Deletes               uint           `json:"deletes"`
	TransactionRetryCount uint           `json:"transaction_retry_count"`
}
