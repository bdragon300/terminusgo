package woql

import (
	"fmt"
	"strings"

	"github.com/bdragon300/terminusgo/woql/schema"
	"github.com/mitchellh/mapstructure"
)

type VariableConvertable interface {
	FromVariable(string)
}

type StringConvertable interface {
	FromString(string)
}

func ParseVariable[ValueT VariableConvertable](expr string, buf ValueT) ValueT {
	parts := strings.SplitN(expr, ":", 2)
	if len(parts) > 1 && parts[0] == "v" {
		buf.FromVariable(parts[1])
	} else if obj, ok := any(buf).(StringConvertable); ok {
		obj.FromString(expr)
	} else {
		panic(fmt.Sprintf("Type %T is not convertable from string", buf))
	}
	return buf
}

func ParseTriplePattern(expr string) (schema.PathPatternType, error) {
	return nil, nil // TODO
}

type NumberConvertable interface {
	FromNumber(value any)
}

func ParseNumber[T any, PT NumberConvertable](value T, buf PT) PT {
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
