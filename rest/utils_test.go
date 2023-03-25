package rest_test

import (
	"reflect"

	"github.com/bdragon300/terminusgo/rest"
	"github.com/stretchr/testify/require"
)

var excludeAssertionFields = []string{"ID"}

// assertSameTerminusObjects recursively compares two terminus objects (the same type and fields values). Fields
// listed in `excludeFields` are skipped during comparison.
func assertSameTerminusObjects(r *require.Assertions, expect, actual rest.TerminusObject, excludeFields []string) {
	_assertSameTerminusObjects(r, expect, actual, excludeFields, "")
}

func _assertSameTerminusObjects(r *require.Assertions, expect, actual rest.TerminusObject, excludeFields []string, _path string) {
	ev := reflect.Indirect(reflect.ValueOf(expect))
	av := reflect.Indirect(reflect.ValueOf(actual))

	if ev.IsValid() != av.IsValid() {
		r.Failf("not equal", "%+v != %+v", expect, actual)
	}
	if ev.Type() != av.Type() {
		r.Failf("non-equal types", "%s: %v != %v", _path, ev.Type(), av.Type())
	}
	for i := 0; i < ev.NumField(); i++ {
		fld := ev.Type().Field(i)
		if !fld.IsExported() { // Exclude these fields from comparison
			continue
		}
		for _, n := range excludeFields {
			if fld.Name == n {
				continue
			}
		}
		p := _path + "." + fld.Name

		efval := ev.Field(i)
		afval := av.Field(i)
		ftyp := fld.Type
		for ftyp.Kind() == reflect.Interface || ftyp.Kind() == reflect.Pointer {
			if efval.IsNil() != afval.IsNil() {
				r.Failf("field not equal", "%s: %+v != %+v", p, efval.Interface(), afval.Interface())
			} else if afval.IsNil() {
				break
			}
			afval = afval.Elem()
			efval = efval.Elem()
			if efval.IsNil() != afval.IsNil() || efval.Type() != afval.Type() {
				r.Failf("field non-equal types", "%s: %v != %v", p, efval.Type(), afval.Type())
			}
			ftyp = efval.Type()
		}
		if ftyp.Kind() == reflect.Struct {
			_assertSameTerminusObjects(r, efval.Interface(), afval.Interface(), excludeFields, p)
		} else {
			r.Equalf(efval.Interface(), afval.Interface(), "field not equal", "%s: %+v != %+v", p, efval.Interface(), afval.Interface())
		}
	}
}
