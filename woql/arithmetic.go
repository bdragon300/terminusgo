package woql

import (
	"github.com/bdragon300/terminusgo/woql/bare"
	"github.com/bdragon300/terminusgo/woql/schema"
)

func Plus(values ...NumberOrVariable) schema.Plus {
	var params []schema.ArithmeticExpressionType
	for _, v := range values {
		params = append(params, *parseVariable(numOrVarWrapper{v}, &schema.ArithmeticValue{}, true))
	}
	return bare.Plus(params...)
}

func Minus(values ...NumberOrVariable) schema.Minus {
	var params []schema.ArithmeticExpressionType
	for _, v := range values {
		params = append(params, *parseVariable(numOrVarWrapper{v}, &schema.ArithmeticValue{}, true))
	}
	return bare.Minus(params...)
}

func Times(values ...NumberOrVariable) schema.Times {
	var params []schema.ArithmeticExpressionType
	for _, v := range values {
		params = append(params, *parseVariable(numOrVarWrapper{v}, &schema.ArithmeticValue{}, true))
	}
	return bare.Times(params...)
}

func Divide(values ...NumberOrVariable) schema.Divide {
	var params []schema.ArithmeticExpressionType
	for _, v := range values {
		params = append(params, *parseVariable(numOrVarWrapper{v}, &schema.ArithmeticValue{}, true))
	}
	return bare.Divide(params...)
}

func Div(values ...NumberOrVariable) schema.Div {
	var params []schema.ArithmeticExpressionType
	for _, v := range values {
		params = append(params, *parseVariable(numOrVarWrapper{v}, &schema.ArithmeticValue{}, true))
	}
	return bare.Div(params...)
}

func Exp(left, right NumberOrVariable) schema.Exp {
	return bare.Exp(
		*parseVariable(numOrVarWrapper{left}, &schema.ArithmeticValue{}, true),
		*parseVariable(numOrVarWrapper{right}, &schema.ArithmeticValue{}, true),
	)
}

func Floor(value NumberOrVariable) schema.Floor {
	return bare.Floor(*parseVariable(numOrVarWrapper{value}, &schema.ArithmeticValue{}, true))
}
