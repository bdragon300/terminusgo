package schema

import (
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
)

type Enum struct {
	ID    string   `mapstructure:"@id"`
	Value []string `mapstructure:"@value" validate:"required"`
}

func GenerateEnum(name string, enumValues []string) (schema Enum) {
	schema.ID = name
	schema.Value = enumValues
	return
}

func (e *Enum) Type() ItemType {
	return EnumSchemaItem
}

func (e *Enum) Name() string {
	return e.ID
}

func (e *Enum) FromRaw(m RawItem) error {
	if !hasType(m, EnumSchemaItem) {
		return errors.New("raw schema has not enum type")
	}
	if err := mapstructure.Decode(m, e); err != nil {
		return err
	}
	return nil
}

func (e *Enum) ToRaw(buf RawItem) error {
	if err := mapstructure.Decode(e, &buf); err != nil {
		return err
	}
	buf["@type"] = EnumSchemaItem
	return nil
}

func (e *Enum) MarshalJSON() ([]byte, error) {
	buf := make(RawItem, 2)
	if err := e.ToRaw(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (e *Enum) UnmarshalJSON(bytes []byte) error {
	buf := make(RawItem)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		return err
	}
	return e.FromRaw(buf)
}
