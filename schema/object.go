package schema

import (
	"errors"
	"reflect"
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
		noType := false
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
			if n, s, ok := getStructFieldSchema(fld); ok {
				fldName = n
				fldTypeName = s.Class
				if s.Type == FieldTypeOptional && fldVal.IsZero() { // Acts like "omitempty" check
					continue
				}
				if _, ok2 := s.Tags["notypeobj"]; ok2 {
					noType = true
				}
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
			if noType || fld.Anonymous {
				delete(m, "@type")
			}
			if fld.Anonymous { // Extend res map
				for k, v := range m {
					buf[k] = v
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
