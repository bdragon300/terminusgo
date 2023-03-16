package srverror

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Union[T0, T1 any] struct {
	V0       T0
	V1       T1
	Selector uint8
}

func (ju *Union[T0, T1]) MarshalJSON() ([]byte, error) {
	switch ju.Selector {
	case 0:
		return json.Marshal(ju.V0)
	case 1:
		return json.Marshal(ju.V1)
	default:
		panic(fmt.Sprintf("Selector can be 0 or 1 only, got %d", ju.Selector))
	}
}

func (ju *Union[T0, T1]) UnmarshalJSON(bytes []byte) error {
	if err := json.Unmarshal(bytes, &ju.V0); err == nil {
		ju.Selector = 0
	} else if err = json.Unmarshal(bytes, &ju.V1); err == nil {
		ju.Selector = 1
	} else {
		return err
	}
	return nil
}

func ToUnion[T0, T1 any](v any, selector uint8) *Union[T0, T1] {
	val := reflect.ValueOf(v)
	zero0 := new(T0)
	zero1 := new(T1)
	if val.CanConvert(reflect.TypeOf(zero0).Elem()) {
		return &Union[T0, T1]{V0: v.(T0), V1: *zero1, Selector: 0}
	}
	if val.CanConvert(reflect.TypeOf(zero1).Elem()) {
		return &Union[T0, T1]{V0: *zero0, V1: v.(T1), Selector: 1}
	}
	panic(fmt.Sprintf("v is not convertable neither to %T nor to %T", zero0, zero1))
}
