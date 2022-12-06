package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type RawConverter interface {
	FromRaw(RawSchemaItem) error
	ToRaw(RawSchemaItem) error // FIXME: maybe make parameter as pointer?
}

type ItemType string

const (
	ClassSchemaItem       ItemType = "Class"
	EnumSchemaItem        ItemType = "Enum"
	TaggedUnionSchemaItem ItemType = "TaggedUnion" // TODO: implement (+ @oneOf)
	UnitSchemaItem        ItemType = "Unit"        // TODO: implement
	ContextSchemaItem     ItemType = "context"     // Ad-hoc type, not used in real schema
)

type RawSchemaItem map[string]any

func (rsi RawSchemaItem) ToSchemaItem(schemaItemBuf RawConverter) error {
	return schemaItemBuf.FromRaw(rsi)
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

type Schema struct {
	Context     Context
	SchemaItems []RawConverter
}

func (s *Schema) FromRawSchema(items []RawSchemaItem) error {
	if len(items) == 0 {
		return errors.New("empty schema")
	}
	for _, item := range items {
		schemaItem, err := produceSchemaItem(item)
		if err != nil {
			return fmt.Errorf("unable to convert raw schema to schema item object: %w", err)
		}
		s.SchemaItems = append(s.SchemaItems, schemaItem)
	}
	return nil
}

func (s *Schema) ToRawSchema(buf []RawSchemaItem) error {
	buf2 := make(RawSchemaItem)
	if err := s.Context.ToRaw(buf2); err != nil {
		return fmt.Errorf("unable to convert context object to raw schema: %w", err)
	}
	buf = append(buf[:0], buf2)

	for ind, v := range s.SchemaItems {
		buf2 = make(RawSchemaItem)
		if err := v.ToRaw(buf2); err != nil {
			return fmt.Errorf("unable to convert schema item object to raw schema at index %d: %w", ind, err)
		}
		buf = append(buf, buf2)
	}
	return nil
}

type IdentityKeeper interface {
	Type() ItemType
	Name() string
}

func (s *Schema) FindModel(model any) int {
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	targetID := modelValue.Type().Name()
	return s.findNameType(targetID, ClassSchemaItem)
}

func (s *Schema) FindEnum(name string) int {
	return s.findNameType(name, EnumSchemaItem)
}

func (s *Schema) findNameType(name string, typ ItemType) int {
	for ind, item := range s.SchemaItems {
		obj := item.(IdentityKeeper)
		if obj.Type() == typ && obj.Name() == name {
			return ind
		}
	}
	return -1
}

func (s *Schema) Validate() error {
	if s.Context.Schema == "" {
		return errors.New("empty context object in schema")
	}
	names := make(map[string]struct{})
	for _, item := range s.SchemaItems {
		obj := item.(IdentityKeeper)
		if _, ok := names[obj.Name()]; ok {
			return fmt.Errorf("duplicate schema item with name %s", obj.Name())
		}
		names[obj.Name()] = struct{}{}
	}
	return nil
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	buf := make([]RawSchemaItem, 0, len(s.SchemaItems)+1)
	if err := s.ToRawSchema(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (s *Schema) UnmarshalJSON(bytes []byte) error {
	buf := make([]RawSchemaItem, 0)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		return err
	}
	return s.FromRawSchema(buf)
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

func produceSchemaItem(m RawSchemaItem) (RawConverter, error) {
	var res RawConverter
	factories := map[ItemType]func() RawConverter{
		ClassSchemaItem: func() RawConverter { return &Class{} },
		EnumSchemaItem:  func() RawConverter { return &Enum{} },
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
	err := res.FromRaw(m)
	return res, err
}
