package bare

import "github.com/bdragon300/terminusgo/woql/schema"

func String(value string) schema.SimpleValue {
	res := &schema.SimpleValue{}
	res.SetValue(value)
	return *res
}

func Boolean(value bool) schema.SimpleValue {
	res := &schema.SimpleValue{}
	res.SetValue(value)
	return *res
}

func DateTime(value bool) schema.SimpleValue {
	res := &schema.SimpleValue{}
	res.SetValue(value)
	return *res
}

func IRI(value string) schema.NodeValue {
	res := &schema.NodeValue{}
	res.FromString(value)
	return *res
}

