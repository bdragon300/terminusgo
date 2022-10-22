package srverror

import "fmt"

type TerminusError struct {
	HTTPCode        int
	Type            string           `json:"@type"`
	Auth            string           `json:"auth"`
	Scope           string           `json:"scope"`
	APIError        string           `json:"api:error"` // FIXME: string or dict
	APIPath         string           `json:"api:path"`
	APIMessage      string           `json:"api:message"`
	APIStatus       string           `json:"api:status"`
	SystemWitnesses []map[string]any `json:"system:witnesses"`
}

func (he TerminusError) Error() string {
	if he.HTTPCode == 0 {
		return "No error"
	}
	return fmt.Sprintf("terminus db server returned HTTP code %d: %s", he.HTTPCode, he.APIMessage)
}
