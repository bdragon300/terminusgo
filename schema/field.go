package schema

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobeam/stringy"
)

const tagName = "terminusgo"

// TODO: figure out what about sys:JSONDocument -- how to use it and implement

type FieldType string

const (
	FieldTypeOptional FieldType = "Optional"
	FieldTypeList     FieldType = "List"
	FieldTypeArray    FieldType = "Array"
	FieldTypeForeign  FieldType = "Foreign"
	FieldTypeSet      FieldType = "Set"
)

type Field struct {
	Type           FieldType `json:"@type,omitempty"`
	Class          string    `json:"@class,omitempty"`           // For all types except Foreign
	ID             string    `json:"@id,omitempty"`              // For Foreign type
	Cardinality    uint      `json:"@cardinality,omitempty"`     // For Set type
	MinCardinality uint      `json:"@min_cardinality,omitempty"` // For Set type
	MaxCardinality uint      `json:"@max_cardinality,omitempty"` // For Set type
	Dimensions     uint      `json:"@dimensions,omitempty"`      // For Array type

	Tags map[string]string `json:"-" mapstructure:"-"`
}

var excludedTypes = []reflect.Type{
	reflect.TypeOf(AbstractModel{}),
	reflect.TypeOf(SubDocumentModel{}),
}

func analyzeModel(mdlTyp reflect.Type) (parents []reflect.Type, grandparents []reflect.Type, fields map[string]Field) {
	// TODO: cache and circular embedding
	fields = make(map[string]Field)

	for i := 0; i < mdlTyp.NumField(); i++ {
		fld := mdlTyp.Field(i)
		fldTyp := fld.Type
		if !fld.IsExported() {
			continue // Skip unexported fields
		}

		if fld.Anonymous {
			switch fldTyp.Kind() {
			case reflect.Interface:
				panic(fmt.Sprintf("Unable to extract model parent type from interface %s.%s", mdlTyp, fldTyp))
			case reflect.Ptr:
				fldTyp = fldTyp.Elem() // Extract a type from pointer anon field
			}
			// Collect fields from all (possibly nested) parent models
			if fldTyp.Kind() == reflect.Struct && !typeExcluded(fldTyp) {
				parents = append(parents, fldTyp)
				ps, gps, fs := analyzeModel(fldTyp)
				grandparents = append(append(grandparents, gps...), ps...)
				for k, v := range fs {
					fields[k] = v
				}
			}
			continue
		}
		if name, s, ok := getFieldSchema(fld); ok {
			fields[name] = s
		}
	}
	return
}

func typeExcluded(typ reflect.Type) bool {
	for _, t := range excludedTypes {
		if t == typ {
			return true
		}
	}
	return false
}

func getFieldSchema(field reflect.StructField) (string, Field, bool) {
	// TODO: cache
	fldTyp := field.Type
	fldName := stringy.New(field.Name).SnakeCase().ToLower()
	schema := Field{}

	if !field.IsExported() {
		return fldName, schema, false // TODO: test it
	}

	schema.Tags = parseTags(field)
	if _, ok := schema.Tags["-"]; ok {
		return fldName, schema, false // Skip by user request
	}

	if fldTyp.Kind() == reflect.Ptr {
		fldTyp = fldTyp.Elem()
		schema.Type = FieldTypeOptional
	}
	for fldTyp.Kind() == reflect.Slice || fldTyp.Kind() == reflect.Array {
		schema.Dimensions++
		schema.Type = FieldTypeList
		fldTyp = fldTyp.Elem()
	}
	if schema.Dimensions > 1 {
		schema.Type = FieldTypeArray
	}
	if t, ok := GetSchemaClass(fldTyp); ok {
		schema.Class = t
	} else if fldTyp.Kind() == reflect.Struct {
		schema.Class = fldTyp.Name()
	}

	applyTags(&schema, schema.Tags)
	if schema.Type != FieldTypeForeign && schema.Class == "" {
		panic(fmt.Sprintf("Unable to determine class for field '%s %s', try to set it manually or mark field as ignored", field.Name, field.Type))
	}

	if n, ok := schema.Tags["name"]; ok {
		fldName = n
	}

	return fldName, schema, true
}

func parseTags(field reflect.StructField) (tags map[string]string) {
	tagVal, ok := field.Tag.Lookup(tagName)
	if !ok {
		return nil
	}

	tags = make(map[string]string)
	for _, part := range strings.Split(tagVal, ",") {
		if part == "" {
			continue
		}
		part = strings.Trim(part, " ")
		k, v, _ := strings.Cut(part, "=")
		tags[strings.Trim(k, " ")] = strings.Trim(v, " '")
	}
	return
}

func applyTags(schema *Field, tags map[string]string) {
	// TODO: add field name (for serialization, instead of json:xxx)
	if _, ok := tags["optional"]; ok {
		schema.Type = FieldTypeOptional
	}
	if _, ok := tags["nooptional"]; ok {
		schema.Type = ""
	}
	if _, ok := tags["foreign"]; ok {
		schema.Type = FieldTypeForeign
		schema.ID = schema.Class
		schema.Class = ""
	}
	containerTags := map[string]FieldType{
		"minCardinality": FieldTypeSet,
		"maxCardinality": FieldTypeSet,
		"cardinality":    FieldTypeSet,
		"dimensions":     FieldTypeList,
	}
	for k, v := range containerTags {
		if val, ok := tags[k]; ok {
			schema.Type = v
			if v, err := strconv.Atoi(val); err == nil {
				k := strings.ToUpper(string(k[0])) + k[1:] // strings.Title is deprecated, don't want to add a library only for one small thing
				reflect.ValueOf(schema).Elem().FieldByName(k).SetUint(uint64(v))
			} else {
				panic(fmt.Sprintf("Tag %s=%s is not an integer: %v", k, val, err))
			}
		}
	}
	if schema.Dimensions > 1 {
		schema.Type = FieldTypeArray
	}
	if val, ok := tags["type"]; ok {
		schema.Type = FieldType(val)
	}
	if val, ok := tags["class"]; ok {
		schema.Class = val
	}
}
