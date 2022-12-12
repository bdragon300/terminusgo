package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type TaggedUnion struct {
	ID            string                   `mapstructure:"@id"`
	Key           ClassKey                 `mapstructure:"@key"`
	Documentation []ClassDocumentationType `mapstructure:"@documentation,omitempty"` // TODO: implement -- how to describe in object
	Base          string                   `mapstructure:"@base"`
	Fields        map[string]Field         `mapstructure:"-"`
}

type TaggedUnionModel struct{}

func (t *TaggedUnion) FromObject(obj any) {
	// TODO: refactor, move to a separate function (with those one in Class's FromObject)
	mdlVal := reflect.ValueOf(obj).Elem()
	if !mdlVal.IsValid() {
		panic("obj is nil")
	}
	mdlType := mdlVal.Type()
	if mdlType.Kind() != reflect.Struct {
		panic("obj must be a struct object")
	}
	ps, gps, fields := analyzeModel(mdlType)
	if !isParent(reflect.TypeOf(TaggedUnionModel{}), ps, gps) {
		panic(fmt.Sprintf("Type %T is not a TaggedUnion", mdlVal))
	}

	t.Fields = fields
	t.ID = mdlType.Name()
	t.Key = ClassKey{Type: ClassKeyRandom}
}

func (t *TaggedUnion) Type() ItemType {
	return TaggedUnionSchemaItem
}

func (t *TaggedUnion) Name() string {
	return t.ID
}

func (t *TaggedUnion) Deserialize(m RawSchemaItem) error {
	if !hasType(m, TaggedUnionSchemaItem) {
		return errors.New("item is not a TaggedUnion")
	}
	if err := mapstructure.Decode(m, t); err != nil {
		return err
	}
	t.Fields = make(map[string]Field)
	for k, v := range m {
		fields, err := parseRawFieldSchema(k, v)
		if err != nil {
			return err
		}
		for k1, v1 := range fields {
			if _, ok := t.Fields[k]; ok {
				return fmt.Errorf("field %s is duplicated", k)
			}
			t.Fields[k1] = v1
		}
	}
	return nil
}

func (t *TaggedUnion) Serialize(buf RawSchemaItem) error {
	if err := mapstructure.Decode(t, &buf); err != nil {
		return err
	}
	buf["@type"] = TaggedUnionSchemaItem
	structFields, oneOfFields := groupStructFields(t.Fields)
	if len(oneOfFields) > 0 {
		buf["@oneOf"] = oneOfFields
	}
	for k, v := range structFields {
		buf[k] = v
	}

	return nil
}

func (t *TaggedUnion) MarshalJSON() ([]byte, error) {
	buf := make(RawSchemaItem)
	if err := t.Serialize(buf); err != nil {
		return nil, err
	}
	return json.Marshal(buf)
}

func (t *TaggedUnion) UnmarshalJSON(bytes []byte) error {
	buf := make(RawSchemaItem)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		return err
	}
	return t.Deserialize(buf)
}
