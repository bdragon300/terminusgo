package schema

import (
	"encoding/base64"
	"math/big"
	"reflect"
	"time"
)

// https://www.w3.org/TR/xmlschema-2/#built-in-datatypes
var primitiveTypeClasses = map[reflect.Kind]string{
	reflect.Bool:    "xsd:boolean",
	reflect.Int:     "xsd:integer",
	reflect.Int8:    "xsd:byte",
	reflect.Int16:   "xsd:short",
	reflect.Int32:   "xsd:int",
	reflect.Int64:   "xsd:long",
	reflect.Uint:    "xsd:nonNegativeInteger",
	reflect.Uint8:   "xsd:unsignedByte",
	reflect.Uint16:  "xsd:unsignedShort",
	reflect.Uint32:  "xsd:unsignedInt",
	reflect.Uint64:  "xsd:unsignedLong",
	reflect.Float32: "xsd:float",
	reflect.Float64: "xsd:double",
	reflect.Map:     "sys:JSON",
	reflect.String:  "xsd:string",
	// TODO: sys.Unit ?
	// TODO: xdd:json ?
	// TODO: any as xsd:anySimpleType?
}

type ComplexTypeConverterFunc func(any) any

type complexTypeDefinition struct {
	schemaName            string
	toSchemaTypeConverter ComplexTypeConverterFunc
}

var complexTypeClasses map[reflect.Type]complexTypeDefinition

func init() {
	complexTypeClasses = map[reflect.Type]complexTypeDefinition{
		reflect.TypeOf(time.Duration(0)): {schemaName: "xsd:duration", toSchemaTypeConverter: func(v any) any {
			return v.(time.Duration).String() // TODO: ISO8601, ex: https://github.com/sosodev/duration
		}},
		reflect.TypeOf(time.Time{}): {schemaName: "xsd:dateTime", toSchemaTypeConverter: func(v any) any {
			return v.(time.Time).Format(time.RFC3339) // TODO: ISO8601
		}},
		reflect.TypeOf([]byte{}): {schemaName: "xsd:base64Binary", toSchemaTypeConverter: func(v any) any {
			return base64.StdEncoding.EncodeToString(v.([]byte))
		}},
		reflect.TypeOf(big.Float{}): {schemaName: "xsd:decimal", toSchemaTypeConverter: func(v any) any {
			val := v.(big.Float)
			res, _ := (&val).Float64()
			return res
		}},
		reflect.TypeOf(big.Int{}): {schemaName: "xsd:integer", toSchemaTypeConverter: func(v any) any {
			val := v.(big.Int)
			return (&val).Int64()
		}},
	}
}

func DefineTypeClass(goType reflect.Type, terminusClass string, toSchemaConverter ComplexTypeConverterFunc) {
	complexTypeClasses[goType] = complexTypeDefinition{
		schemaName:            terminusClass,
		toSchemaTypeConverter: toSchemaConverter,
	}
}

func DefinePrimitiveTypeClass(goKind reflect.Kind, terminusClass string) {
	primitiveTypeClasses[goKind] = terminusClass
}

func GetSchemaClass(typ reflect.Type) (string, bool) {
	if t, ok := complexTypeClasses[typ]; ok {
		return t.schemaName, true
	} else if t, ok := primitiveTypeClasses[typ.Kind()]; ok {
		return t, true
	}
	return "", false
}

func GetConverter(typ reflect.Type) (ComplexTypeConverterFunc, bool) {
	if t, ok := complexTypeClasses[typ]; ok {
		return t.toSchemaTypeConverter, true
	}
	return nil, false
}
