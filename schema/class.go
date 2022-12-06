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
	RawModel         struct {
		Type string `mapstructure:"@type" json:"@type" terminusgo:"name=@type"` // TODO: implement
	}
	// TODO: implement Documentation
)

type FieldSerializeMapCallback func(field reflect.StructField, name string, typeName string, value any) any

func SerializeObject(buf map[string]any, object any, mapCallback FieldSerializeMapCallback) error {
	val := reflect.Indirect(reflect.ValueOf(object))
	if !val.IsValid() {
		return errors.New("object is nil")
	}
	typeName := val.Type().Name()
	return serializeObject(buf, val, typeName, mapCallback)
}

func serializeObject(buf map[string]any, object reflect.Value, typeName string, mapCallback FieldSerializeMapCallback) error {
	objectTyp := object.Type()

	for i := 0; i < object.NumField(); i++ {
		fld := objectTyp.Field(i)
		fldVal := object.Field(i)
		fldName := fld.Name
		fldTypeName := typeName // `object`'s "@type" value
		if !fld.IsExported() {
			continue
		}

		// Extract a value from interface
		if fldVal.Kind() == reflect.Interface {
			fldVal = fldVal.Elem()
			if !fldVal.IsValid() {
				continue // Interface is contained nil
			}
		}
		fldVal = reflect.Indirect(fldVal)
		if !fldVal.IsValid() {
			continue // Skip nil pointers
		}
		fldRealVal := fldVal.Interface()

		if !fld.Anonymous {
			if n, s, ok := getFieldSchema(fld); ok {
				fldName = n
				fldTypeName = s.Class
			} else {
				continue
			}
			if mapCallback != nil {
				fldRealVal = mapCallback(fld, fldName, fldTypeName, fldRealVal)
				fldVal = reflect.ValueOf(fldRealVal)
			}
		}

		switch fldVal.Kind() {
		case reflect.Struct:
			m := make(map[string]any)
			err := serializeObject(m, fldVal, fldTypeName, mapCallback)
			if err != nil {
				return err
			}
			if fld.Anonymous { // Extend res map
				for k, v := range m {
					if k != "@type" {
						buf[k] = v
					}
				}
			} else { // Set res in a field
				buf[fldName] = m
			}
		case reflect.Map:
			m := make(map[string]any)
			iter := fldVal.MapRange()
			for iter.Next() {
				m[iter.Key().String()] = iter.Value()
			}
			buf[fldName] = m
		default:
			buf[fldName] = fldRealVal
		}
	}
	// Set @type only if no such field in struct
	if _, ok := buf["@type"]; !ok {
		buf["@type"] = typeName
	}

	return nil
}

func (c *Class) FromValue(obj any) {
	mdlType := reflect.ValueOf(obj).Type()
	if mdlType.Kind() != reflect.Struct {
		panic("obj must be a struct")
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

func (c *Class) FromRaw(m RawSchemaItem) error {
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

func (c *Class) ToRaw(buf RawSchemaItem) error {
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
	buf := make(RawSchemaItem, 7+len(c.Fields))
	if err := c.ToRaw(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (c *Class) UnmarshalJSON(bytes []byte) error {
	buf := make(RawSchemaItem)
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
