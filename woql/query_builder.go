package woql

import (
	"fmt"
	"strings"

	"github.com/bdragon300/terminusgo/woql/bare"
	"github.com/bdragon300/terminusgo/woql/schema"
)

// TODO: check query builder against javascript client and terminus query schema
// TODO: in some methods the parameters can actually be a one var, var list or string. Do smth with this

func NewSimpleQueryBuilder() *QueryBuilder {
	return &QueryBuilder{bare.NewQueryBuilder()}
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

func (b *QueryBuilder) And() *QueryBuilder {
	return wrapBareQB(b.Bare.And())
}

func (b *QueryBuilder) Not() *QueryBuilder {
	return wrapBareQB(b.Bare.Not())
}

func (b *QueryBuilder) Select(vars ...string) *QueryBuilder {
	params := make([]schema.NodeValue, 0)
	for _, v := range vars {
		params = append(params, *parseVariable(v, &schema.NodeValue{}, true))
	}
	return wrapBareQB(b.Bare.Select(params...))
}

func (b *QueryBuilder) From(graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.From(graph))
}

func (b *QueryBuilder) Using(collection string) *QueryBuilder {
	return wrapBareQB(b.Bare.Using(collection))
}

func (b *QueryBuilder) Distinct(vars ...string) *QueryBuilder {
	params := make([]schema.NodeValue, 0)
	for _, v := range vars {
		params = append(params, *parseVariable(v, &schema.NodeValue{}, true))
	}
	return wrapBareQB(b.Bare.Distinct(params...))
}

func (b *QueryBuilder) Into(graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.Into(graph))
}

func (b *QueryBuilder) OrderBy(vars map[string]schema.OrderDirection) *QueryBuilder {
	ordering := make([]schema.OrderTemplate, 0)
	for k, v := range vars {
		ordering = append(ordering, schema.OrderTemplate{
			Order:    v,
			Variable: extractVariableName(k),
		})
	}
	return wrapBareQB(b.Bare.OrderBy(ordering))
}

func (b *QueryBuilder) GroupBy(groupVars []string, templateVar string, outputVar string) *QueryBuilder {
	grpVars := make([]schema.NodeValue, 0)
	for _, v := range groupVars {
		grpVars = append(grpVars, *parseVariable(v, &schema.NodeValue{}, true))
	}
	return wrapBareQB(b.Bare.GroupBy(
		grpVars,
		*parseVariable(templateVar, &schema.Value{}, true),
		*parseVariable(outputVar, &schema.Value{}, true),
	))
}

func (b *QueryBuilder) TripleCount(graph, countVar string) *QueryBuilder {
	return wrapBareQB(b.Bare.TripleCount(
		graph,
		*parseVariable(countVar, &schema.DataValue{}, true), // FIXME: countVar in python client has type Value, figure out why
	))
}

func (b *QueryBuilder) Triple(subject, predicate, object string) *QueryBuilder {
	return wrapBareQB(b.Bare.Triple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) AddTriple(subject, predicate, object string) *QueryBuilder {
	return wrapBareQB(b.Bare.AddTriple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) DeleteTriple(subject, predicate, object string) *QueryBuilder {
	return wrapBareQB(b.Bare.DeleteTriple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) AddedTriple(subject, predicate, object string) *QueryBuilder {
	return wrapBareQB(b.Bare.AddedTriple(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
	))
}

func (b *QueryBuilder) Quad(subject, predicate, object, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.Quad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) AddQuad(subject, predicate, object, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.AddQuad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) DeleteQuad(subject, predicate, object, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.DeleteQuad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) AddedQuad(subject, predicate, object, graph string) *QueryBuilder {
	return wrapBareQB(b.Bare.AddedQuad(
		*parseVariable(subject, &schema.NodeValue{}, false),
		*parseVariable(predicate, &schema.NodeValue{}, false),
		*parseVariable(object, &schema.Value{}, false),
		graph,
	))
}

func (b *QueryBuilder) Subsumption(parent, child string) *QueryBuilder {
	return wrapBareQB(b.Bare.Subsumption(
		*parseVariable(parent, &schema.NodeValue{}, false),
		*parseVariable(child, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) Equals(left, right string) *QueryBuilder {
	return wrapBareQB(b.Bare.Equals(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) SubString(string string, length uint, substring string) *QueryBuilder {
	return wrapBareQB(b.Bare.SubString(
		*parseVariable(string, &schema.DataValue{}, false),
		*parseNumber(length, &schema.DataValue{}),
		*parseVariable(substring, &schema.DataValue{}, false),
		*parseNumber(uint(0), &schema.DataValue{}),
		*parseNumber(uint(0), &schema.DataValue{}),
	))
}

func (b *QueryBuilder) SubStringBeforeAfter(string string, length uint, substring string, before, after uint) *QueryBuilder {
	return wrapBareQB(b.Bare.SubString(
		*parseVariable(string, &schema.DataValue{}, false),
		*parseNumber(length, &schema.DataValue{}),
		*parseVariable(substring, &schema.DataValue{}, false),
		*parseNumber(before, &schema.DataValue{}),
		*parseNumber(after, &schema.DataValue{}),
	))
}

func (b *QueryBuilder) UpdateDocument(document map[string]any, iri string) *QueryBuilder {
	// TODO: document is any object, but schema dictates to use Value, where there is not place for object. Figure out what can we do with this
	return nil
}

func (b *QueryBuilder) InsertDocument(document map[string]any, iri string) *QueryBuilder {
	// TODO: document is any object, but schema dictates to use Value, where there is not place for object. Figure out what can we do with this
	return nil
}

func (b *QueryBuilder) DeleteDocument(iri string) *QueryBuilder {
	return wrapBareQB(b.Bare.DeleteDocument(
		*parseVariable(iri, &schema.NodeValue{}, false),
	))
}

func (b *QueryBuilder) ReadDocument(iri, outputVar string) *QueryBuilder {
	return wrapBareQB(b.Bare.ReadDocument(
		*parseVariable(iri, &schema.NodeValue{}, false),
		*parseVariable(outputVar, &schema.Value{}, true),
	))
}

func (b *QueryBuilder) File(uri string, options *schema.QueryResourceOptions) *QueryBuilder {
	source := schema.Source{URI: uri} // FIXME: in python client the "source" key is used, which is not matched to the schema
	return wrapBareQB(b.Bare.File(source, options))
}

func (b *QueryBuilder) Once() *QueryBuilder {
	return wrapBareQB(b.Bare.Once())
}

func (b *QueryBuilder) Remote(uri string, options *schema.QueryResourceOptions) *QueryBuilder {
	source := schema.Source{URI: uri} // FIXME: in python client the "source" key is used, which is not matched to the schema
	return wrapBareQB(b.Bare.Remote(source, options))
}

func (b *QueryBuilder) Post(uri string, options *schema.QueryResourceOptions) *QueryBuilder {
	source := schema.Source{URI: uri} // FIXME: in python client the "source" key is used, which is not matched to the schema
	return wrapBareQB(b.Bare.Post(source, options))
}

func (b *QueryBuilder) Trim(untrimmed, trimmed string) *QueryBuilder {
	return wrapBareQB(b.Bare.Trim(
		*parseVariable(untrimmed, &schema.DataValue{}, false),
		*parseVariable(trimmed, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Upper(left, right string) *QueryBuilder {
	return wrapBareQB(b.Bare.Upper(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Lower(left, right string) *QueryBuilder {
	return wrapBareQB(b.Bare.Lower(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
	))
}

// TODO: maybe it's needed a version with var list and one var, not only with DataValue.List
func (b *QueryBuilder) Concatenate(varExpr, result string) *QueryBuilder {
	var params []schema.DataValue

	// Split up a string to variables and the regular text blocks and push to param list with appropriate type
	// E.g "  lorem ipsum v:var1v:var2 dolor sit      v:  var3 amet " transforms to
	// "  lorem ipsum" + "v:var1" + "v:var2" + " dolor sit      " + "v:var3" + " amet "
	parts := strings.Split(varExpr, "v:")
	for _, part := range parts {
		fields := strings.Fields(part)
		if len(fields) < 1 {
			panic("unnamed variable in concatenate expression") // TODO: return error instead of panic
		}
		params = append(params, *parseVariable("v:"+fields[0], &schema.DataValue{}, true))

		part = part[strings.Index(part, fields[0])+len(fields[0]):]
		if part != "" {
			params = append(params, *parseVariable(part, &schema.DataValue{}, false))
		}
	}
	return wrapBareQB(b.Bare.Concatenate(
		schema.DataValue{List: params},
		*parseVariable(result, &schema.DataValue{}, false),
	))
}

// TODO: maybe it's needed a version with var list and one var, not only with DataValue.List
func (b *QueryBuilder) Join(varList []string, separator, result string) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range varList {
		params = append(params, *parseVariable(v, &schema.DataValue{}, false))
	}
	return wrapBareQB(b.Bare.Join(
		schema.DataValue{List: params},
		*parseVariable(separator, &schema.DataValue{}, false),
		*parseVariable(result, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Split(varName, separator, resultVar string) *QueryBuilder {
	return wrapBareQB(b.Bare.Split(
		*parseVariable(varName, &schema.DataValue{}, false),
		*parseVariable(separator, &schema.DataValue{}, false),
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Regexp(pattern, str string, result string) *QueryBuilder {
	var resParam *schema.DataValue
	if result != "" {
		resParam = parseVariable(result, &schema.DataValue{}, false)
	}
	return wrapBareQB(b.Bare.Regexp(
		*parseVariable(pattern, &schema.DataValue{}, false),
		*parseVariable(str, &schema.DataValue{}, false),
		resParam,
	))
}

func (b *QueryBuilder) Eval(arith schema.ArithmeticExpressionType, result string) *QueryBuilder {
	return wrapBareQB(b.Bare.Eval(
		arith,
		*parseVariable(result, &schema.ArithmeticValue{}, false),
	))
}

func (b *QueryBuilder) IsA(element, typ string) *QueryBuilder {
	return wrapBareQB(b.Bare.IsA(
		*parseVariable(element, &schema.NodeValue{}, false),
		*parseVariable(typ, &schema.NodeValue{}, false),
	))
}

// TODO: make strict type for similarity -- it can be a variable or number
func (b *QueryBuilder) Like(left, right string, similarity any) *QueryBuilder {
	return wrapBareQB(b.Bare.Like(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
		*parseVariable(similarity, &schema.DataValue{}, false),
	))
}

// TODO: make strict type for params -- they can be a variable or number
func (b *QueryBuilder) Less(left, right any) *QueryBuilder {
	return wrapBareQB(b.Bare.Less(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
	))
}

// TODO: make strict type for params -- they can be a variable or number
func (b *QueryBuilder) Greater(left, right any) *QueryBuilder {
	return wrapBareQB(b.Bare.Greater(
		*parseVariable(left, &schema.DataValue{}, false),
		*parseVariable(right, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Optional() *QueryBuilder {
	return wrapBareQB(b.Bare.Optional())
}

func (b *QueryBuilder) Unique(prefix string, keyList []string, resultVar string) *QueryBuilder {
	var keys []schema.DataValue
	for _, v := range keyList {
		keys = append(keys, *parseVariable(v, &schema.DataValue{}, false))
	}
	return wrapBareQB(b.Bare.HashKey(
		*parseVariable(prefix, &schema.DataValue{}, false),
		keys,
		*parseVariable(resultVar, &schema.NodeValue{}, true),
	))
}

func (b *QueryBuilder) IDGen(prefix string, keyList []string, resultVar string) *QueryBuilder {
	var keys []schema.DataValue
	for _, v := range keyList {
		keys = append(keys, *parseVariable(v, &schema.DataValue{}, false))
	}
	return wrapBareQB(b.Bare.LexicalKey(
		*parseVariable(prefix, &schema.DataValue{}, false),
		keys,
		*parseVariable(resultVar, &schema.NodeValue{}, true),
	))
}

// TODO: in schema there are two parameters, in python client there are three -- figure out why
func (b *QueryBuilder) RandomIDGen(prefix string, resultVar string) *QueryBuilder {
	return wrapBareQB(b.Bare.RandomKey(
		*parseVariable(prefix, &schema.DataValue{}, false),
		*parseVariable(resultVar, &schema.NodeValue{}, true),
	))
}

// TODO: make strict type for length -- they can be a variable or number
func (b *QueryBuilder) Pad(input, pad string, length any, resultVar string) *QueryBuilder {
	return wrapBareQB(b.Bare.Pad(
		*parseVariable(input, &schema.DataValue{}, false),
		*parseVariable(pad, &schema.DataValue{}, false),
		*parseVariable(length, &schema.DataValue{}, false),
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Dot(document, field, value string) *QueryBuilder {
	return wrapBareQB(b.Bare.Dot(
		*parseVariable(document, &schema.DataValue{}, false),
		*parseVariable(field, &schema.DataValue{}, false),
		*parseVariable(value, &schema.DataValue{}, false),
	))
}

// TODO: list can be either a variable or a list, see docs
func (b *QueryBuilder) Member(member any, list string) *QueryBuilder {
	return wrapBareQB(b.Bare.Member(
		*parseVariable(member, &schema.DataValue{}, false),
		*parseVariable(list, &schema.DataValue{}, false),
	))
}

func (b *QueryBuilder) Sum(list []any, outputVar string) *QueryBuilder {
	var params []schema.DataValue
	for _, v := range list {
		params = append(params, *parseVariable(v, &schema.DataValue{}, false))
	}
	return wrapBareQB(b.Bare.Sum(
		schema.DataValue{List: params},
		*parseVariable(outputVar, &schema.DataValue{}, true),
	))
}

// TODO: in docs/python client start parameter is a number or variable
func (b *QueryBuilder) Start(start uint) *QueryBuilder {
	return wrapBareQB(b.Bare.Start(start))
}

// TODO: in docs/python client start parameter is a number or variable
func (b *QueryBuilder) Limit(limit uint) *QueryBuilder {
	return wrapBareQB(b.Bare.Limit(limit))
}

func (b *QueryBuilder) Length(varList, resultVar string) *QueryBuilder {
	return wrapBareQB(b.Bare.Length(
		*parseVariable(varList, &schema.DataValue{}, false),
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func (b *QueryBuilder) Immediately() *QueryBuilder {
	return wrapBareQB(b.Bare.Immediately())
}

// TODO: countVar can be either a number or variable -- make its type more strict
func (b *QueryBuilder) Count(countVar any) *QueryBuilder {
	return wrapBareQB(b.Bare.Count(*parseVariable(countVar, &schema.DataValue{}, false)))
}

func (b *QueryBuilder) TypeCast(value any, typ, resultVar string) *QueryBuilder {
	if !schema.ValidateLiteralType(typ) {
		panic(fmt.Sprintf("Invalid literal type %q", typ)) // TODO: error instead of panic
	}
	typParam := &schema.NodeValue{}
	typParam.FromString(typ)
	return wrapBareQB(b.Bare.TypeCast(
		*parseVariable(value, &schema.Value{}, false),
		*typParam,
		*parseVariable(resultVar, &schema.Value{}, false),
	))
}

// TODO: maybe it's needed to validate `typ` somehow?
func (b *QueryBuilder) TypeOF(value, typ string) *QueryBuilder {
	typParam := &schema.NodeValue{}
	typParam.FromString(typ)
	return wrapBareQB(b.Bare.TypeOf(
		*parseVariable(value, &schema.Value{}, false),
		*typParam,
	))
}

func (b *QueryBuilder) True() *QueryBuilder {
	return wrapBareQB(b.Bare.True())
}

func (b *QueryBuilder) Path(subj, pattern, obj, resultVar string) *QueryBuilder {
	p, err := parseTriplePattern(pattern)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing pattern %q: %s", pattern, err))
	}
	var resVar *schema.Value
	if resultVar != "" {
		resVar = parseVariable(resultVar, &schema.Value{}, true)
	}
	return wrapBareQB(b.Bare.Path(
		*parseVariable(subj, &schema.Value{}, false),
		p,
		*parseVariable(obj, &schema.Value{}, false),
		resVar,
	))
}

func (b *QueryBuilder) Size(graph, resultVar string) *QueryBuilder {
	return wrapBareQB(b.Bare.Size(
		graph,
		*parseVariable(resultVar, &schema.DataValue{}, true),
	))
}

func wrapBareQB(qb *bare.QueryBuilder) *QueryBuilder {
	return &QueryBuilder{qb}
}
