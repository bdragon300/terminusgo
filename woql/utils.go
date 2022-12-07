package woql

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/bdragon300/terminusgo/woql/schema"
)

var defaultVocabulary = map[string]string{
	"type":             "rdf:type",
	"label":            "rdfs:label",
	"Class":            "owl:Class",
	"DatatypeProperty": "owl:DatatypeProperty",
	"ObjectProperty":   "owl:ObjectProperty",
	"Document":         "terminus:Document",
	"abstract":         "terminus:Document",
	"comment":          "rdfs:comment",
	"range":            "rdfs:range",
	"domain":           "rdfs:domain",
	"subClassOf":       "rdfs:subClassOf",
	"boolean":          "xsd:boolean",
	"string":           "xsd:string",
	"integer":          "xsd:integer",
	"decimal":          "xsd:decimal",
	"email":            "xdd:email",
	"json":             "xdd:json",
	"dateTime":         "xsd:dateTime",
	"date":             "xsd:date",
	"coordinate":       "xdd:coordinate",
	"line":             "xdd:coordinatePolyline",
	"polygon":          "xdd:coordinatePolygon",
}

func fromVocab[T ~string](qb *QueryBuilder, val T) T {
	if v, ok := qb.vocabulary[string(val)]; ok {
		return T(v)
	}
	return val
}

type variableConvertable interface {
	FromVariableName(string)
}

type stringConvertable interface {
	FromString(string, bool)
}

type anyConvertable interface {
	FromAnyValue(value any)
}

func parseVariable[ValueT variableConvertable](expr any, buf ValueT, preferLiteral bool) ValueT { // TODO: remove onlyVariable
	// TODO: refactor, move common parts to a separate function
	switch v := expr.(type) {
	case string, Variable, StringOrVariable:
		strExpr := v.(string)
		varName, err := extractVariableName(strExpr)
		_, varOnly := v.(Variable)
		if err == nil {
			buf.FromVariableName(varName)
		} else if obj, ok1 := any(buf).(stringConvertable); ok1 && !varOnly {
			parseString(strExpr, obj, preferLiteral)
		} else {
			panic(fmt.Sprintf("Type %T is not convertable from variable or string", buf))
		}

	case intOrVarWrapper:
		switch v2 := v.v.(type) {
		case string:
			if varName, err := extractVariableName(v2); err != nil {
				panic(fmt.Sprintf("String %q is not a variable expression", expr)) // FIXME: return error instead of panic
			} else {
				buf.FromVariableName(varName)
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, big.Int:
			if obj, ok := any(buf).(anyConvertable); ok {
				obj.FromAnyValue(expr)
			} else {
				panic(fmt.Sprintf("Type %T is not convertable from integer", buf))
			}
		default:
			panic(fmt.Sprintf("%v is not an integer or variable", expr)) // FIXME: return error instead of panic
		}

	case numOrVarWrapper:
		switch v2 := v.v.(type) {
		case string:
			if varName, err := extractVariableName(v2); err != nil {
				panic(fmt.Sprintf("String %q is not a variable expression", expr)) // FIXME: return error instead of panic
			} else {
				buf.FromVariableName(varName)
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, big.Int, float32, float64, big.Float:
			if obj, ok := any(buf).(anyConvertable); ok {
				obj.FromAnyValue(expr)
			} else {
				panic(fmt.Sprintf("Type %T is not convertable from number", buf))
			}
		default:
			panic(fmt.Sprintf("%v is not a number or variable", expr)) // FIXME: return error instead of panic
		}

	default:
		switch v2 := v.(type) {
		case string:
			if varName, err := extractVariableName(v2); err == nil {
				buf.FromVariableName(varName)
			} else if obj, ok1 := any(buf).(stringConvertable); ok1 {
				parseString(v2, obj, preferLiteral)
			} else {
				panic(fmt.Sprintf("Type %T is not convertable from variable or string", buf))
			}
		default:
			if obj, ok := any(buf).(anyConvertable); ok {
				obj.FromAnyValue(v2)
			} else {
				panic(fmt.Sprintf("Type %T is not convertable to literal", buf))
			}
		}
	}
	return buf
}

func parseString[T stringConvertable](str string, buf T, preferLiteral bool) T {
	// `v\:` prefix instead of `v:` prevents string interpretation as variable
	if strings.HasPrefix(str, "v\\:") { // TODO: move to FromString
		str = strings.Replace(str, "v\\:", "v:", 1)
	}
	buf.FromString(str, preferLiteral)
	return buf
}

type variableTypes interface {
	string | Variable | StringOrVariable
}

func extractVariableName[T variableTypes](expr T) (string, error) {
	// TODO: validate variable name (not empty, valid characters etc.)
	parts := strings.SplitN(strings.TrimSpace(string(expr)), ":", 2)
	if len(parts) > 1 && parts[0] == "v" {
		return strings.TrimSpace(parts[1]), nil
	}
	return "", fmt.Errorf("string %q is not a variable expression", expr)
}

func parseTriplePattern(expr string) (schema.PathPatternType, error) {
	return nil, nil // TODO
}
