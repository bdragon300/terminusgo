package woql

import (
	"github.com/bdragon300/terminusgo/woql/bare"
	"github.com/bdragon300/terminusgo/woql/schema"
)

// TODO: implement WOQLLibrary
// TODO: see if it's needed to implement woqlDoc
// TODO: implement Get/Put (require .As() to be implemented)

func NewSimpleQueryBuilder() *QueryBuilder {
	return &QueryBuilder{Bare: bare.NewQueryBuilder(), vocabulary: defaultVocabulary}
}

type (
	Variable          string
	StringOrVariable  string
	IntegerOrVariable any
	NumberOrVariable  any
	AnyOrVariable     any
)

// intOrVarWrapper is a wrapper for IntegerOrVariable variable, this helps to distinguish a IntegerOrVariable type
// from other "any" aliases
type intOrVarWrapper struct {
	v any
}

// numOrVarWrapper is a wrapper for NumberOrVariable variable, this helps to distinguish a NumberOrVariable type
// from other "any" aliases
type numOrVarWrapper struct {
	v any
}

type QueryBuilder struct {
	Bare       *bare.QueryBuilder
	vocabulary map[string]string
}

func (b *QueryBuilder) WithVocabulary(vocab map[string]string) *QueryBuilder {
	b.vocabulary = vocab
	return b
}

func (b *QueryBuilder) ToRaw(buf *bare.RawQuery) error {
	return b.Bare.ToRaw(buf)
}

func (b *QueryBuilder) GetQueryData() any {
	return b.Bare.GetQueryData()
}

func (b *QueryBuilder) Clone() *QueryBuilder {
	return b.clone(b.Bare.Clone())
}

func (b *QueryBuilder) And(subQueries ...schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.And(subQueries...))
}

func (b *QueryBuilder) Not(query schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Not(query))
}

func (b *QueryBuilder) Select(vars ...Variable) *QueryBuilder {
	params := make([]string, 0)
	for _, v := range vars {
		if name, err := extractVariableName(v); err == nil {
			params = append(params, name)
		} else {
			panic(err)
		}
	}
	return b.clone(b.Bare.Select(params))
}

func (b *QueryBuilder) From(graph string, subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.From(graph, subQuery))
}

func (b *QueryBuilder) Using(collection string, subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Using(collection, subQuery))
}

func (b *QueryBuilder) Distinct(vars ...Variable) *QueryBuilder {
	params := make([]string, 0)
	for _, v := range vars {
		if name, err := extractVariableName(v); err == nil {
			params = append(params, name)
		} else {
			panic(err)
		}
	}
	return b.clone(b.Bare.Distinct(params))
}

func (b *QueryBuilder) Into(graph string, subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Into(graph, subQuery))
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
			panic(err)
		}
	}
	return b.clone(b.Bare.OrderBy(ordering))
}

func (b *QueryBuilder) GroupBy(groupVars, templateVars []Variable, outputVar Variable, subQuery schema.Querier) *QueryBuilder {
	grpVars := make([]string, 0)
	for _, v := range groupVars {
		varName, err := extractVariableName(v)
		if err != nil {
			panic(err)
		}
		grpVars = append(grpVars, varName)
	}
	tplVars := make([]schema.Value, 0)
	for _, v := range templateVars {
		tplVars = append(tplVars, *parseVariable(v, &schema.Value{}, false))
	}
	return b.clone(b.Bare.GroupBy(
		grpVars,
		schema.Value{List: tplVars},
		*parseVariable(outputVar, &schema.Value{}, false),
		subQuery,
	))
}

func (b *QueryBuilder) TripleCount(resourceID string, countVar IntegerOrVariable) *QueryBuilder {
	return b.clone(b.Bare.TripleCount(
		resourceID,
		*parseVariable(intOrVarWrapper{countVar}, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Triple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Triple(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) AddTriple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.AddTriple(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) AddedTriple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.AddedTriple(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) DeleteTriple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.DeleteTriple(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) RemovedTriple(subject, predicate, object StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.RemovedTriple(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) UpdateTriple(subject, predicate, newObject StringOrVariable) *QueryBuilder {
	return b.And(
		b.Query().Optional(
			b.Query().
				Triple(subject, predicate, "v:AnyObject").
				DeleteTriple(subject, predicate, "v:AnyObject").
				Not(nil).
				Triple(subject, predicate, newObject),
		),
		b.Query().AddTriple(subject, predicate, newObject),
	)
}

func (b *QueryBuilder) Quad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return b.clone(b.Bare.Quad(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) AddQuad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return b.clone(b.Bare.AddQuad(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph, // TODO: maybe it's needed to transform graph id? See cleanGraph comment in js client
	))
}

func (b *QueryBuilder) AddedQuad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return b.clone(b.Bare.AddedQuad(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) DeleteQuad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return b.clone(b.Bare.DeleteQuad(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) RemovedQuad(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return b.clone(b.Bare.RemovedQuad(
		*parseVariable(fromVocab(b, subject), &schema.NodeValue{}, false),
		*parseVariable(fromVocab(b, predicate), &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) UpdateQuad(subject, predicate, newObject StringOrVariable, graph string) *QueryBuilder {
	return b.And(
		b.Query().Optional(
			b.Query().
				Quad(subject, predicate, "v:AnyObject", graph).
				DeleteQuad(subject, predicate, "v:AnyObject", graph).
				Not(nil).
				Quad(subject, predicate, newObject, graph),
		),
		b.Query().AddQuad(subject, predicate, newObject, graph),
	)
}

func (b *QueryBuilder) Subsumption(parent, child StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Subsumption(
		*parseVariable(parent, &schema.NodeValue{}, false),
		*parseVariable(child, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) Equals(left, right AnyOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Equals(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) SubString(string StringOrVariable, length IntegerOrVariable, substring StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.SubString(
		*parseVariable(string, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{length}, &schema.DataValue{}, true),
		*parseVariable(substring, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{uint(0)}, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{uint(0)}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) SubStringBeforeAfter(string StringOrVariable, before, length, after IntegerOrVariable, substring StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.SubString(
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
	return b.clone(b.Bare.UpdateDocument(
		*parseVariable(document, &schema.Value{}, false),
		param,
	))
}

func (b *QueryBuilder) InsertDocument(document map[string]any, iri StringOrVariable) *QueryBuilder {
	var param *schema.NodeValue
	if iri != "" {
		param = parseVariable(iri, &schema.NodeValue{}, false)
	}
	return b.clone(b.Bare.InsertDocument(
		*parseVariable(document, &schema.Value{}, false),
		param,
	))
}

func (b *QueryBuilder) DeleteDocument(iri StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.DeleteDocument(
		*parseVariable(iri, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) ReadDocument(iri StringOrVariable, outputVar Variable) *QueryBuilder {
	return b.clone(b.Bare.ReadDocument(
		*parseVariable(iri, &schema.NodeValue{}, false),
		*parseVariable(outputVar, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) Once(subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Once(subQuery))
}

func (b *QueryBuilder) Trim(untrimmed, trimmed StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Trim(
		*parseVariable(untrimmed, &schema.DataValue{}, true),
		*parseVariable(trimmed, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Upper(left StringOrVariable, right Variable) *QueryBuilder {
	return b.clone(b.Bare.Upper(
		*parseVariable(left, &schema.DataValue{}, true),
		*parseVariable(right, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Lower(left StringOrVariable, right Variable) *QueryBuilder {
	return b.clone(b.Bare.Lower(
		*parseVariable(left, &schema.DataValue{}, true),
		*parseVariable(right, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) ConcatenateList(stringsOrVars []StringOrVariable, result StringOrVariable) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range stringsOrVars {
		params = append(params, *parseVariable(v, &schema.DataValue{}, true))
	}
	return b.clone(b.Bare.Concatenate(
		schema.DataValue{List: params},
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Concatenate(listVar Variable, result StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Concatenate(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) JoinList(stringsOrVars []StringOrVariable, separator, result StringOrVariable) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range stringsOrVars {
		params = append(params, *parseVariable(v, &schema.DataValue{}, true))
	}
	return b.clone(b.Bare.Join(
		schema.DataValue{List: params},
		*parseVariable(separator, &schema.DataValue{}, true),
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Join(listVar Variable, separator, result StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Join(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(separator, &schema.DataValue{}, true),
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Split(varName, separator, resultVar StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Split(
		*parseVariable(varName, &schema.DataValue{}, true),
		*parseVariable(separator, &schema.DataValue{}, true),
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Regexp(pattern, str StringOrVariable, result Variable) *QueryBuilder {
	return b.clone(b.Bare.Regexp(
		*parseVariable(pattern, &schema.DataValue{}, true),
		*parseVariable(str, &schema.DataValue{}, true),
		*parseVariable(result, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) RegexpToList(pattern, str StringOrVariable, result []StringOrVariable) *QueryBuilder {
	resParam := schema.DataValue{}
	if len(result) > 0 {
		for _, v := range result {
			resParam.List = append(resParam.List, *parseVariable(v, &schema.DataValue{}, true))
		}
	}
	return b.clone(b.Bare.Regexp(
		*parseVariable(pattern, &schema.DataValue{}, true),
		*parseVariable(str, &schema.DataValue{}, true),
		resParam,
	))
}

func (b *QueryBuilder) Eval(arith schema.ArithmeticExpressionType, result Variable) *QueryBuilder {
	return b.clone(b.Bare.Eval(
		arith,
		*parseVariable(result, &schema.ArithmeticValue{}, true),
	))
}

func (b *QueryBuilder) IsA(element, typ StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.IsA(
		*parseVariable(element, &schema.NodeValue{}, false),
		*parseVariable(typ, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) Like(left, right StringOrVariable, similarity NumberOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Like(
		*parseVariable(left, &schema.DataValue{}, true),
		*parseVariable(right, &schema.DataValue{}, true),
		*parseVariable(numOrVarWrapper{similarity}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Less(left, right NumberOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Less(
		*parseVariable(numOrVarWrapper{left}, &schema.DataValue{}, false),
		*parseVariable(numOrVarWrapper{right}, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Greater(left, right NumberOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Greater(
		*parseVariable(numOrVarWrapper{left}, &schema.DataValue{}, false),
		*parseVariable(numOrVarWrapper{right}, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Optional(subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Optional(subQuery))
}

func (b *QueryBuilder) Unique(prefix string, keyList []StringOrVariable, resultVar Variable) *QueryBuilder {
	var keys []schema.DataValue
	for _, v := range keyList {
		keys = append(keys, *parseVariable(v, &schema.DataValue{}, true))
	}
	return b.clone(b.Bare.HashKey(
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
	return b.clone(b.Bare.LexicalKey(
		*parseString(prefix, &schema.DataValue{}, true),
		keys,
		*parseVariable(resultVar, &schema.NodeValue{}, true),
	))
}

func (b *QueryBuilder) Pad(input, pad StringOrVariable, length IntegerOrVariable, resultVar Variable) *QueryBuilder {
	return b.clone(b.Bare.Pad(
		*parseVariable(input, &schema.DataValue{}, true),
		*parseVariable(pad, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{length}, &schema.DataValue{}, true),
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Dot(document, field, value StringOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Dot(
		*parseVariable(document, &schema.DataValue{}, false),
		*parseVariable(field, &schema.DataValue{}, true),
		*parseVariable(value, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Member(member, list AnyOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Member(
		*parseVariable(member, &schema.DataValue{}, false),
		*parseVariable(list, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) MemberOfList(member AnyOrVariable, list []AnyOrVariable) *QueryBuilder {
	var paramList []schema.DataValue
	for _, v := range list {
		paramList = append(paramList, *parseVariable(v, &schema.DataValue{}, false))
	}
	return b.clone(b.Bare.Member(
		*parseVariable(member, &schema.DataValue{}, false),
		schema.DataValue{List: paramList},
	))
}

func (b *QueryBuilder) SumList(list []NumberOrVariable, outputVar Variable) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range list {
		params = append(params, *parseVariable(numOrVarWrapper{v}, &schema.DataValue{}, true))
	}
	return b.clone(b.Bare.Sum(
		schema.DataValue{List: params},
		*parseVariable(outputVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Sum(listVar, outputVar Variable) *QueryBuilder {
	return b.clone(b.Bare.Sum(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(outputVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Start(start uint, subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Start(start, subQuery))
}

func (b *QueryBuilder) Limit(limit uint, subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Limit(limit, subQuery))
}

func (b *QueryBuilder) Length(listVar Variable, resultVar IntegerOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Length(
		*parseVariable(listVar, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{resultVar}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) LengthList(list []AnyOrVariable, resultVar IntegerOrVariable) *QueryBuilder {
	return b.clone(b.Bare.Length(
		*parseVariable(list, &schema.DataValue{}, true),
		*parseVariable(intOrVarWrapper{resultVar}, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Immediately(subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Immediately(subQuery))
}

func (b *QueryBuilder) Count(countVar IntegerOrVariable, subQuery schema.Querier) *QueryBuilder {
	return b.clone(b.Bare.Count(*parseVariable(intOrVarWrapper{countVar}, &schema.DataValue{}, false), subQuery))
}

func (b *QueryBuilder) TypeCast(value AnyOrVariable, typ StringOrVariable, resultVar Variable) *QueryBuilder {
	return b.clone(b.Bare.TypeCast(
		*parseVariable(value, &schema.Value{}, false),
		*parseVariable(typ, &schema.NodeValue{}, false),
		*parseVariable(resultVar, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) TypeOf(value, typ string) *QueryBuilder {
	return b.clone(b.Bare.TypeOf(
		*parseString(value, &schema.Value{}, true),
		*parseString(typ, &schema.NodeValue{}, true),
	))
}

func (b *QueryBuilder) True() *QueryBuilder {
	return b.clone(b.Bare.True())
}

func (b *QueryBuilder) Path(subj StringOrVariable, pattern string, obj StringOrVariable, resultVar Variable) *QueryBuilder {
	parsedPattern, err := parseTriplePattern(pattern) // TODO: add Validate method
	if err != nil {
		panic(err)
	}
	var resVar *schema.Value
	if resultVar != "" {
		resVar = parseVariable(resultVar, &schema.Value{}, false)
	}
	return b.clone(b.Bare.Path(
		*parseVariable(subj, &schema.Value{}, false),
		parsedPattern,
		*parseVariable(obj, &schema.Value{}, false),
		resVar,
	))
}

func (b *QueryBuilder) Size(graph string, resultVar Variable) *QueryBuilder {
	return b.clone(b.Bare.Size(
		graph,
		*parseVariable(resultVar, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Star(graph string, subject, predicate, object StringOrVariable) *QueryBuilder {
	if subject == "" {
		subject = "v:Subject"
	}
	if predicate == "" {
		predicate = "v:Predicate"
	}
	if object == "" {
		object = "v:Object"
	}
	if graph != "" {
		return b.Quad(subject, predicate, object, graph)
	}
	return b.Triple(subject, predicate, object)
}

func (b *QueryBuilder) All(subject, predicate, object StringOrVariable, graph string) *QueryBuilder {
	return b.Star(graph, subject, predicate, object)
}

func (b *QueryBuilder) Insert(id, typ StringOrVariable, graph string) *QueryBuilder {
	if graph != "" {
		return b.AddQuad(id, "rdf:type", "@schema:"+typ, graph)
	}
	return b.AddTriple(id, "rdf:type", "@schema:"+typ)
}

func (b *QueryBuilder) Nuke(graph string) *QueryBuilder {
	if graph != "" {
		return b.Quad("v:A", "v:B", "v:C", graph).DeleteQuad("v:A", "v:B", "v:C", graph)
	}
	return b.Triple("v:A", "v:B", "v:C").DeleteTriple("v:A", "v:B", "v:C")
}

func (b *QueryBuilder) Query() *QueryBuilder {
	return NewSimpleQueryBuilder()
}

func (b *QueryBuilder) clone(bareQB *bare.QueryBuilder) *QueryBuilder {
	return &QueryBuilder{Bare: bareQB, vocabulary: b.vocabulary}
}
