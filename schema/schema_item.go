package schema

import (
	"errors"
	"fmt"
)

type ItemType string

const (
	ClassSchemaItem       ItemType = "Class"
	EnumSchemaItem        ItemType = "Enum"
	TaggedUnionSchemaItem ItemType = "TaggedUnion"
	UnitSchemaItem        ItemType = "Unit"    // TODO: implement
	ContextSchemaItem     ItemType = "context" // Ad-hoc type, not used in real schema
)

type RawSchemaItem map[string]any

func (rsi RawSchemaItem) ToSchemaItem(schemaItemBuf Serializable) error {
	return schemaItemBuf.Deserialize(rsi)
}

func (rsi RawSchemaItem) Type() (ItemType, error) {
	variants := [5]ItemType{ContextSchemaItem, ClassSchemaItem, EnumSchemaItem, TaggedUnionSchemaItem, UnitSchemaItem}
	for _, t := range variants {
		if hasType(rsi, t) {
			return t, nil
		}
	}
	return "", errors.New("cannot determine schema item type")
}

func hasType(item RawSchemaItem, typ ItemType) bool {
	if val, ok := item["@type"]; ok {
		return val.(string) == string(typ)
	}
	if _, ok := item["@schema"]; ok {
		return typ == ContextSchemaItem
	}
	return false
}

func deserializeSchemaItem(m RawSchemaItem) (Serializable, error) {
	// TODO: implement inline @class definition, see `/schema` endpoint response
	var res Serializable
	factories := map[ItemType]func() Serializable{
		ClassSchemaItem: func() Serializable { return &Class{} },
		EnumSchemaItem:  func() Serializable { return &Enum{} },
	}
	if hasType(m, ContextSchemaItem) {
		res = &Context{}
	} else {
		typ, ok := m["@type"]
		if !ok {
			return nil, fmt.Errorf("empty @type field in schema item")
		}
		factory, ok := factories[ItemType(typ.(string))]
		if !ok {
			return nil, fmt.Errorf("unknown @type field value %v", typ.(string))
		}
		res = factory()
	}
	err := res.Deserialize(m)
	return res, err
}

type ClassDocumentationType struct {
	Language   string                                      `json:"@language,omitempty"`
	Label      string                                      `json:"@label"`
	Comment    string                                      `json:"@comment"`
	Properties map[string]ClassDocumentationPropertiesItem `json:"@properties,omitempty"`
	Values     map[string]string                           `json:"@values,omitempty"`
}

type ClassDocumentationPropertiesItem struct {
	Label   string `json:"@label"`
	Comment string `json:"@comment"`
}
