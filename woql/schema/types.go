package schema

import (
	"fmt"
	"reflect"

	"github.com/bdragon300/terminusgo/schema"
)

type DictionaryTemplate struct {
	*schema.SubDocumentModel
	Data []FieldValuePair `terminusgo:"type=Set"`
}

type FieldValuePair struct {
	*schema.SubDocumentModel
	Field string
	Value Value
}

type Value struct {
	*schema.TaggedUnionModel
	*schema.SubDocumentModel
	Dictionary DictionaryTemplate // TODO: field does not used anywhere
	List       []Value
	Node       string
	Variable   string
	Data       any `terminusgo:"class=xsd:anySimpleType"`
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
	*schema.TaggedUnionModel
	*schema.SubDocumentModel
	Node     string
	Variable string
}

func (v *NodeValue) FromVariableName(value string) {
	v.Variable = value
}

func (v *NodeValue) FromString(value string, _ bool) {
	v.Node = value
}

type DataValue struct {
	*schema.TaggedUnionModel
	*schema.SubDocumentModel
	List     []DataValue
	Data     any `terminusgo:"class=xsd:anySimpleType"`
	Variable string
}

func (v *DataValue) FromVariableName(value string) {
	v.Variable = value
}

func (v *DataValue) FromString(value string, preferLiteral bool) {
	if preferLiteral {
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
	*schema.TaggedUnionModel
	*schema.SubDocumentModel
	Name  string
	Index uint
}

type Column struct {
	*schema.SubDocumentModel
	Indicator Indicator
	Variable  string
	Type      *string
}

type Source struct {
	*schema.TaggedUnionModel
	*schema.SubDocumentModel
	Post string
	URL  string
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
	Separator     string
	IgnoreQuoutes bool
	Strip         bool
	SkipHeader    bool
	Convert       bool
	Case          CSVCase
	Functor       string
	Arity         uint
	MatchArity    bool
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
	BaseURL    string
	AnonPrefix string
	Format     TurtleFormat
	Resources  TurtleResources
	OnError    TurtleOnError
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
	Order    OrderDirection
	Variable string
}

type Literal struct {
	Type  string `terminusgo:"name=@type"`
	Value any
}

func (s *Literal) FromAnyValue(value any) {
	typ := reflect.ValueOf(value).Type()
	if cls, ok := schema.GetSchemaClass(typ); ok {
		s.Type = cls
		s.Value = value
		if conv, ok := schema.GetConverter(typ); ok {
			s.Value = conv(s.Value)
		}
		return
	}
	panic(fmt.Sprintf(
		"Cannot determine schema type of value with type %T, "+
			"maybe it's needed to define a type (see schema.DefineTypeClass() or schema.DefinePrimitiveTypeClass)?", value,
	))
}
