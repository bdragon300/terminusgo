package bare

import "github.com/bdragon300/terminusgo/woql/schema"

func String(value string) schema.Literal {
	res := &schema.Literal{}
	res.FromAnyValue(value)
	return *res
}

func Boolean(value bool) schema.Literal {
	res := &schema.Literal{}
	res.FromAnyValue(value)
	return *res
}

func DateTime(value bool) schema.Literal {
	res := &schema.Literal{}
	res.FromAnyValue(value)
	return *res
}

func IRI(value string) schema.NodeValue {
	res := &schema.NodeValue{}
	res.FromString(value, false)
	return *res
}

