package bare

import "github.com/bdragon300/terminusgo/woql/schema"

func And(queries ...schema.Querier) schema.And {
	var query schema.Querier
	res := schema.And{}
	for _, query = range queries {
		switch v := query.(type) {
		case schema.And:
			res.And = append(res.And, v.And...) // Flatten ANDs
		default:
			res.And = append(res.And, v)
		}
	}
	return res
}

func Or(queries ...schema.Querier) schema.Or {
	var query schema.Querier
	res := schema.Or{}
	for _, query = range queries {
		switch v := query.(type) {
		case schema.Or:
			res.Or = append(res.Or, v.Or...) // Flatten ORs
		default:
			res.Or = append(res.Or, v)
		}
	}
	return res
}
