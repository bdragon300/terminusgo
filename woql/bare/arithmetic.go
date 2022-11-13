package bare

import "github.com/bdragon300/terminusgo/woql/schema"

func Plus(values ...schema.ArithmeticExpressionType) schema.Plus {
	return makeValuesNested(
		func(l, r schema.ArithmeticExpressionType) schema.ArithmeticExpressionType {
			return schema.Plus{Left: l, Right: r}
		},
		values...,
	).(schema.Plus)
}

func Minus(values ...schema.ArithmeticExpressionType) schema.Minus {
	return makeValuesNested(
		func(l, r schema.ArithmeticExpressionType) schema.ArithmeticExpressionType {
			return schema.Minus{Left: l, Right: r}
		},
		values...,
	).(schema.Minus)
}

func Times(values ...schema.ArithmeticExpressionType) schema.Times {
	return makeValuesNested(
		func(l, r schema.ArithmeticExpressionType) schema.ArithmeticExpressionType {
			return schema.Times{Left: l, Right: r}
		},
		values...,
	).(schema.Times)
}

func Divide(values ...schema.ArithmeticExpressionType) schema.Divide {
	return makeValuesNested(
		func(l, r schema.ArithmeticExpressionType) schema.ArithmeticExpressionType {
			return schema.Divide{Left: l, Right: r}
		},
		values...,
	).(schema.Divide)
}

func Div(values ...schema.ArithmeticExpressionType) schema.Div {
	return makeValuesNested(
		func(l, r schema.ArithmeticExpressionType) schema.ArithmeticExpressionType {
			return schema.Div{Left: l, Right: r}
		},
		values...,
	).(schema.Div)
}

func Exp(left, right schema.ArithmeticExpressionType) schema.Exp {
	return schema.Exp{
		Left:  left,
		Right: right,
	}
}

func Floor(value schema.ArithmeticExpressionType) schema.Floor {
	return schema.Floor{
		Argument: value,
	}
}

type arithWrapCallback func(l, r schema.ArithmeticExpressionType) schema.ArithmeticExpressionType

func makeValuesNested(wrapCb arithWrapCallback, values ...schema.ArithmeticExpressionType) schema.ArithmeticExpressionType {
	if len(values) < 2 {
		panic("Count of parameters must be 2 or more")
	}
	res := values[len(values)-1]
	for i := len(values) - 2; i >= 0; i-- {
		res = wrapCb(values[i], res)
	}
	return res
}
