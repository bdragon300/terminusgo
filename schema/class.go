package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type ClassKeyType string

const (
	RandomClassKey    ClassKeyType = "Random"
	LexicalClassKey   ClassKeyType = "Lexical"
	HashClassKey      ClassKeyType = "Hash"
	ValueHashClassKey ClassKeyType = "ValueHash"
)

type ClassKey struct {
	Type   ClassKeyType `json:"@type"`
	Fields []string     `json:"@fields,omitempty"`
}

type Class struct {
	ID            string                   `json:"@id" copier:"@id"`
	Key           ClassKey                 `json:"@key" copier:"@key"`
	Documentation []ClassDocumentationType `json:"@documentation" copier:"@documentation"`
	Base          string                   `json:"@base" copier:"@base"`
	SubDocument   bool                     `json:"-" copier:"-"`
	Abstract      bool                     `json:"-" copier:"-"`
	Inherits      string                   `json:"-" copier:"-"`
	InheritsMany  []string                 `json:"-" copier:"-"` // TODO: make ORM-like type declarations instead of specifying dependencies in fields here
}

func (ct *Class) Validate() error {
	if ct.Documentation != nil {
		defaultLang := false
		for _, v := range ct.Documentation {
			if v.Language == "" {
				if defaultLang {
					return errors.New("class documentation has several entries with default language (no language specified)")
				}
				defaultLang = true
			}
		}
	}
	if ct.Inherits != "" && len(ct.InheritsMany) > 0 {
		return errors.New("both Inherits and InheritsMany fields are set")
	}
	return nil
}

func (ct *Class) CopyFrom(item map[string]any) error {
	if err := copier.Copy(*ct, item); err != nil {
		return err
	}
	ct.SubDocument = item["@subdocument"] != nil
	ct.Abstract = item["@abstract"] != nil
	if v, ok := item["@inherits"]; ok {
		switch v.(type) {
		case string, []string:
			item["@inherits"] = v
		default:
			return fmt.Errorf("unknown value type in @inherits field: %T", v)
		}
	}
	return nil
}

func (ct *Class) MarshalJSON() ([]byte, error) {
	var body struct {
		*Class
		SubDocumentT *[]any   `json:"@subdocument,omitempty"`
		AbstractT    *[]any   `json:"@abstract,omitempty"`
		TypeT        ItemType `json:"@type"`
		InheritsT    any      `json:"@inherits,omitempty"`
	}
	body.Class = ct
	body.TypeT = ItemTypeClass
	if body.SubDocument {
		body.SubDocumentT = &[]any{}
	}
	if body.Abstract {
		body.AbstractT = &[]any{}
	}
	if body.Inherits != "" {
		body.InheritsT = body.Inherits
	} else if len(body.InheritsMany) > 0 {
		body.InheritsT = body.InheritsMany
	}
	return json.Marshal(body)
}

type ClassDocumentationType struct {
	Language   string                        `json:"@language,omitempty"`
	Label      string                        `json:"@label"`
	Comment    string                        `json:"@comment"`
	Properties *ClassDocumentationProperties `json:"@properties,omitempty" validate:"required_without=Values"`
	Values     *ClassDocumentationEnum       `json:"@values,omitempty" validate:"required_without=Properties"`
}

type ClassDocumentationProperties map[string]ClassDocumentationPropertiesItem

type ClassDocumentationPropertiesItem struct {
	Label   string `json:"@label"`
	Comment string `json:"@comment"`
}

type ClassDocumentationEnum map[string]string
