package schema

type ContextDocumentation struct {
	Title       string `json:"@title"`
	Description string `json:"@description"`
	Authors     string `json:"@authors"`
}

type Context struct {
	Schema        string                `json:"@schema" validate:"url,required"`
	Base          string                `json:"@base" validate:"url,required"`
	Documentation *ContextDocumentation `json:"@documentation,omitempty"`
}
