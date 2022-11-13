package bare

import (
	"github.com/bdragon300/terminusgo/woql/schema"
)

// TODO: maybe it's needed a package for JSON-LD

type AggFunctionKind uint

const (
	AggAnd AggFunctionKind = iota
	AggOr
)

type AggFunction struct {
	Kind  AggFunctionKind
	Items []schema.Querier
}

func (a AggFunction) GetQuery() schema.Querier {
	if len(a.Items) > 1 {
		switch a.Kind {
		// TODO: flatten nested Ands in one And, nested Ors in one Or (using appropriate functions)
		case AggAnd:
			return schema.And{And: a.Items}
		case AggOr:
			return schema.Or{Or: a.Items}
		}
	} else if len(a.Items) == 1 {
		return a.Items[0]
	}
	return nil
}

func NewBareQueryBuilder() *QueryBuilder {
	aggFunction := AggFunction{}
	return &QueryBuilder{
		root:   aggFunction,
		cursor: &aggFunction,
		errors: nil,
	}
}

type QueryBuilder struct {
	root   AggFunction
	cursor *AggFunction
	errors error
}

func (b *QueryBuilder) GetQuery() schema.Querier {
	return b.root
}

func (b *QueryBuilder) And() *QueryBuilder {
	if b.cursor.Kind == AggAnd {
		return b
	}
	b.cursor = &AggFunction{
		Kind:  AggAnd,
		Items: []schema.Querier{*b.cursor},
	}
	return b
}

func (b *QueryBuilder) Or() *QueryBuilder {
	if b.cursor.Kind == AggOr {
		return b
	}
	b.cursor = &AggFunction{
		Kind:  AggOr,
		Items: []schema.Querier{*b.cursor},
	}
	return b
}

func (b *QueryBuilder) Not() *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Not{SubQuery: agg})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) Select(variables ...schema.NodeValue) *QueryBuilder {
	agg := AggFunction{}
	var varNames []string
	for _, varName := range variables {
		varNames = append(varNames, varName.Variable)
	}
	b.cursor.Items = append(b.cursor.Items, schema.Select{
		SubQuery:  agg,
		Variables: varNames,
	})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) From(graph string) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.From{SubQuery: agg, Graph: graph})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) Using(collection string) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Using{SubQuery: agg, Collection: collection})
	b.cursor = &agg
	return b
}

// TODO: Implement Comment function

func (b *QueryBuilder) Distinct(variables ...schema.NodeValue) *QueryBuilder {
	agg := AggFunction{}
	var varNames []string
	for _, varName := range variables {
		varNames = append(varNames, varName.Variable)
	}
	b.cursor.Items = append(b.cursor.Items, schema.Distinct{
		SubQuery:  agg,
		Variables: varNames,
	})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) Into(graph string) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Into{SubQuery: agg, Graph: graph})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) OrderBy(ordering []schema.OrderTemplate) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.OrderBy{
		SubQuery: agg,
		Ordering: ordering,
	})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) GroupBy(template schema.Value, groupBy []string, grouped schema.Value) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.GroupBy{
		Template: template,
		GroupBy:  groupBy,
		SubQuery: agg,
		Grouped:  grouped,
	})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) TripleCount(resource string, count schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.TripleCount{
		Resource: resource,
		Count:    count,
	})
	return b
}

func (b *QueryBuilder) Triple(subj, pred schema.NodeValue, obj schema.Value) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Triple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
	})
	return b
}

func (b *QueryBuilder) AddTriple(subject, predicate schema.NodeValue, object schema.Value) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.AddTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	})
	return b
}

func (b *QueryBuilder) DeleteTriple(subject, predicate schema.NodeValue, object schema.Value) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.DeleteTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	})
	return b
}

func (b *QueryBuilder) AddedTriple(subj, pred schema.NodeValue, obj schema.Value) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.AddedTriple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
	})
	return b
}

// TODO: implement RemovedTriple

func (b *QueryBuilder) Quad(subj, pred schema.NodeValue, obj schema.Value, graph string) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Triple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
		Graph:     &graph,
	})
	return b
}

func (b *QueryBuilder) AddQuad(subject, predicate schema.NodeValue, object schema.Value, graph string) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.AddTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Graph:     &graph,
	})
	return b
}

func (b *QueryBuilder) DeleteQuad(subject, predicate schema.NodeValue, object schema.Value, graph string) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.DeleteTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Graph:     &graph,
	})
	return b
}

func (b *QueryBuilder) AddedQuad(subj, pred schema.NodeValue, obj schema.Value, graph string) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.AddedTriple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
		Graph:     &graph,
	})
	return b
}

// TODO: implement RemovedQuad

func (b *QueryBuilder) Subsumption(parent, child schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Subsumption{
		Child:  child,
		Parent: parent,
	})
	return b
}

func (b *QueryBuilder) Equals(left, right schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Equals{ // TODO: in python client left\right are schema.Value, figure out why
		Left:  left,
		Right: right,
	})
	return b
}

func (b *QueryBuilder) SubString(string, length, substring, before, after schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.SubString{
		String:    string,
		Length:    length,
		Before:    before,
		After:     after,
		SubString: substring,
	})
	return b
}

func (b *QueryBuilder) UpdateDocument(document schema.Value, identifier *schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.UpdateDocument{
		Document:   document,
		Identifier: identifier,
	})
	return b
}

func (b *QueryBuilder) InsertDocument(document schema.Value, identifier *schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.InsertDocument{
		Document:   document,
		Identifier: identifier,
	})
	return b
}

func (b *QueryBuilder) DeleteDocument(identifier schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.DeleteDocument{
		Identifier: identifier,
	})
	return b
}

func (b *QueryBuilder) ReadDocument(document schema.Value, identifier schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.ReadDocument{
		Document:   document,
		Identifier: identifier,
	})
	return b
}

func (b *QueryBuilder) Get(columns []schema.Column, resource schema.QueryResource) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Get{
		Columns:  columns,
		Resource: resource,
		// HasHeader: nil,  // TODO: HasHeader is not used in python client, figure out why
	}) // FIXME: in python client cursor moves to resource, figure out what it means
	return b
}

// TODO: implement Put

// FIXME: understand if File/Remote/Post are actually needed instead of a single function
func (b *QueryBuilder) File(source schema.Source, options *schema.QueryResourceOptions) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, queryResource(source, options))
	return b
}

func (b *QueryBuilder) Once() *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Once{SubQuery: agg})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) Remote(source schema.Source, options *schema.QueryResourceOptions) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, queryResource(source, options))
	return b
}

func (b *QueryBuilder) Post(source schema.Source, options *schema.QueryResourceOptions) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, queryResource(source, options))
	return b
}

func (b *QueryBuilder) Trim(trimmed, untrimmed schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Trim{
		Untrimmed: untrimmed,
		Trimmed:   trimmed,
	})
	return b
}

func (b *QueryBuilder) Upper(mixed, upper schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Upper{ // TODO: figure out why parameters are called left\right in python client
		Mixed: mixed,
		Upper: upper,
	})
	return b
}

func (b *QueryBuilder) Lower(mixed, lower schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Lower{ // TODO: figure out why parameters are called left\right in python client
		Mixed: mixed,
		Lower: lower,
	})
	return b
}

func (b *QueryBuilder) Concatenate(list, result schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Concatenate{
		List:   list,
		Result: result,
	})
	return b
}

func (b *QueryBuilder) Join(list, separator, result schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Join{
		List:      list,
		Separator: separator,
		Result:    result,
	})
	return b
}

func (b *QueryBuilder) Regexp(pattern, str schema.DataValue, result *schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Regexp{
		Pattern: pattern,
		String:  str,
		Result:  result,
	})
	return b
}

func (b *QueryBuilder) Eval(expression schema.ArithmeticExpression, result schema.ArithmeticValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Eval{
		Expression: expression,
		Result:     result,
	})
	return b
}

func (b *QueryBuilder) IsA(element, typ schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.IsA{
		Element: element,
		Type:    typ,
	})
	return b
}

func (b *QueryBuilder) Like(left, right, similarity schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Like{
		Left:       left,
		Right:      right,
		Similarity: similarity,
	})
	return b
}

func (b *QueryBuilder) Less(left, right schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Less{
		Left:  left,
		Right: right,
	})
	return b
}

func (b *QueryBuilder) Greater(left, right schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Greater{
		Left:  left,
		Right: right,
	})
	return b
}

func (b *QueryBuilder) Optional() *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Optional{SubQuery: agg})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) HashKey(base schema.DataValue, keyList []schema.DataValue, uri schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.HashKey{
		Base:    base,
		KeyList: keyList,
		URI:     uri,
	})
	return b
}

func (b *QueryBuilder) LexicalKey(base schema.DataValue, keyList []schema.DataValue, uri schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.LexicalKey{
		Base:    base,
		KeyList: keyList,
		URI:     uri,
	})
	return b
}

func (b *QueryBuilder) RandomKey(base schema.DataValue, uri schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.RandomKey{
		Base: base,
		URI:  uri,
	})
	return b
}

func (b *QueryBuilder) Pad(str, char, times, result schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Pad{
		String: str,
		Char:   char,
		Times:  times,
		Result: result,
	})
	return b
}

func (b *QueryBuilder) Split(str, pattern, list schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Split{
		String:  str,
		Pattern: pattern,
		List:    list,
	})
	return b
}

func (b *QueryBuilder) Dot(document, field, value schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Dot{
		Document: document,
		Field:    field,
		Value:    value,
	})
	return b
}

func (b *QueryBuilder) Member(member, list schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Member{
		Member: member,
		List:   list,
	})
	return b
}

func (b *QueryBuilder) Sum(list, result schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Sum{
		List:   list,
		Result: result,
	})
	return b
}

func (b *QueryBuilder) Start(start uint) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Start{
		SubQuery: agg,
		Start:    start,
	})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) Limit(limit uint) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Limit{
		SubQuery: agg,
		Limit:    limit,
	})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) Length(list, length schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Length{
		List:   list,
		Length: length, // FIXME: in python client this field has type Value, figure out why
	})
	return b
}

func (b *QueryBuilder) Immediately() *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Immediately{SubQuery: agg})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) Count(count schema.DataValue) *QueryBuilder {
	agg := AggFunction{}
	b.cursor.Items = append(b.cursor.Items, schema.Count{SubQuery: agg, Count: count})
	b.cursor = &agg
	return b
}

func (b *QueryBuilder) TypeCast(value schema.Value, typ schema.NodeValue, result schema.Value) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Typecast{
		Value:  value,
		Type:   typ,
		Result: result,
	})
	return b
}

func (b *QueryBuilder) TypeOf(value schema.Value, typ schema.NodeValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.TypeOf{
		Value: value,
		Type:  typ,
	})
	return b
}

func (b *QueryBuilder) True() *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.True{})
	return b
}

func (b *QueryBuilder) Path(subj schema.Value, pattern schema.PathPatternType, obj schema.Value, path *schema.Value) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Path{
		Subject: subj,
		Pattern: pattern,
		Object:  obj,
		Path:    path,
	})
	return b
}

func (b *QueryBuilder) Size(resource string, size schema.DataValue) *QueryBuilder {
	b.cursor.Items = append(b.cursor.Items, schema.Size{
		Resource: resource,
		Size:     size,
	})
	return b
}

func queryResource(source schema.Source, options *schema.QueryResourceOptions) schema.QueryResource {
	return schema.QueryResource{
		Source:  source,               // FIXME: figure out why we pass also a string here in python client
		Format:  schema.FormatTypeCSV, // FIXME: in python client this field is modified based on options (only in File)
		Options: options,
	}
}
