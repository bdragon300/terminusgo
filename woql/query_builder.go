package woql

import (
	"fmt"

	"github.com/bdragon300/terminusgo/woql/bare"
	"github.com/bdragon300/terminusgo/woql/schema"
)

// TODO: check query builder against javascript client and terminus query schema
// TODO: in some methods the parameters can actually be a one var, var list or string. Do smth with this

func NewSimpleQueryBuilder() *QueryBuilder {
	return &QueryBuilder{bare.NewQueryBuilder()}
}

type (
	Variable          string
	StringOrVariable  string
	IntegerOrVariable any // Must be wrapped by intOrVarWrapper before parseVariable call
	NumberOrVariable  any //  Must be wrapped by numOrVarWrapper before parseVariable call
	AnyOrVariable     any
)

// intOrVarWrapper is a wrapper for IntegerOrVariable variable, this helps to distinguish a IntegerOrVariable type
// from other "any" aliases
type intOrVarWrapper struct {
	v any
}

type numOrVarWrapper struct {
	v any
}

type QueryBuilder struct {
	Bare *bare.QueryBuilder
	// TODO: errors member (multierror package)
}

func (b *QueryBuilder) GetQuery() schema.Querier {
	return b.Bare.GetQuery()
}

func (b *QueryBuilder) Clone() *QueryBuilder {
	return wrapBareQB(b.Bare.Clone())
}

func (b *QueryBuilder) And(subQueries ...schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.And(subQueries...))
}

func (b *QueryBuilder) Not(query schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Not(query))
}

func (b *QueryBuilder) Select(vars ...Variable) *QueryBuilder {
	params := make([]string, 0)
	for _, v := range vars {
		if name, err := extractVariableName(v); err == nil {
			params = append(params, name)
		} else {
			panic(fmt.Sprintf("Value %q is not a variable", v)) // TODO: keep error instead of panic
		}
	}
	return wrapBareQB(b.Bare.Select(params))
}

func (b *QueryBuilder) From(graph string, subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.From(graph, subQuery))
}

func (b *QueryBuilder) Using(collection string, subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Using(collection, subQuery))
}

func (b *QueryBuilder) Distinct(vars ...Variable) *QueryBuilder {
	params := make([]string, 0)
	for _, v := range vars {
		if name, err := extractVariableName(v); err == nil {
			params = append(params, name)
		} else {
			panic(fmt.Sprintf("Value %q is not a variable", v)) // TODO: keep error instead of panic
		}
	}
	return wrapBareQB(b.Bare.Distinct(params))
}

func (b *QueryBuilder) Into(graph string, subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Into(graph, subQuery))
}

func (b *QueryBuilder) OrderBy(vars map[Variable]schema.OrderDirection) *QueryBuilder {
	ordering := make([]schema.OrderTemplate, 0)
	for k, v := range vars {
		if name, err := extractVariableName(k); err == nil {
			ordering = append(ordering, schema.OrderTemplate{
				Order:    v,
				Variable: name,
			})
		} else {
			panic(fmt.Sprintf("Value %q is not a variable", v)) // TODO: keep error instead of panic
		}
	}
	return wrapBareQB(b.Bare.OrderBy(ordering))
}

func (b *QueryBuilder) GroupBy(groupVars []Variable, templateVars []Variable, outputVar Variable, subQuery schema.Querier) *QueryBuilder {
	grpVars := make([]schema.NodeValue, 0)
	for _, v := range groupVars {
		grpVars = append(grpVars, *parseVariable(v, &schema.NodeValue{}, false))
	}
	tplVars := make([]schema.Value, 0)
	for _, v := range templateVars {
		tplVars = append(tplVars, *parseVariable(v, &schema.Value{}, false))
	}
	return wrapBareQB(b.Bare.GroupBy(
		grpVars,
		schema.Value{List: tplVars},
		*parseVariable(outputVar, &schema.Value{}, false),
		subQuery,
	))
}

func (b *QueryBuilder) TripleCount(resourceID string, countVar IntegerOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.TripleCount(
		resourceID,
		*parseVariable(intOrVarWrapper{countVar}, &schema.DataValue{}, false), // FIXME: countVar in python client has type Value, figure out why
	))
}

func (b *QueryBuilder) Triple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Triple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) AddTriple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.AddTriple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) DeleteTriple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.DeleteTriple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) AddedTriple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.AddedTriple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) Quad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.Quad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) AddQuad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.AddQuad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
		graph, // TODO: maybe it's needed to transform graph id? See cleanGraph comment in js client
	))
}

func (b *QueryBuilder) DeleteQuad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.DeleteQuad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) AddedQuad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.AddedQuad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false), // TODO: vocab
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) Subsumption(parent, child StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Subsumption(
		*parseVariable(parent, &schema.NodeValue{}, false),
		*parseVariable(child, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) Equals(left, right AnyOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Equals(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) SubString(string StringOrVariable, length IntegerOrVariable, substring StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.SubString(
		*parseVariable(string, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{length}, &schema.DataValue{}, true),
		*parseVariable(substring, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{uint(0)}, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{uint(0)}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) SubStringBeforeAfter(string StringOrVariable, before, length, after IntegerOrVariable, substring StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.SubString(
		*parseVariable(string, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{length}, &schema.DataValue{}, true),
		*parseVariable(substring, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{before}, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{after}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) UpdateDocument(document map[string]any, iri StringOrVariable) *QueryBuilder {
	var param *schema.NodeValue
	if iri != "" {
		param = parseVariable(iri, &schema.NodeValue{}, false)
	}
	return wrapBareQB(b.Bare.UpdateDocument(
		*parseVariable(document, &schema.Value{}, false),
		param,
	))
}

func (b *QueryBuilder) InsertDocument(document map[string]any, iri StringOrVariable) *QueryBuilder {
	var param *schema.NodeValue
	if iri != "" {
		param = parseVariable(iri, &schema.NodeValue{}, false)
	}
	return wrapBareQB(b.Bare.InsertDocument(
		*parseVariable(document, &schema.Value{}, false),
		param,
	))
}

func (b *QueryBuilder) DeleteDocument(iri StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.DeleteDocument(
		*parseVariable(iri, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) ReadDocument(iri StringOrVariable, outputVar Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.ReadDocument(
		*parseVariable(iri, &schema.NodeValue{}, false),
		*parseVariable(outputVar, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) File(url string, options schema.FileOptions) *QueryBuilder {
	source := schema.Source{File: url}
	format := schema.FormatTypeCSV
	if options != nil {
		format = options.FileFormatType()
	}
	return wrapBareQB(b.Bare.QueryResource(source, format, options))
}

func (b *QueryBuilder) Once(subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Once(subQuery))
}

func (b *QueryBuilder) Remote(url string, options schema.FileOptions) *QueryBuilder {
	source := schema.Source{URI: url}
	format := schema.FormatTypeCSV
	if options != nil {
		format = options.FileFormatType()
	}
	return wrapBareQB(b.Bare.QueryResource(source, format, options))
}

func (b *QueryBuilder) Post(url string, options schema.FileOptions) *QueryBuilder {
	source := schema.Source{Post: url}
	format := schema.FormatTypeCSV
	if options != nil {
		format = options.FileFormatType()
	}
	return wrapBareQB(b.Bare.QueryResource(source, format, options))
}

func (b *QueryBuilder) Trim(untrimmed, trimmed StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Trim(
		*parseVariable(untrimmed, &schema.DataValue{}, true),
		*parseVariable(trimmed, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Upper(left StringOrVariable, right Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.Upper(
		*parseVariable(left, &schema.DataValue{}, true),
		*parseVariable(right, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Lower(left StringOrVariable, right Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.Lower(
		*parseVariable(left, &schema.DataValue{}, true),
		*parseVariable(right, &schema.DataValue{}, true),
	))
}

//// TODO: maybe it's needed a version with var list and one var, not only with DataValue.List
// func (b *QueryBuilder) ConcatenateTemplate(varExpr string, result StringOrVariable) *QueryBuilder {
//	var params []schema.DataValue
//
//	// Split up a string to variables and the regular text blocks and push to param list with appropriate type
//	// E.g "  lorem ipsum v:var1v:var2 dolor sit      v:  var3 amet " is transformed to
//	// "  lorem ipsum" + "v:var1" + "v:var2" + " dolor sit      " + "v:var3" + " amet "
//	parts := strings.Split(varExpr, "v:")
//	for _, part := range parts {
//		fields := strings.Fields(part)
//		if len(fields) < 1 {
//			panic("unnamed variable in concatenate expression") // TODO: return error instead of panic
//		}
//		params = append(params, *parseVariable("v:"+fields[0], &schema.DataValue{}, true))
//
//		part = part[strings.Index(part, fields[0])+len(fields[0]):]
//		if part != "" {
//			params = append(params, *parseVariable(part, &schema.DataValue{}, true))
//		}
//	}
//	return wrapBareQB(b.Bare.Concatenate(
//		schema.DataValue{List: params},
//		*parseVariable(result, &schema.DataValue{}, true),
//	))
//}

func (b *QueryBuilder) ConcatenateList(stringsOrVars []StringOrVariable, result StringOrVariable) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range stringsOrVars {
		params = append(params, *parseVariable(v, &schema.DataValue{}, true))
	}
	return wrapBareQB(b.Bare.Concatenate(
		schema.DataValue{List: params},
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Concatenate(listVar Variable, result StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Concatenate(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) JoinList(stringsOrVars []StringOrVariable, separator, result StringOrVariable) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range stringsOrVars {
		params = append(params, *parseVariable(v, &schema.DataValue{}, true))
	}
	return wrapBareQB(b.Bare.Join(
		schema.DataValue{List: params},
		*parseVariable(separator, &schema.DataValue{}, true),
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Join(listVar Variable, separator, result StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Join(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(separator, &schema.DataValue{}, true),
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Split(varName, separator, resultVar StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Split(
		*parseVariable(varName, &schema.DataValue{}, true),
		*parseVariable(separator, &schema.DataValue{}, true),
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Regexp(pattern, str StringOrVariable, result Variable) *QueryBuilder {
	var resParam *schema.DataValue
	if result != "" {
		resParam = parseVariable(result, &schema.DataValue{}, true)
	}
	return wrapBareQB(b.Bare.Regexp(
		*parseVariable(pattern, &schema.DataValue{}, true),
		*parseVariable(str, &schema.DataValue{}, true),
		resParam,
	))
}

func (b *QueryBuilder) RegexpToList(pattern, str StringOrVariable, result []StringOrVariable) *QueryBuilder {
	var resParam *schema.DataValue
	if len(result) > 0 {
		resParam = &schema.DataValue{}
		for _, v := range result {
			resParam.List = append(resParam.List, *parseVariable(v, &schema.DataValue{}, true))
		}
	}
	return wrapBareQB(b.Bare.Regexp(
		*parseVariable(pattern, &schema.DataValue{}, true),
		*parseVariable(str, &schema.DataValue{}, true),
		resParam,
	))
}

func (b *QueryBuilder) Eval(arith schema.ArithmeticExpressionType, result Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.Eval(
		arith,
		*parseVariable(result, &schema.ArithmeticValue{}, true),
	))
}

func (b *QueryBuilder) IsA(element, typ StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.IsA(
		*parseVariable(element, &schema.NodeValue{}, false),
		*parseVariable(typ, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) Like(left, right StringOrVariable, similarity NumberOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Like(
		*parseVariable(left, &schema.DataValue{}, true),
		*parseVariable(right, &schema.DataValue{}, true),
		*parseVariable(numOrVarWrapper{similarity}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Less(left, right NumberOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Less(
		*parseVariable(numOrVarWrapper{left}, &schema.DataValue{}, false),
		*parseVariable(numOrVarWrapper{right}, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Greater(left, right NumberOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Greater(
		*parseVariable(numOrVarWrapper{left}, &schema.DataValue{}, false),
		*parseVariable(numOrVarWrapper{right}, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Optional(subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Optional(subQuery))
}

func (b *QueryBuilder) Unique(prefix string, keyList []StringOrVariable, resultVar Variable) *QueryBuilder {
	var keys []schema.DataValue
	for _, v := range keyList {
		keys = append(keys, *parseVariable(v, &schema.DataValue{}, true))
	}
	return wrapBareQB(b.Bare.HashKey(
		*parseString(prefix, &schema.DataValue{}, true),
		keys,
		*parseVariable(resultVar, &schema.NodeValue{}, true),
	))
}

func (b *QueryBuilder) IDGen(prefix string, keyList []StringOrVariable, resultVar Variable) *QueryBuilder {
	var keys []schema.DataValue
	for _, v := range keyList {
		keys = append(keys, *parseVariable(v, &schema.DataValue{}, true))
	}
	return wrapBareQB(b.Bare.LexicalKey(
		*parseString(prefix, &schema.DataValue{}, true),
		keys,
		*parseVariable(resultVar, &schema.NodeValue{}, true),
	))
}

func (b *QueryBuilder) Pad(input, pad StringOrVariable, length IntegerOrVariable, resultVar Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.Pad(
		*parseVariable(input, &schema.DataValue{}, true),
		*parseVariable(pad, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{length}, &schema.DataValue{}, true),
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Dot(document, field, value StringOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Dot(
		*parseVariable(document, &schema.DataValue{}, false),
		*parseVariable(field, &schema.DataValue{}, true),
		*parseVariable(value, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Member(member, list AnyOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Member(
		*parseVariable(member, &schema.DataValue{}, false),
		*parseVariable(list, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) MemberOfList(member, list []AnyOrVariable) *QueryBuilder {
	var paramList []schema.DataValue
	for _, v := range list {
		paramList = append(paramList, *parseVariable(v, &schema.DataValue{}, false))
	}
	return wrapBareQB(b.Bare.Member(
		*parseVariable(member, &schema.DataValue{}, false),
		schema.DataValue{List: paramList},
	))
}

func (b *QueryBuilder) SumList(list []NumberOrVariable, outputVar Variable) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range list {
		params = append(params, *parseVariable(numOrVarWrapper{v}, &schema.DataValue{}, true))
	}
	return wrapBareQB(b.Bare.Sum(
		schema.DataValue{List: params},
		*parseVariable(outputVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Sum(listVar, outputVar Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.Sum(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(outputVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Start(start uint, subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Start(start, subQuery))
}

func (b *QueryBuilder) Limit(limit uint, subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Limit(limit, subQuery))
}

func (b *QueryBuilder) Length(listVar Variable, resultVar IntegerOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Length(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{resultVar}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) LengthList(list []AnyOrVariable, resultVar IntegerOrVariable) *QueryBuilder {
	return wrapBareQB(b.Bare.Length(
		*parseVariable(list, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{resultVar}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Immediately(subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Immediately(subQuery))
}

func (b *QueryBuilder) Count(countVar IntegerOrVariable, subQuery schema.Querier) *QueryBuilder {
	return wrapBareQB(b.Bare.Count(*parseVariable(countVar, &schema.DataValue{}, false), subQuery))
}

func (b *QueryBuilder) TypeCast(value AnyOrVariable, typ StringOrVariable, resultVar Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.TypeCast(
		*parseVariable(value, &schema.Value{}, false),
		*parseVariable(typ, &schema.NodeValue{}, false),
		*parseVariable(resultVar, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) TypeOf(value, typ string) *QueryBuilder {
	return wrapBareQB(b.Bare.TypeOf(
		*parseString(value, &schema.Value{}, true),
		*parseString(typ, &schema.NodeValue{}, true),
	))
}

func (b *QueryBuilder) True() *QueryBuilder {
	return wrapBareQB(b.Bare.True())
}

func (b *QueryBuilder) Path(subj StringOrVariable, pattern string, obj StringOrVariable, resultVar Variable) *QueryBuilder {
	parsedPattern, err := parseTriplePattern(pattern)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing pattern %q: %s", pattern, err)) // TODO: return error instead of panic
	}
	var resVar *schema.Value
	if resultVar != "" {
		resVar = parseVariable(resultVar, &schema.Value{}, false)
	}
	return wrapBareQB(b.Bare.Path(
		*parseVariable(subj, &schema.Value{}, false),
		parsedPattern,
		*parseVariable(obj, &schema.Value{}, false),
		resVar,
	))
}

func (b *QueryBuilder) Size(graph string, resultVar Variable) *QueryBuilder {
	return wrapBareQB(b.Bare.Size(
		graph,
		*parseVariable(resultVar, &schema.DataValue{}, false),
	))
}

func wrapBareQB(qb *bare.QueryBuilder) *QueryBuilder {
	return &QueryBuilder{qb}
}
