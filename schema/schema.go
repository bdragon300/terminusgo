package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Serializable interface {
	Deserialize(RawSchemaItem) error
	Serialize(RawSchemaItem) error // FIXME: maybe make parameter as pointer?
}

type Schema struct {
	Context     Context
	SchemaItems []Serializable
}

func (s *Schema) Deserialize(items []RawSchemaItem) error {
	if len(items) == 0 {
		return errors.New("empty schema")
	}
	for _, item := range items {
		schemaItem, err := deserializeSchemaItem(item)
		if err != nil {
			return fmt.Errorf("unable to convert raw schema to schema item object: %w", err)
		}
		s.SchemaItems = append(s.SchemaItems, schemaItem)
	}
	return nil
}

func (s *Schema) Serialize(buf []RawSchemaItem) error {
	buf2 := make(RawSchemaItem)
	if err := s.Context.Serialize(buf2); err != nil {
		return fmt.Errorf("unable to convert context object to raw schema: %w", err)
	}
	buf = append(buf[:0], buf2)

	for ind, v := range s.SchemaItems {
		buf2 = make(RawSchemaItem)
		if err := v.Serialize(buf2); err != nil {
			return fmt.Errorf("unable to convert schema item object to raw schema at index %d: %w", ind, err)
		}
		buf = append(buf, buf2)
	}
	return nil
}

func (s *Schema) FindModel(model any) int {
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	targetID := modelValue.Type().Name()
	return s.findNameType(targetID, ClassSchemaItem)
}

func (s *Schema) FindEnum(name string) int {
	return s.findNameType(name, EnumSchemaItem)
}

type IdentityKeeper interface {
	Type() ItemType
	Name() string
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
	if err := s.Serialize(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (s *Schema) UnmarshalJSON(bytes []byte) error {
	buf := make([]RawSchemaItem, 0)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		return err
	}
	return s.Deserialize(buf)
}
