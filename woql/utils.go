package woql

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/bdragon300/terminusgo/woql/schema"
	"github.com/mitchellh/mapstructure"
)

type variableConvertable interface {
	FromVariableName(string)
}

type stringConvertable interface {
	FromString(string)
}

type numberConvertable interface {
	FromNumber(value any)
}

func parseVariable[ValueT variableConvertable](expr any, buf ValueT, requireVariable bool) ValueT {
	switch v := expr.(type) {
	case string:
		parts := strings.SplitN(v, ":", 2)
		if len(parts) > 1 && parts[0] == "v" {
			// TODO: validate variable name (not empty, valid characters, etc.)
			buf.FromVariableName(strings.TrimSpace(parts[1]))
		} else if requireVariable {
			panic(fmt.Sprintf("String %q is not a variable expression", expr)) // FIXME: return error instead of panic
		} else if obj, ok1 := any(buf).(stringConvertable); ok1 {
			if parts[0] == "v\\" {
				// Handle `v\:` expression (preventing string interpretation as variable)
				obj.FromString("v:" + parts[1])
			} else {
				obj.FromString(v)
			}
		} else {
			panic(fmt.Sprintf("Type %T is not convertable from variable expression", buf))
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, big.Float, big.Int:
		if obj, ok := any(buf).(numberConvertable); ok {
			obj.FromNumber(expr)
		} else {
			panic(fmt.Sprintf("Type %T is not convertable from number", buf))
		}
	default:
		panic(fmt.Sprintf("Type %T is not convertable from variable expression", expr))
	}
	return buf
}

func extractVariableName(expr string) string {
	parts := strings.SplitN(expr, ":", 2)
	if len(parts) > 1 && parts[0] == "v" {
		return parts[1]
	}
	panic(fmt.Sprintf("String %q is not a variable expression", expr)) // FIXME: return error instead of panic
}

func parseTriplePattern(expr string) (schema.PathPatternType, error) {
	return nil, nil // TODO
}

func parseNumber[T any, PT numberConvertable](value T, buf PT) PT {
	if bufVal, ok := any(buf).(*T); ok {
		*bufVal = value
		return buf
	} else if _, ok = any(value).(PT); ok {
		if err := mapstructure.Decode(value, buf); err != nil {
			panic(fmt.Sprintf("Error while copying struct: %v", err))
		}
	}

	buf.FromNumber(value)
	return buf
}
