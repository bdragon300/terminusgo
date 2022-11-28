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
	Dictionary DictionaryTemplate `json:"dictionary"` // TODO: field does not used anywhere
	List       []Value            `json:"list"`
	Node       string             `json:"node"`
	Variable   string             `json:"variable"`
	Data       any                `json:"data" terminusgo:"class=xsd:anySimpleType"`
}

func (v *Value) FromVariableName(value string) {
	v.Variable = value
}

func (v *Value) FromString(value string, forceLiteral bool) {
	if forceLiteral {
		newVal := &Literal{}
		newVal.FromAnyValue(value)
		v.Data = *newVal
	} else {
		v.Node = value
	}
}

func (v *Value) FromAnyValue(value any) {
	// TODO: nested value parse (map, list, etc)
}

type NodeValue struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Node     string `json:"node"`
	Variable string `json:"variable"`
}

func (v *NodeValue) FromVariableName(value string) {
	v.Variable = value
}

func (v *NodeValue) FromString(value string, _ bool) {
	v.Node = value
}

type DataValue struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	List     []DataValue `json:"list"`
	Data     any         `json:"data" terminusgo:"class=xsd:anySimpleType"`
	Variable string      `json:"variable"`
}

func (v *DataValue) FromVariableName(value string) {
	v.Variable = value
}

func (v *DataValue) FromString(value string, forceLiteral bool) {
	if forceLiteral {
		newVal := &Literal{}
		newVal.FromAnyValue(value)
		v.Data = *newVal
	} else {
		v.Data = value
	}
}

func (v *DataValue) FromAnyValue(value any) {
	newVal := &Literal{}
	newVal.FromAnyValue(value)
	v.Data = *newVal
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
	URL  string `json:"url"`
}

type FileOptions interface {
	FileFormatType() FormatType
}

type FormatType string

const (
	FormatTypeCSV    FormatType = "csv"
	FormatTypeTurtle FormatType = "turtle"
	FormatTypePanda  FormatType = "panda"
)

type CSVCase string

const (
	CSVCasePreserve CSVCase = "preserve"
	CSVCaseUp       CSVCase = "up"
	CSVCaseDown     CSVCase = "down"
)

// FileOptionsCSV is options list for reading CSV file by TerminusDB
// For more info see https://www.swi-prolog.org/pldoc/man?predicate=csv//2
type FileOptionsCSV struct {
	Separator     string  `json:"separator,omitempty"`
	IgnoreQuoutes bool    `json:"ignore_quoutes,omitempty"`
	Strip         bool    `json:"strip,omitempty"`
	SkipHeader    bool    `json:"skip_header"`
	Convert       bool    `json:"convert"`
	Case          CSVCase `json:"case,omitempty"`
	Functor       string  `json:"functor,omitempty"`
	Arity         uint    `json:"arity,omitempty"`
	MatchArity    bool    `json:"match_arity"`
}

func (i FileOptionsCSV) FileFormatType() FormatType {
	return FormatTypeCSV
}

type TurtleFormat string

const (
	TurtleFormatAuto   TurtleFormat = "auto"
	TurtleFormatTurtle TurtleFormat = "turtle"
	TurtleFormatTrig   TurtleFormat = "trig"
)

type TurtleResources string

const (
	TurtleResourcesURI TurtleResources = "uri"
	TurtleResourcesIRI TurtleResources = "iri"
)

type TurtleOnError string

const (
	TurtleOnErrorWarning TurtleOnError = "warning"
	TurtleOnErrorError   TurtleOnError = "error"
)

// FileOptionsTurtle is options list for reading CSV file by TerminusDB
// For more info see https://www.swi-prolog.org/pldoc/man?predicate=rdf_read_turtle/3
type FileOptionsTurtle struct {
	BaseURL    string          `json:"base_url,omitempty"`
	AnonPrefix string          `json:"anon_prefix,omitempty"`
	Format     TurtleFormat    `json:"format,omitempty"`
	Resources  TurtleResources `json:"resources,omitempty"`
	OnError    TurtleOnError   `json:"on_error,omitempty"`
}

func (i FileOptionsTurtle) FileFormatType() FormatType {
	return FormatTypeTurtle
}

type OrderDirection string

const (
	OrderAscending  OrderDirection = "asc"
	OrderDescending OrderDirection = "desc"
)

type OrderTemplate struct {
	*schema.SubDocumentModel
	Order    OrderDirection `json:"order"`
	Variable string         `json:"variable"`
}

type Literal struct {
	schema.RawModel
	Value any `json:"@value"`
}

func (s *Literal) FromAnyValue(value any) {
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
