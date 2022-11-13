package schema

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
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
	Type           FieldType `json:"@type,omitempty" mapstructure:"@type,omitempty"`
	Class          string    `json:"@class,omitempty" mapstructure:"@class,omitempty"`                     // For all types except Foreign
	ID             string    `json:"@id,omitempty" mapstructure:"@id,omitempty"`                           // For Foreign type
	Cardinality    uint      `json:"@cardinality,omitempty" mapstructure:"@cardinality,omitempty"`         // For Set type
	MinCardinality uint      `json:"@min_cardinality,omitempty" mapstructure:"@min_cardinality,omitempty"` // For Set type
	MaxCardinality uint      `json:"@max_cardinality,omitempty" mapstructure:"@max_cardinality,omitempty"` // For Set type
	Dimensions     uint      `json:"@dimensions,omitempty" mapstructure:"@dimensions,omitempty"`           // For Array type
}

func analyzeModel(mdlTyp reflect.Type) (parents []reflect.Type, grandparents []reflect.Type, fields map[string]Field) {
	// TODO: cache and circular embedding
	fields = make(map[string]Field)

	for i := 0; i < mdlTyp.NumField(); i++ {
		fld := mdlTyp.Field(i)
		fldTyp := fld.Type

		if fld.Anonymous {
			// Extract a type from pointer anon field
			if fldTyp.Kind() == reflect.Ptr {
				fldTyp = fldTyp.Elem()
			}
			// Collect fields from all (possibly nested) parent models
			if fldTyp.Kind() == reflect.Struct {
				parents = append(parents, fldTyp) // FIXME: exclude AbstractModel, SubDocumentModel, etc.
				ps, gps, fs := analyzeModel(fldTyp)
				grandparents = append(append(grandparents, gps...), ps...)
				for k, v := range fs {
					fields[k] = v
				}
			}
			continue
		}
		if s, ok := getFieldSchema(fld); ok {
			fields[fld.Name] = s
		}
	}
	return
}

func getFieldSchema(field reflect.StructField) (Field, bool) {
	fltTyp := field.Type
	schema := Field{}

	if !unicode.IsUpper(rune(field.Name[0])) {
		return schema, false // Skip private fields
	}

	opts := parseTags(field)
	if _, ok := opts["-"]; ok {
		return schema, false // Skip by user request
	}

	if fltTyp.Kind() == reflect.Ptr {
		fltTyp = fltTyp.Elem()
		schema.Type = FieldTypeOptional
	}
	for fltTyp.Kind() == reflect.Slice || fltTyp.Kind() == reflect.Array {
		schema.Dimensions++
		schema.Type = FieldTypeList
		fltTyp = fltTyp.Elem()
	}
	if schema.Dimensions > 1 {
		schema.Type = FieldTypeArray
	}
	if t, ok := GetSchemaClass(fltTyp); ok {
		schema.Class = t
	} else if fltTyp.Kind() == reflect.Struct {
		schema.Class = fltTyp.Name()
	}

	applyTags(&schema, opts)
	if schema.Type != FieldTypeForeign && schema.Class == "" {
		panic(fmt.Sprintf("Unable to determine class for field '%s %s', try to set it manually or mark field as ignored", field.Name, field.Type))
	}

	return schema, true
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

func applyTags(schema *Field, opts map[string]string) {
	// TODO: add nooptional
	if _, ok := opts["optional"]; ok {
		schema.Type = FieldTypeOptional
	}
	if _, ok := opts["foreign"]; ok {
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
		if val, ok := opts[k]; ok {
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
	if val, ok := opts["type"]; ok {
		schema.Type = FieldType(val)
	}
	if val, ok := opts["class"]; ok {
		schema.Class = val
	}
}
