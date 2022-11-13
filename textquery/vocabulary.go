package textquery

import (
	"github.com/bdragon300/terminusgo/woql"
)

type vocabularyItem func(query *woql.SimpleQueryBuilder, params []callParam) callResult

var functionsVocabulary = map[string]vocabularyItem{
	// "div": func(_ *woql.SimpleQueryBuilder, p []callParam) callResult {
	//	return query.Div(*woql.ParseNumber(p[0], &schema.ArithmeticValue{}), *woql.ParseNumber(p[1], &schema.ArithmeticValue{}))
	// },
	// "path": func(q *woql.SimpleQueryBuilder, p []callParam) callResult {
	//	return q.Path(p[0].(string), p[1].(string), p[2].(string))
	// },
	// "and": func(q *woql.SimpleQueryBuilder, p []callParam) callResult {
	//	return q.And(castToQueryMakerList(p)...)
	// },
	// "triple": func(q *woql.SimpleQueryBuilder, p []callParam) callResult {
	//	return q.Triple(p[0].(string), p[1].(string), p[2].(string))
	// },
}

// func castToQueryMakerList(values []callParam) (res []woql.QueryMaker) {
//	for _, v := range values {
//		res = append(res, v.(woql.QueryMaker))
//	}
//	return
// }
