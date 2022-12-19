package textquery

import (
	"reflect"

	"github.com/bdragon300/terminusgo/woql"
	"github.com/bdragon300/terminusgo/woql/schema"
)

type vocabularyItem func(query *woql.QueryBuilder, params []callParam) callResult

// TODO: Querier (and possible other interface params) may be nil -- consider that
var functionsVocabulary = map[string]vocabularyItem{
	"plus": func(_ *woql.QueryBuilder, p []callParam) callResult {
		return woql.Plus(paramVariadic[woql.NumberOrVariable](p)...)
	},
	"minus": func(_ *woql.QueryBuilder, p []callParam) callResult {
		return woql.Minus(paramVariadic[woql.NumberOrVariable](p)...)
	},
	"times": func(_ *woql.QueryBuilder, p []callParam) callResult {
		return woql.Times(paramVariadic[woql.NumberOrVariable](p)...)
	},
	"divide": func(_ *woql.QueryBuilder, p []callParam) callResult {
		return woql.Divide(paramVariadic[woql.NumberOrVariable](p)...)
	},
	"div": func(_ *woql.QueryBuilder, p []callParam) callResult {
		return woql.Div(paramVariadic[woql.NumberOrVariable](p)...)
	},
	"exp": func(_ *woql.QueryBuilder, p []callParam) callResult {
		return woql.Exp(param2[woql.NumberOrVariable, woql.NumberOrVariable](p))
	},
	"floor": func(_ *woql.QueryBuilder, p []callParam) callResult {
		return woql.Floor(param1[woql.NumberOrVariable](p))
	},
	// TODO: implement parsing dicts
	// "file": func(_ *woql.QueryBuilder, p []callParam) callResult {
	//	return woql.File(param2[string, schema.FileOptions](p))
	// },
	// "remote": func(_ *woql.QueryBuilder, p []callParam) callResult {
	//	return woql.Remote(param2[string, schema.FileOptions](p))
	// },
	// "post": func(_ *woql.QueryBuilder, p []callParam) callResult {
	//	return woql.Post(param2[string, schema.FileOptions](p))
	// },
	"and": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.And(paramVariadic[schema.Querier](p)...)
	},
	"not": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Not(param1[schema.Querier](p))
	},
	"select": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Select(paramVariadic[woql.Variable](p)...)
	},
	"from": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.From(param2[string, schema.Querier](p))
	},
	"using": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Using(param2[string, schema.Querier](p))
	},
	"distinct": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Distinct(paramVariadic[woql.Variable](p)...)
	},
	"into": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Into(param2[string, schema.Querier](p))
	},
	// TODO: implement parsing dict
	// "order_by": func(q *woql.QueryBuilder, p []callParam) callResult {
	//	return q.OrderBy(param1[map[woql.Variable]schema.OrderDirection](p))
	// },
	// TODO: implement parsing list
	// "group_by": func(q *woql.QueryBuilder, p []callParam) callResult {
	//	return q.GroupBy(param4[[]woql.Variable, []woql.Variable, woql.Variable, schema.Querier](p))
	// },
	"triple_count": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.TripleCount(param2[string, woql.IntegerOrVariable](p))
	},
	"triple": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Triple(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"add_triple": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.AddTriple(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"added_triple": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.AddedTriple(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"delete_triple": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.DeleteTriple(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"removed_triple": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.RemovedTriple(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"update_triple": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.UpdateTriple(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"quad": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Quad(param4[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"add_quad": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.AddQuad(param4[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"added_quad": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.AddQuad(param4[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"delete_quad": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.DeleteQuad(param4[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"removed_quad": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.RemovedQuad(param4[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"update_quad": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.UpdateQuad(param4[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"subsumption": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Subsumption(param2[woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"equals": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Equals(param2[woql.AnyOrVariable, woql.AnyOrVariable](p))
	},
	"substring": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{3, 5})
		switch len(p) {
		case 3:
			return q.SubString(param3[woql.StringOrVariable, woql.IntegerOrVariable, woql.StringOrVariable](p))
		case 5:
			return q.SubStringBeforeAfter(param5[woql.StringOrVariable, woql.IntegerOrVariable, woql.IntegerOrVariable, woql.IntegerOrVariable, woql.StringOrVariable](p))
		}
		panic("We must not be there, there is a programming flaw")
	},
	// TODO: implement parsing dict
	// "update_document": func(q *woql.QueryBuilder, p []callParam) callResult {
	//	return q.UpdateDocument(param2[map[string]any, woql.StringOrVariable](p))
	// },
	// TODO: implement parsing dict
	// "insert_document": func(q *woql.QueryBuilder, p []callParam) callResult {
	//	return q.InsertDocument(param2[map[string]any, woql.StringOrVariable](p))
	// },
	"delete_document": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.DeleteDocument(param1[woql.StringOrVariable](p))
	},
	"read_document": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.ReadDocument(param2[woql.StringOrVariable, woql.Variable](p))
	},
	"once": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Once(param1[schema.Querier](p))
	},
	"trim": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Trim(param2[woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"upper": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Upper(param2[woql.StringOrVariable, woql.Variable](p))
	},
	"lower": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Lower(param2[woql.StringOrVariable, woql.Variable](p))
	},
	"concatenate": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{2})
		p0Kind := reflect.TypeOf(p[0]).Kind()
		switch p0Kind {
		case reflect.Slice, reflect.Array:
			return q.ConcatenateList(param2[[]woql.StringOrVariable, woql.StringOrVariable](p))
		default:
			return q.Concatenate(param2[woql.Variable, woql.StringOrVariable](p))
		}
	},
	"join": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{3})
		p0Kind := reflect.TypeOf(p[0]).Kind()
		switch p0Kind {
		case reflect.Slice, reflect.Array:
			return q.JoinList(param3[[]woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
		default:
			return q.Join(param3[woql.Variable, woql.StringOrVariable, woql.StringOrVariable](p))
		}
	},
	"split": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Split(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"regexp": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{3})
		p2Kind := reflect.TypeOf(p[2]).Kind()
		switch p2Kind {
		case reflect.Slice, reflect.Array:
			return q.RegexpToList(param3[woql.StringOrVariable, woql.StringOrVariable, []woql.StringOrVariable](p))
		default:
			return q.Regexp(param3[woql.StringOrVariable, woql.StringOrVariable, woql.Variable](p))
		}
	},
	"eval": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Eval(param2[schema.ArithmeticExpressionType, woql.Variable](p))
	},
	"isa": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.IsA(param2[woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"like": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Like(param3[woql.StringOrVariable, woql.StringOrVariable, woql.NumberOrVariable](p))
	},
	"less": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Less(param2[woql.NumberOrVariable, woql.NumberOrVariable](p))
	},
	"greater": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Greater(param2[woql.NumberOrVariable, woql.NumberOrVariable](p))
	},
	"optional": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Optional(param1[schema.Querier](p))
	},
	"unique": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Unique(param3[string, []woql.StringOrVariable, woql.Variable](p))
	},
	"idgen": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.IDGen(param3[string, []woql.StringOrVariable, woql.Variable](p))
	},
	"pad": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Pad(param4[woql.StringOrVariable, woql.StringOrVariable, woql.IntegerOrVariable, woql.Variable](p))
	},
	"dot": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Dot(param3[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"member": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{2})
		p1Kind := reflect.TypeOf(p[1]).Kind()
		switch p1Kind {
		case reflect.Slice, reflect.Array:
			return q.MemberOfList(param2[woql.AnyOrVariable, []woql.AnyOrVariable](p))
		default:
			return q.Member(param2[woql.AnyOrVariable, woql.AnyOrVariable](p))
		}
	},
	"sum": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{2})
		p0Kind := reflect.TypeOf(p[0]).Kind()
		switch p0Kind {
		case reflect.Slice, reflect.Array:
			return q.SumList(param2[[]woql.NumberOrVariable, woql.Variable](p))
		default:
			return q.Sum(param2[woql.Variable, woql.Variable](p))
		}
	},
	"start": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Start(param2[uint, schema.Querier](p))
	},
	"limit": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Limit(param2[uint, schema.Querier](p))
	},
	"length": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{2})
		p0Kind := reflect.TypeOf(p[0]).Kind()
		switch p0Kind {
		case reflect.Slice, reflect.Array:
			return q.LengthList(param2[[]woql.AnyOrVariable, woql.IntegerOrVariable](p))
		default:
			return q.Length(param2[woql.Variable, woql.IntegerOrVariable](p))
		}
	},
	"immediately": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Immediately(param1[schema.Querier](p))
	},
	"count": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Count(param2[woql.IntegerOrVariable, schema.Querier](p))
	},
	"typecast": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.TypeCast(param3[woql.AnyOrVariable, woql.StringOrVariable, woql.Variable](p))
	},
	"typeof": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.TypeOf(param2[string, string](p))
	},
	"true": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{0})
		return q.True()
	},
	"path": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Path(param4[woql.StringOrVariable, string, woql.StringOrVariable, woql.Variable](p))
	},
	"size": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Size(param2[string, woql.Variable](p))
	},
	"star": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Star(param4[string, woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable](p))
	},
	"all": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.All(param4[woql.StringOrVariable, woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"insert": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Insert(param3[woql.StringOrVariable, woql.StringOrVariable, string](p))
	},
	"nuke": func(q *woql.QueryBuilder, p []callParam) callResult {
		return q.Nuke(param1[string](p))
	},
	"query": func(q *woql.QueryBuilder, p []callParam) callResult {
		checkParamCount(p, []int{0})
		return q.Query()
	},
}
