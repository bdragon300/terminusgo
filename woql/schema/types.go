package schema

import (
	"fmt"
	"reflect"

	"github.com/bdragon300/terminusgo/schema"
)

type DictionaryTemplate struct {
	*schema.SubDocumentModel
	Data []FieldValuePair `json:"data" terminusgo:"type=Set"`
}

type FieldValuePair struct {
	*schema.SubDocumentModel
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Value struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Dictionary DictionaryTemplate `json:"dictionary"`
	List       []Value            `json:"list"`
	Node       string             `json:"node"`
	Variable   string             `json:"variable"`
	Data       any                `json:"data" terminusgo:"class=xsd:anySimpleType"`
}

func (v *Value) FromVariable(value string) {
	v.Variable = value
}

func (v *Value) FromString(value string) {
	v.Node = value
}

type NodeValue struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Node     string `json:"node"`
	Variable string `json:"variable"`
}

func (v *NodeValue) FromVariable(value string) {
	v.Variable = value
}

func (v *NodeValue) FromString(value string) {
	v.Node = value
}

type DataValue struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	List     []DataValue `json:"list"`
	Data     any         `json:"data" terminusgo:"class=xsd:anySimpleType"`
	Variable string      `json:"variable"`
}

func (v *DataValue) FromVariable(value string) {
	v.Variable = value
}

type Indicator struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Name  string `json:"name"`
	Index uint   `json:"index"`
}

type Column struct {
	*schema.SubDocumentModel
	Indicator Indicator `json:"indicator"`
	Variable  string    `json:"variable"`
	Type      *string   `json:"type"`
}

type Source struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Post string `json:"post"`
	URI  string `json:"uri"`
}

type FormatType string

const FormatTypeCSV FormatType = "csv"

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type OrderTemplate struct {
	*schema.SubDocumentModel
	Order    Order  `json:"order"`
	Variable string `json:"variable"`
}

type SimpleValue struct {
	schema.RawModel
	Value any `json:"@value"`
}

func (s *SimpleValue) SetValue(value any) {
	typ := reflect.ValueOf(value).Type()
	if cls, ok := schema.GetSchemaClass(typ); ok {
		s.RawModel = schema.RawModel{Type: cls}
		s.Value = value
		if conv, ok := schema.GetConverter(typ); ok {
			s.Value = conv(s.Value)
		}
		return
	}
	panic(fmt.Sprintf(
		"Cannot determine schema type of value with type %T, "+
			"maybe it's needed to define type (see schema.DefineTypeClass() or schema.DefinePrimitiveTypeClass)?", value,
	))
}
