package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type ClassKeyType string

const (
	ClassKeyRandom    ClassKeyType = "Random"
	ClassKeyLexical   ClassKeyType = "Lexical"
	ClassKeyHash      ClassKeyType = "Hash"
	ClassKeyValueHash ClassKeyType = "ValueHash"
)

type ClassKey struct {
	Type   ClassKeyType `mapstructure:"@type"`
	Fields []string     `mapstructure:"@fields,omitempty"`
}

type Class struct {
	ID            string                   `mapstructure:"@id"`
	Key           ClassKey                 `mapstructure:"@key"`
	Documentation []ClassDocumentationType `mapstructure:"@documentation,omitempty"`
	Base          string                   `mapstructure:"@base"`
	Inherits      []string                 `mapstructure:"-"`
	Fields        map[string]Field         `mapstructure:"-"`
	SubDocument   bool                     `mapstructure:"-"`
	Abstract      bool                     `mapstructure:"-"`
	// TODO: implement @metadata
}

type (
	AbstractModel    struct{} // TODO: implement
	SubDocumentModel struct{}
	// TODO: implement Documentation
)

func (c *Class) FromObject(obj any) {
	mdlVal := reflect.ValueOf(obj).Elem()
	if !mdlVal.IsValid() {
		panic("obj is nil")
	}
	mdlType := mdlVal.Type()
	if mdlType.Kind() != reflect.Struct {
		panic("obj must be a struct object")
	}
	ps, gps, fields := analyzeModel(mdlType)

	c.Fields = fields
	c.ID = mdlType.Name()
	c.Key = ClassKey{Type: ClassKeyRandom}

	abTyp := reflect.TypeOf(AbstractModel{})
	sdTyp := reflect.TypeOf(SubDocumentModel{})
	for _, typ := range ps {
		c.Abstract = c.Abstract || typ == abTyp
		c.SubDocument = c.SubDocument || typ == sdTyp
		c.Inherits = append(c.Inherits, typ.Name())
	}
	for _, typ := range gps {
		c.SubDocument = c.SubDocument || typ == sdTyp
	}
}

func (c *Class) Type() ItemType {
	return ClassSchemaItem
}

func (c *Class) Name() string {
	return c.ID
}

func (c *Class) Deserialize(m RawSchemaItem) error {
	if !hasType(m, ClassSchemaItem) {
		return errors.New("item is not a Class")
	}
	if err := mapstructure.Decode(m, c); err != nil {
		return err
	}
	if _, ok := m["@subdocument"]; ok {
		c.SubDocument = true
	}
	if _, ok := m["@abstract"]; ok {
		c.Abstract = true
	}
	if v, ok := m["@inherits"]; ok {
		switch vt := v.(type) {
		case string:
			c.Inherits = []string{vt}
		case []string:
			c.Inherits = vt
		default:
			return fmt.Errorf("unknown value type in @inherits field: %T", v)
		}
	}
	c.Fields = make(map[string]Field)
	for k, v := range m {
		if strings.HasPrefix(k, "@") {
			continue
		}
		switch vt := v.(type) {
		case string:
			c.Fields[k] = Field{Class: vt}
		default:
			fschema := Field{}
			if err := mapstructure.Decode(vt, &fschema); err != nil {
				return err
			}
			c.Fields[k] = fschema
		}
	}

	return nil
}

func (c *Class) Serialize(buf RawSchemaItem) error {
	if err := mapstructure.Decode(c, &buf); err != nil {
		return err
	}
	buf["@type"] = ClassSchemaItem
	if len(c.Inherits) == 1 {
		buf["@inherits"] = c.Inherits[0]
	} else if len(c.Inherits) > 1 {
		buf["@inherits"] = c.Inherits
	}
	if c.SubDocument {
		buf["@subdocument"] = &[]any{}
	}
	if c.Abstract {
		buf["@abstract"] = &[]any{}
	}

	structFields, oneOfFields := groupStructFields(c.Fields)
	if len(oneOfFields) > 0 {
		buf["@oneOf"] = oneOfFields
	}
	for k, v := range structFields {
		buf[k] = v
	}

	return nil
}

func (c *Class) MarshalJSON() ([]byte, error) {
	buf := make(RawSchemaItem)
	if err := c.Serialize(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (c *Class) UnmarshalJSON(bytes []byte) error {
	buf := make(RawSchemaItem)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		return err
	}
	return c.Deserialize(buf)
}

