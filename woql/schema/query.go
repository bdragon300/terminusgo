package schema

import "github.com/bdragon300/terminusgo/schema"

// TODO: check if all schema structs are actually used

// TODO: add json marshaler/unmarshaler which is not require json tags everywhere (and applies `omitempty` everywhere)

// TODO: not used
type NamedQuery struct {
	Name  string  `json:"name"`
	Query Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

// TODO: not used
type NamedParametricQuery struct {
	Name       string   `json:"name"`
	Parameters []string `json:"parameters"`
	Query      Querier  `json:"query" terminusgo:"type=Class,class=Query"`
}

type Querier interface {
	GetQuery() Querier
}

type Query struct {
	*schema.AbstractModel
	*schema.SubDocumentModel
}

func (q Query) GetQuery() Querier {
	return &q
}

type Path struct {
	Query
	Subject Value           `json:"subject"`
	Pattern PathPatternType `json:"pattern" terminusgo:"type=Class,class=PathPattern"`
	Object  Value           `json:"object"`
	Path    *Value          `json:"path"`
}

type Using struct {
	Query
	SubQuery   Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Collection string  `json:"collection"`
}

type Select struct {
	Query
	SubQuery  Querier  `json:"query" terminusgo:"type=Class,class=Query"`
	Variables []string `json:"variables"`
}

type Distinct struct {
	Query
	SubQuery  Querier  `json:"query" terminusgo:"type=Class,class=Query"`
	Variables []string `json:"variables"`
}

type And struct {
	Query
	And []Querier `json:"and" terminusgo:"type=List,class=Query"`
}

type Or struct {
	Query
	Or []Querier `json:"or" terminusgo:"type=List,class=Query"`
}

type From struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Graph    string  `json:"graph"`
}

type Into struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Graph    string  `json:"graph"`
}

type Triple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     *string   `json:"graph"`
}

type AddTriple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     *string   `json:"graph"`
}

type AddedTriple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     *string   `json:"graph"`
}

type DeleteTriple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type DeletedTriple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type Link struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type Data struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    DataValue `json:"object"`
	Graph     *string   `json:"graph"`
}

type Subsumption struct {
	Query
	Child  NodeValue `json:"child"`
	Parent NodeValue `json:"parent"`
}

type Equals struct {
	Query
	Left  DataValue `json:"left"`
	Right DataValue `json:"right"`
}

type SubString struct {
	Query
	String    DataValue `json:"string"`
	Length    DataValue `json:"length"`
	Before    DataValue `json:"before"`
	After     DataValue `json:"after"`
	SubString DataValue `json:"substring"`
}

type ReadDocument struct {
	Query
	Document   Value     `json:"document"`
	Identifier NodeValue `json:"identifier"`
}

type UpdateDocument struct {
	Query
	Document   Value      `json:"document"`
	Identifier *NodeValue `json:"identifier"`
}

type InsertDocument struct {
	Query
	Document   Value      `json:"document"`
	Identifier *NodeValue `json:"identifier"`
}

type DeleteDocument struct {
	Query
	Identifier NodeValue `json:"identifier"`
}

type QueryResource struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Source  Source                `json:"source"`
	Format  FormatType            `json:"format"`
	Options *QueryResourceOptions ` json:"options" terminusgo:"class=xdd:json"`
}

// FIXME: hack to comply Querier, figure out why it is a part of query in python client, but in not in schema
func (q QueryResource) GetQuery() Querier {
	return q
}

// FIXME: figure out full list of options, what they mean and why they are used in python client
type QueryResourceOptions struct {
	Format       string `json:"format"`
	FormatHeader bool   `json:"format_header"`
}

type Get struct {
	Query
	Columns   []Column      `json:"columns"`
	Resource  QueryResource `json:"resource"`
	HasHeader *bool         `json:"has_header"`
}

// TODO: not used
type AddedData struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    DataValue `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type AddLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type AddedLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type AddData struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    DataValue `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type DeleteLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type DeletedLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

// TODO: not used
type If struct {
	Query
	Test Querier `json:"test" terminusgo:"type=Class,class=Query"`
	Then Querier `json:"then" terminusgo:"type=Class,class=Query"`
	Else Querier `json:"else" terminusgo:"type=Class,class=Query"`
}

type Trim struct {
	Query
	Untrimmed DataValue `json:"untrimmed"`
	Trimmed   DataValue `json:"trimmed"`
}

type IsA struct {
	Query
	Element NodeValue `json:"element"`
	Type    NodeValue `json:"type"`
}

type Like struct {
	Query
	Left       DataValue `json:"left"`
	Right      DataValue `json:"right"`
	Similarity DataValue `json:"similarity"` // FIXME: figure out why it's of type Value in python client
}

type Less struct {
	Query
	Left  DataValue `json:"left"` // FIXME: figure out why these fields have Value type in python client
	Right DataValue `json:"right"`
}

type Greater struct {
	Query
	Left  DataValue `json:"left"` // FIXME: figure out why these fields have Value type in python client
	Right DataValue `json:"right"`
}

type Optional struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

type LexicalKey struct {
	Query
	Base    DataValue   `json:"base"`
	KeyList []DataValue `json:"key_list"`
	URI     NodeValue   `json:"uri"`
}

type RandomKey struct {
	Query
	Base DataValue `json:"base"`
	URI  NodeValue `json:"uri"`
}

type HashKey struct {
	Query
	Base    DataValue   `json:"base"`
	KeyList []DataValue `json:"key_list"`
	URI     NodeValue   `json:"uri"`
}

type Upper struct {
	Query
	Mixed DataValue `json:"mixed"`
	Upper DataValue `json:"upper"`
}

type Lower struct {
	Query
	Mixed DataValue `json:"mixed"`
	Lower DataValue `json:"lower"`
}

type Pad struct {
	Query
	String DataValue `json:"string"`
	Char   DataValue `json:"char"`
	Times  DataValue `json:"times"`
	Result DataValue `json:"result"`
}

type Split struct {
	Query
	String  DataValue `json:"string"`
	Pattern DataValue `json:"pattern"`
	List    DataValue `json:"list"`
}

type Member struct {
	Query
	Member DataValue `json:"member"`
	List   DataValue `json:"list"`
}

type Concatenate struct {
	Query
	List   DataValue `json:"list"`
	Result DataValue `json:"result"`
}

type Join struct {
	Query
	List      DataValue `json:"list"`
	Separator DataValue `json:"separator"`
	Result    DataValue `json:"result"`
}

type Sum struct {
	Query
	List   DataValue `json:"list"`
	Result DataValue `json:"result"`
}

type Start struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Start    uint    `json:"start"`
}

type Limit struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Limit    uint    `json:"limit"`
}

type Regexp struct {
	Query
	Pattern DataValue  `json:"pattern"`
	String  DataValue  `json:"string"`
	Result  *DataValue `json:"result"`
}

type True struct {
	Query
}

type OrderBy struct {
	Query
	SubQuery Querier         `json:"query" terminusgo:"type=Class,class=Query"`
	Ordering []OrderTemplate `json:"ordering"`
}

type GroupBy struct {
	Query
	Template Value    `json:"template"`
	GroupBy  []string `json:"group_by"`
	SubQuery Querier  `json:"query" terminusgo:"type=Class,class=Query"`
	Grouped  Value    `json:"grouped"`
}

type Length struct {
	Query
	List   DataValue `json:"list"`
	Length DataValue `json:"length"`
}

type Not struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

type Once struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

type Immediately struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

type Count struct {
	Query
	SubQuery Querier   `json:"query" terminusgo:"type=Class,class=Query"`
	Count    DataValue `json:"count"`
}

type Typecast struct {
	Query
	Value  Value     `json:"value"`
	Type   NodeValue `json:"type"`
	Result Value     `json:"result"`
}

type Dot struct {
	Query
	Document DataValue `json:"document"`
	Field    DataValue `json:"field"`
	Value    DataValue `json:"value"`
}

type Size struct {
	Query
	Resource string    `json:"resource"`
	Size     DataValue `json:"size"`
}

type TripleCount struct {
	Query
	Resource string    `json:"resource"`
	Count    DataValue `json:"count"`
}

type TypeOf struct {
	Query
	Value Value     `json:"value"`
	Type  NodeValue `json:"type"`
}

type Eval struct {
	Query
	Expression ArithmeticExpression `json:"expression"`
	Result     ArithmeticValue      `json:"result"`
}
