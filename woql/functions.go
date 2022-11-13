package woql

import (
	"math/big"

	"github.com/bdragon300/terminusgo/woql/schema"
	"golang.org/x/exp/constraints"
)

type Numbers interface {
	constraints.Integer | constraints.Float | big.Float | big.Int | schema.ArithmeticValue
}

func Div[T1, T2 Numbers](left T1, right T2) schema.Div {
	return schema.Div{
		Left:  *ParseNumber(left, &schema.ArithmeticValue{}),
		Right: *ParseNumber(right, &schema.ArithmeticValue{}),
	}
}
