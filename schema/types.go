package schema

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

const tagName = "terminusgo"

type GenericType any // FieldSchema or string

type BaseModel struct{}

var baseModelTyp reflect.Type

// TODO: figure out what about sys:JSONDocument -- how to use it and implement

type FieldSchema struct {
	Type           string `mapstructure:"@type,omitempty"`
	Class          string `mapstructure:"@class,omitempty"`           // For all types except Foreign
	ID             string `mapstructure:"@id,omitempty"`              // For Foreign type
	Cardinality    uint   `mapstructure:"@cardinality,omitempty"`     //  For Set type
	MinCardinality uint   `mapstructure:"@min_cardinality,omitempty"` // For Set type
	MaxCardinality uint   `mapstructure:"@max_cardinality,omitempty"` // For Set type
	Dimensions     uint   `mapstructure:"@dimensions,omitempty"`      // For Array type
}

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
	// TODO: check how it checks type aliases, e.g. `type x int`
}

var complexTypeClasses map[reflect.Type]string

func init() {
	baseModelTyp = reflect.TypeOf(BaseModel{})
	complexTypeClasses = map[reflect.Type]string{
		reflect.TypeOf(time.Duration(0)): "xsd:duration",
		reflect.TypeOf(time.Time{}):      "xsd:dateTime",
		reflect.TypeOf([]byte{}):         "xsd:base64Binary",
	}
}

func SetGoTypeClass(goType reflect.Type, terminusClass string) {
	complexTypeClasses[goType] = terminusClass
}

func SetGoKindClass(goKind reflect.Kind, terminusClass string) {
	primitiveTypeClasses[goKind] = terminusClass
}

func GetModelSchema(mdl any) (string, map[string]GenericType, bool) {
	// TODO: cache and circular
	res := make(map[string]GenericType)
	typ := reflect.TypeOf(mdl)

	if getModelSchemaRecursive(typ, res) {
		return typ.Name(), res, true
	}
	return "", nil, false
}

func getFieldSchema(field reflect.StructField) GenericType {
	fTyp := field.Type
	schema := FieldSchema{}
	termClass := ""

	opts := parseTags(field)
	if _, ok := opts["-"]; ok {
		return nil // Skip by user request
	}

	if fTyp.Kind() == reflect.Ptr {
		fTyp = fTyp.Elem()
		schema.Type = "Optional"
	}
	for fTyp.Kind() == reflect.Slice || fTyp.Kind() == reflect.Array {
		schema.Dimensions++
		schema.Type = "List"
		fTyp = fTyp.Elem()
	}
	if schema.Dimensions > 1 {
		schema.Type = "Array"
	}
	if t, ok := complexTypeClasses[fTyp]; ok {
		termClass = t
	} else if t, ok := primitiveTypeClasses[fTyp.Kind()]; ok {
		termClass = t
	} else if fTyp.Kind() == reflect.Struct {
		termClass = fTyp.Name()
	}

	applyTags(&schema, &termClass, opts)
	if termClass == "" {
		panic(fmt.Sprintf("Unable to determine class for field '%s %s', try to set it manually or mark field to skip in schema", field.Name, field.Type))
	}

	schemaMap := make(map[string]any)
	if err := mapstructure.Decode(schema, &schemaMap); err != nil {
		panic(err)
	}
	if len(schemaMap) == 0 {
		return termClass // Only class has been set
	}
	schemaMap["@class"] = termClass
	return schemaMap
}

func getModelSchemaRecursive(typ reflect.Type, buf map[string]GenericType) bool {
	localBuf := make(map[string]GenericType)
	isModel := false

	for i := 0; i < typ.NumField(); i++ {
		fv := typ.Field(i)
		fvTyp := fv.Type

		if fv.Anonymous {
			if fvTyp.Kind() == reflect.Ptr {
				fvTyp = fvTyp.Elem()
			}
			if fvTyp == baseModelTyp {
				isModel = true
			} else if fvTyp.Kind() == reflect.Struct {
				isModel = isModel || getModelSchemaRecursive(fvTyp, localBuf)
			}
		} else {
			localBuf[fv.Name] = getFieldSchema(fv)
		}
	}

	if isModel {
		for k, v := range localBuf {
			if v != nil {
				buf[k] = v
			}
		}
	}
	return isModel
}

func parseTags(field reflect.StructField) (res map[string]string) {
	tagVal, ok := field.Tag.Lookup(tagName)
	if !ok {
		return nil
	}

	res = make(map[string]string)
	for _, part := range strings.Split(tagVal, ",") {
		if part == "" {
			continue
		}
		part = strings.Trim(part, " ")
		k, v, _ := strings.Cut(part, "=")
		res[strings.Trim(k, " ")] = strings.Trim(v, " '")
	}
	return
}

func applyTags(schema *FieldSchema, termClass *string, opts map[string]string) {
	if _, ok := opts["optional"]; ok {
		schema.Type = "Optional"
	}
	if _, ok := opts["foreign"]; ok {
		schema.Type = "Foreign"
		schema.ID = *termClass
		*termClass = ""
	}
	containerTags := [][2]string{
		{"Set", "minCardinality"}, {"Set", "maxCardinality"}, {"Set", "cardinality"}, {"List", "dimensions"},
	}
	for _, ct := range containerTags {
		if val, ok := opts[ct[1]]; ok {
			schema.Type = ct[0]
			if v, err := strconv.Atoi(val); err == nil {
				k := strings.ToUpper(string(ct[1][0])) + ct[1][1:] // strings.Title is deprecated, don't want to add a library only for one small thing
				reflect.ValueOf(schema).Elem().FieldByName(k).SetUint(uint64(v))
			} else {
				panic(fmt.Sprintf("Tag %s=%s is not an integer: %v", ct[1], val, err))
			}
		}
	}
	if schema.Dimensions > 1 {
		schema.Type = "Array"
	}
	if val, ok := opts["type"]; ok {
		schema.Type = val
	}
	if val, ok := opts["class"]; ok {
		*termClass = val
	}
}
