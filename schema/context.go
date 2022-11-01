package schema

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

type ContextDocumentation struct {
	Title       string   `json:"@title"`
	Description string   `json:"@description"`
	Authors     []string `json:"@authors"`
}

type Context struct {
	Schema        string                `json:"@schema" mapstructure:"@schema" validate:"url,required"`
	Base          string                `json:"@base" mapstructure:"@base" validate:"url,required"`
	Documentation *ContextDocumentation `json:"@documentation,omitempty" mapstructure:"@documentation,omitempty"`
}

func (c *Context) Type() ItemType {
	return EnumSchemaItem
}

func (c *Context) FromRaw(m RawItem) error {
	if !hasType(m, ContextSchemaItem) {
		return errors.New("raw schema has not context type")
	}
	if err := mapstructure.Decode(m, c); err != nil { // TODO: check if mapstructure resets all fields even if they not present in map. @documentation here, for instance
		return err
	}
	return nil
}

func (c *Context) ToRaw(buf RawItem) error {
	if err := mapstructure.Decode(c, &buf); err != nil {
		return err
	}
	return nil
}
