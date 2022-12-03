package schema

import (
	"github.com/bdragon300/terminusgo/schema"
)

type ArithmeticExpressionType interface {
	Arithmetic()
}

type ArithmeticExpression struct {
	*schema.SubDocumentModel
	*schema.AbstractModel
}

func (a ArithmeticExpression) Arithmetic() {}

type ArithmeticValue struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	ArithmeticExpression
	Data     any    `json:"data" terminusgo:"class=xsd:anySimpleType"`
	Variable string `json:"variable"`
}

func (v *ArithmeticValue) FromVariableName(value string) {
	v.Variable = value
}

func (v *ArithmeticValue) FromAnyValue(value any) {
	newVal := &Literal{}
	newVal.FromAnyValue(value)
	v.Data = *newVal
}

func (v *ArithmeticValue) FromString(value string, forceLiteral bool) {
	if forceLiteral {
		newVal := &Literal{}
		newVal.FromAnyValue(value)
		v.Data = *newVal
	} else {
		v.Data = value
	}
}

type Plus struct {
	ArithmeticExpression
	Left  ArithmeticExpressionType `json:"left" terminusgo:"type=Class,class=ArithmeticExpression"`
	Right ArithmeticExpressionType `json:"right" terminusgo:"type=Class,class=ArithmeticExpression"`
}

type Minus struct {
	ArithmeticExpression
	Left  ArithmeticExpressionType `json:"left" terminusgo:"type=Class,class=ArithmeticExpression"`
	Right ArithmeticExpressionType `json:"right" terminusgo:"type=Class,class=ArithmeticExpression"`
}

type Times struct {
	ArithmeticExpression
	Left  ArithmeticExpressionType `json:"left" terminusgo:"type=Class,class=ArithmeticExpression"`
	Right ArithmeticExpressionType `json:"right" terminusgo:"type=Class,class=ArithmeticExpression"`
}

type Divide struct {
	ArithmeticExpression
	Left  ArithmeticExpressionType `json:"left" terminusgo:"type=Class,class=ArithmeticExpression"`
	Right ArithmeticExpressionType `json:"right" terminusgo:"type=Class,class=ArithmeticExpression"`
}

type Div struct {
	ArithmeticExpression
	Left  ArithmeticExpressionType `json:"left" terminusgo:"type=Class,class=ArithmeticExpression"`
	Right ArithmeticExpressionType `json:"right" terminusgo:"type=Class,class=ArithmeticExpression"`
}

type Exp struct {
	ArithmeticExpression
	Left  ArithmeticExpressionType `json:"left" terminusgo:"type=Class,class=ArithmeticExpression"`
	Right ArithmeticExpressionType `json:"right" terminusgo:"type=Class,class=ArithmeticExpression"`
}

type Floor struct {
	ArithmeticExpression
	Argument ArithmeticExpressionType `json:"argument" terminusgo:"type=Class,class=ArithmeticExpression"`
}