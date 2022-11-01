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
}

type (
	AbstractModel    struct{} // TODO: implement
	SubDocumentModel struct{}
)

func GenerateClass(model any) (schema Class) {
	mdlType := reflect.ValueOf(model).Type()
	if mdlType.Kind() != reflect.Struct {
		panic("model must be a struct")
	}
	ps, gps, fields := analyzeModel(mdlType)

	schema.Fields = fields
	schema.ID = mdlType.Name()
	schema.Key = ClassKey{Type: ClassKeyRandom}

	abTyp := reflect.TypeOf(AbstractModel{})
	sdTyp := reflect.TypeOf(SubDocumentModel{})
	for _, typ := range ps {
		schema.Abstract = schema.Abstract || typ == abTyp
		schema.SubDocument = schema.SubDocument || typ == sdTyp
		schema.Inherits = append(schema.Inherits, typ.Name())
	}
	for _, typ := range gps {
		schema.SubDocument = schema.SubDocument || typ == sdTyp
	}

	return
}

func (c *Class) Type() ItemType {
	return ClassSchemaItem
}

func (c *Class) Name() string {
	return c.ID
}

func (c *Class) Validate() error {
	// TODO: call go-validate
	if c.Documentation != nil {
		defaultLang := false
		for _, v := range c.Documentation {
			if v.Language == "" {
				if defaultLang {
					return errors.New("class documentation has several entries with default language (no language specified)")
				}
				defaultLang = true
			}
		}
	}
	return nil
}

func (c *Class) FromRaw(m RawItem) error {
	if !hasType(m, ClassSchemaItem) {
		return errors.New("raw schema has not class type")
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

func (c *Class) ToRaw(buf RawItem) error {
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
	for k, v := range c.Fields {
		if v.Type == "" {
			buf[k] = v.Class
		} else {
			buf[k] = v
		}
	}
	return nil
}

func (c *Class) MarshalJSON() ([]byte, error) {
	buf := make(RawItem, 7+len(c.Fields))
	if err := c.ToRaw(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (c *Class) UnmarshalJSON(bytes []byte) error {
	buf := make(RawItem)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		return err
	}
	return c.FromRaw(buf)
}

type ClassDocumentationType struct {
	Language   string                                      `json:"@language,omitempty"`
	Label      string                                      `json:"@label"`
	Comment    string                                      `json:"@comment"`
	Properties map[string]ClassDocumentationPropertiesItem `json:"@properties,omitempty" validate:"required_without=Values"`
	Values     map[string]string                           `json:"@values,omitempty" validate:"required_without=Properties"`
}

type ClassDocumentationPropertiesItem struct {
	Label   string `json:"@label"`
	Comment string `json:"@comment"`
}
