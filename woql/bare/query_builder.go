package bare

import (
	"github.com/bdragon300/terminusgo/woql/schema"
	"github.com/huandu/go-clone"
)

// TODO: maybe it's needed a package for JSON-LD
// TODO: compare querybuilder interface (including parameter names) with javascript client interface
// TODO: decide if we want accept another Querier parameters in methods, like in python client

type AggFunction struct {
	up    *AggFunction // Pointer to AggFunction on the previous subquery level; self for the topmost object
	Items []schema.Querier
}

func (a AggFunction) GetQuery() schema.Querier {
	if len(a.Items) > 1 {
		// TODO: flatten nested Ands in one And, nested Ors in one Or (using appropriate functions)
		return &schema.And{And: a.Items}
	} else if len(a.Items) == 1 {
		return a.Items[0]
	}
	return &schema.True{} // TODO: figure out if it's ok
}

func NewQueryBuilder() *QueryBuilder {
	aggFunction := AggFunction{}
	aggFunction.up = &aggFunction
	return &QueryBuilder{
		root:   aggFunction,
		cursor: &aggFunction,
	}
}

type QueryBuilder struct {
	root   AggFunction
	cursor *AggFunction
}

func (b *QueryBuilder) GetQuery() schema.Querier {
	return &b.root
}

func (b *QueryBuilder) Clone() *QueryBuilder {
	val := clone.Slowly(b)
	return val.(*QueryBuilder)
}

func (b *QueryBuilder) And() *QueryBuilder {
	bc := b.Clone()
	if len(bc.cursor.Items) == 0 {
		// Do next fluent calls on the previous subquery level
		bc.cursor = bc.cursor.up
	}
	return bc
}

func (b *QueryBuilder) Not() *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Not{SubQuery: &agg})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) Select(variables ...schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	var varNames []string
	for _, varName := range variables {
		varNames = append(varNames, varName.Variable)
	}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Select{
		SubQuery:  agg,
		Variables: varNames,
	})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) From(graph string) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.From{SubQuery: &agg, Graph: graph})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) Using(collection string) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Using{SubQuery: &agg, Collection: collection})
	bc.cursor = &agg
	return bc
}

// TODO: Implement Comment function

func (b *QueryBuilder) Distinct(variables ...schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	var varNames []string
	for _, varName := range variables {
		varNames = append(varNames, varName.Variable)
	}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Distinct{
		SubQuery:  agg,
		Variables: varNames,
	})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) Into(graph string) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Into{SubQuery: &agg, Graph: graph})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) OrderBy(ordering []schema.OrderTemplate) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.OrderBy{
		SubQuery: &agg,
		Ordering: ordering,
	})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) GroupBy(groupBy []schema.NodeValue, template schema.Value, grouped schema.Value) *QueryBuilder {
	var varNames []string
	for _, varName := range groupBy {
		varNames = append(varNames, varName.Variable)
	}
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.GroupBy{
		Template: template, // FIXME: it can be a list in python client, figure out how to interpret template -- list or string
		GroupBy:  varNames,
		SubQuery: &agg,
		Grouped:  grouped,
	})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) TripleCount(resource string, count schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.TripleCount{
		Resource: resource,
		Count:    count,
	})
	return bc
}

func (b *QueryBuilder) Triple(subject, predicate schema.NodeValue, object schema.Value) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Triple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	})
	return bc
}

func (b *QueryBuilder) AddTriple(subject, predicate schema.NodeValue, object schema.Value) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.AddTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	})
	return bc
}

func (b *QueryBuilder) DeleteTriple(subject, predicate schema.NodeValue, object schema.Value) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.DeleteTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	})
	return bc
}

func (b *QueryBuilder) AddedTriple(subj, pred schema.NodeValue, obj schema.Value) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.AddedTriple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
	})
	return bc
}

// TODO: implement RemovedTriple

func (b *QueryBuilder) Quad(subj, pred schema.NodeValue, obj schema.Value, graph string) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Triple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
		Graph:     &graph,
	})
	return bc
}

func (b *QueryBuilder) AddQuad(subject, predicate schema.NodeValue, object schema.Value, graph string) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.AddTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Graph:     &graph,
	})
	return bc
}

func (b *QueryBuilder) DeleteQuad(subject, predicate schema.NodeValue, object schema.Value, graph string) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.DeleteTriple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Graph:     &graph,
	})
	return bc
}

func (b *QueryBuilder) AddedQuad(subj, pred schema.NodeValue, obj schema.Value, graph string) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.AddedTriple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
		Graph:     &graph,
	})
	return bc
}

// TODO: implement RemovedQuad

func (b *QueryBuilder) Subsumption(parent, child schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Subsumption{
		Child:  child,
		Parent: parent,
	})
	return bc
}

func (b *QueryBuilder) Equals(left, right schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Equals{ // TODO: in python client left\right are schema.Value, figure out why
		Left:  left,
		Right: right,
	})
	return bc
}

func (b *QueryBuilder) SubString(string, length, substring, before, after schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.SubString{
		String:    string,
		Length:    length,
		Before:    before,
		After:     after,
		SubString: substring,
	})
	return bc
}

func (b *QueryBuilder) UpdateDocument(document schema.Value, identifier *schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.UpdateDocument{
		Document:   document,
		Identifier: identifier,
	})
	return bc
}

func (b *QueryBuilder) InsertDocument(document schema.Value, identifier *schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.InsertDocument{
		Document:   document,
		Identifier: identifier,
	})
	return bc
}

func (b *QueryBuilder) DeleteDocument(identifier schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.DeleteDocument{
		Identifier: identifier,
	})
	return bc
}

func (b *QueryBuilder) ReadDocument(identifier schema.NodeValue, document schema.Value) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.ReadDocument{
		Document:   document,
		Identifier: identifier,
	})
	return bc
}

// FIXME: maybe it's better to have it as function, not method?
func (b *QueryBuilder) Get(columns []schema.Column, resource schema.QueryResource) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Get{
		Columns:  columns,
		Resource: resource,
		// HasHeader: nil,  // TODO: HasHeader is not used in python client, figure out why
	}) // FIXME: in python client cursor moves to resource, figure out what it means
	return bc
}

// TODO: implement Put

// FIXME: understand if File/Remote/Post are actually needed instead of a single function
func (b *QueryBuilder) File(source schema.Source, options *schema.QueryResourceOptions) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, queryResource(source, options))
	return bc
}

func (b *QueryBuilder) Once() *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Once{SubQuery: &agg})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) Remote(source schema.Source, options *schema.QueryResourceOptions) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, queryResource(source, options))
	return bc
}

func (b *QueryBuilder) Post(source schema.Source, options *schema.QueryResourceOptions) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, queryResource(source, options))
	return bc
}

func (b *QueryBuilder) Trim(untrimmed, trimmed schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Trim{
		Untrimmed: untrimmed,
		Trimmed:   trimmed,
	})
	return bc
}

func (b *QueryBuilder) Upper(mixed, upper schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Upper{ // TODO: figure out why parameters are called left\right in python client
		Mixed: mixed,
		Upper: upper,
	})
	return bc
}

func (b *QueryBuilder) Lower(mixed, lower schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Lower{ // TODO: figure out why parameters are called left\right in python client
		Mixed: mixed,
		Lower: lower,
	})
	return bc
}

func (b *QueryBuilder) Concatenate(list, result schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Concatenate{
		List:   list,
		Result: result,
	})
	return bc
}

func (b *QueryBuilder) Join(list, separator, result schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Join{
		List:      list,
		Separator: separator,
		Result:    result,
	})
	return bc
}

func (b *QueryBuilder) Split(str, pattern, list schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Split{
		String:  str,
		Pattern: pattern,
		List:    list,
	})
	return bc
}

func (b *QueryBuilder) Regexp(pattern, str schema.DataValue, result *schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Regexp{
		Pattern: pattern,
		String:  str,
		Result:  result,
	})
	return bc
}

func (b *QueryBuilder) Eval(expression schema.ArithmeticExpressionType, result schema.ArithmeticValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Eval{
		Expression: expression,
		Result:     result,
	})
	return bc
}

func (b *QueryBuilder) IsA(element, typ schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.IsA{
		Element: element,
		Type:    typ,
	})
	return bc
}

func (b *QueryBuilder) Like(left, right, similarity schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Like{
		Left:       left,
		Right:      right,
		Similarity: similarity,
	})
	return bc
}

func (b *QueryBuilder) Less(left, right schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Less{
		Left:  left,
		Right: right,
	})
	return bc
}

func (b *QueryBuilder) Greater(left, right schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Greater{
		Left:  left,
		Right: right,
	})
	return bc
}

func (b *QueryBuilder) Optional() *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Optional{SubQuery: &agg})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) HashKey(base schema.DataValue, keyList []schema.DataValue, uri schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.HashKey{
		Base:    base,
		KeyList: keyList,
		URI:     uri,
	})
	return bc
}

func (b *QueryBuilder) LexicalKey(base schema.DataValue, keyList []schema.DataValue, uri schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.LexicalKey{
		Base:    base,
		KeyList: keyList,
		URI:     uri,
	})
	return bc
}

func (b *QueryBuilder) RandomKey(base schema.DataValue, uri schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.RandomKey{
		Base: base,
		URI:  uri,
	})
	return bc
}

func (b *QueryBuilder) Pad(str, char, times, result schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Pad{
		String: str,
		Char:   char,
		Times:  times,
		Result: result,
	})
	return bc
}

func (b *QueryBuilder) Dot(document, field, value schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Dot{
		Document: document,
		Field:    field,
		Value:    value,
	})
	return bc
}

func (b *QueryBuilder) Member(member, list schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Member{
		Member: member,
		List:   list,
	})
	return bc
}

func (b *QueryBuilder) Sum(list, result schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Sum{
		List:   list,
		Result: result,
	})
	return bc
}

func (b *QueryBuilder) Start(start uint) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Start{
		SubQuery: &agg,
		Start:    start,
	})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) Limit(limit uint) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Limit{
		SubQuery: &agg,
		Limit:    limit,
	})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) Length(list, length schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Length{
		List:   list,
		Length: length, // FIXME: in python client this field has type Value, figure out why
	})
	return bc
}

func (b *QueryBuilder) Immediately() *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Immediately{SubQuery: &agg})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) Count(count schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	agg := AggFunction{up: bc.cursor}
	bc.cursor.Items = append(bc.cursor.Items, &schema.Count{SubQuery: &agg, Count: count})
	bc.cursor = &agg
	return bc
}

func (b *QueryBuilder) TypeCast(value schema.Value, typ schema.NodeValue, result schema.Value) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Typecast{
		Value:  value,
		Type:   typ,
		Result: result,
	})
	return bc
}

func (b *QueryBuilder) TypeOf(value schema.Value, typ schema.NodeValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.TypeOf{
		Value: value,
		Type:  typ,
	})
	return bc
}

func (b *QueryBuilder) True() *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.True{})
	return bc
}

func (b *QueryBuilder) Path(subj schema.Value, pattern schema.PathPatternType, obj schema.Value, path *schema.Value) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Path{
		Subject: subj,
		Pattern: pattern,
		Object:  obj,
		Path:    path,
	})
	return bc
}

func (b *QueryBuilder) Size(resource string, size schema.DataValue) *QueryBuilder {
	bc := b.Clone()
	bc.cursor.Items = append(bc.cursor.Items, &schema.Size{
		Resource: resource,
		Size:     size,
	})
	return bc
}

func queryResource(source schema.Source, options *schema.QueryResourceOptions) *schema.QueryResource {
	return &schema.QueryResource{
		Source:  source,               // FIXME: figure out why we pass also a string here in python client
		Format:  schema.FormatTypeCSV, // FIXME: in python client this field is modified based on options (only in File)
		Options: options,
	}
}
