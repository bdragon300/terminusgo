package schema

import (
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
)

type Enum struct {
	ID    string   `mapstructure:"@id"`
	Value []string `mapstructure:"@value"`
	// TODO: implement documentation
}

func (e *Enum) FromValue(name string, enumValues []string) {
	e.ID = name
	e.Value = enumValues
}

func (e *Enum) Type() ItemType {
	return EnumSchemaItem
}

func (e *Enum) Name() string {
	return e.ID
}

func (e *Enum) Deserialize(m RawSchemaItem) error {
	if !hasType(m, EnumSchemaItem) {
		return errors.New("item is not a Elass")
	}
	if err := mapstructure.Decode(m, e); err != nil {
		return err
	}
	return nil
}

func (e *Enum) Serialize(buf RawSchemaItem) error {
	if err := mapstructure.Decode(e, &buf); err != nil {
		return err
	}
	buf["@type"] = EnumSchemaItem
	return nil
}

func (e *Enum) MarshalJSON() ([]byte, error) {
	buf := make(RawSchemaItem, 2)
	if err := e.Serialize(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (e *Enum) UnmarshalJSON(bytes []byte) error {
	buf := make(RawSchemaItem)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		return err
	}
	return e.Deserialize(buf)
}
