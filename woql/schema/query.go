package schema

import "github.com/bdragon300/terminusgo/schema"

// TODO: check if all schema structs are actually used

// TODO: add json marshaler/unmarshaler which is not require json tags everywhere (and applies `omitempty` everywhere)
// TODO: check schema against json_woql.pl

type NamedQuery struct {
	Name     string
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
}

func (n NamedQuery) GetSubQuery() Querier {
	return n.SubQuery
}

type NamedParametricQuery struct {
	Name       string
	Parameters []string
	SubQuery   Querier `json:"query" terminusgo:"type=Class,class=Query,name=query"`
}

func (n NamedParametricQuery) GetSubQuery() Querier {
	return n.SubQuery
}

type Querier interface {
	GetQueryData() any
}

type Query struct {
	*schema.AbstractModel
	*schema.SubDocumentModel
}

func (q Query) GetQueryData() any {
	return &q
}

func (q Query) QueryType() {}

type SubQuerier interface { // FIXME: seems that not needed
	GetSubQuery() Querier
}

type Path struct {
	Query
	Subject Value
	Pattern PathPatternType `terminusgo:"type=Class,class=PathPattern"`
	Object  Value
	Path    *Value
}

type Using struct {
	Query
	SubQuery   Querier `terminusgo:"type=Class,class=Query,name=query"`
	Collection string
}

func (n Using) GetSubQuery() Querier {
	return n.SubQuery
}

type Select struct {
	Query
	SubQuery  Querier `json:"query" terminusgo:"type=Class,class=Query,name=query"`
	Variables []string
}

func (n Select) GetSubQuery() Querier {
	return n.SubQuery
}

type Distinct struct {
	Query
	SubQuery  Querier `terminusgo:"type=Class,class=Query,name=query"`
	Variables []string
}

func (n Distinct) GetSubQuery() Querier {
	return n.SubQuery
}

type And struct {
	Query
	And []Querier `terminusgo:"type=List,class=Query"`
}

type Or struct {
	Query
	Or []Querier `terminusgo:"type=List,class=Query"`
}

type From struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
	Graph    string
}

func (n From) GetSubQuery() Querier {
	return n.SubQuery
}

type Into struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
	Graph    string
}

func (n Into) GetSubQuery() Querier {
	return n.SubQuery
}

type Triple struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
	Graph     *string
}

type AddTriple struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
}

type AddedTriple struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
}

type DeleteTriple struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
	Graph     *string
}

type RemovedTriple struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
}

type AddQuad struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
	Graph     string
}

type AddedQuad struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
	Graph     string
}

type RemovedQuad struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    Value
	Graph     string
}

// TODO: not used
type Link struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    NodeValue
	Graph     *string
}

type Data struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    DataValue
	Graph     *string
}

type Subsumption struct {
	Query
	Child  NodeValue
	Parent NodeValue
}

type Equals struct {
	Query
	Left  DataValue
	Right DataValue
}

type SubString struct {
	Query
	String    DataValue
	Length    DataValue
	Before    DataValue
	After     DataValue
	SubString DataValue `terminusgo:"name=substring"`
}

type ReadDocument struct {
	Query
	Document   Value
	Identifier NodeValue
}

type UpdateDocument struct {
	Query
	Document   Value
	Identifier *NodeValue
}

type InsertDocument struct {
	Query
	Document   Value
	Identifier *NodeValue
}

type DeleteDocument struct {
	Query
	Identifier NodeValue
}

type ReadObject struct {
	Query
	Identifier NodeValue
	Document   Value
}

type UpdateObject struct {
	Query
	Document Value
}

type DeleteObject struct {
	Query
	Identifier NodeValue
}

type QueryResource struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Source  Source
	Format  FormatType
	Options FileOptions `terminusgo:"class=xdd:json,optional,untyped_object"`
}

type Get struct {
	Query
	Columns   []Column
	Resource  QueryResource
	HasHeader *bool
}

type Put struct {
	Query
	Columns  []Column
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
	Resource QueryResource
}

type AddedData struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    DataValue
	Graph     *string
}

type AddLink struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    NodeValue
	Graph     *string
}

type AddedLink struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    NodeValue
	Graph     *string
}

type AddData struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    DataValue
	Graph     *string
}

type DeleteLink struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    NodeValue
	Graph     *string
}

type DeletedLink struct {
	Query
	Subject   NodeValue
	Predicate NodeValue
	Object    NodeValue
	Graph     *string
}

type If struct {
	Query
	Test Querier `terminusgo:"type=Class,class=Query"`
	Then Querier `terminusgo:"type=Class,class=Query"`
	Else Querier `terminusgo:"type=Class,class=Query"`
}

// TODO: implement somehow
// func (n If) GetSubQuery() Querier {
//	return n.SubQuery
//}

type When struct {
	Query
	SubQuery   Querier `terminusgo:"type=Class,class=Query,name=query"`
	Consequent Querier `terminusgo:"type=Class,class=Query"`
}

// TODO: implement SubQuerier

type Trim struct {
	Query
	Untrimmed DataValue
	Trimmed   DataValue
}

type IsA struct {
	Query
	Element NodeValue
	Type    NodeValue
}

type Like struct {
	Query
	Left       DataValue
	Right      DataValue
	Similarity DataValue // FIXME: figure out why it's of type Value in python client
}

type Less struct {
	Query
	Left  DataValue // FIXME: figure out why these fields have Value type in python client
	Right DataValue
}

type Greater struct {
	Query
	Left  DataValue // FIXME: figure out why these fields have Value type in python client
	Right DataValue
}

type Optional struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
}

func (n Optional) GetSubQuery() Querier {
	return n.SubQuery
}

type LexicalKey struct {
	Query
	Base    DataValue
	KeyList []DataValue
	URI     NodeValue
}

type RandomKey struct {
	Query
	Base DataValue
	URI  NodeValue
}

type HashKey struct {
	Query
	Base    DataValue
	KeyList []DataValue
	URI     NodeValue
}

type Upper struct {
	Query
	Mixed DataValue
	Upper DataValue
}

type Lower struct {
	Query
	Mixed DataValue
	Lower DataValue
}

type Pad struct {
	Query
	String DataValue
	Char   DataValue
	Times  DataValue
	Result DataValue
}

type Split struct {
	Query
	String  DataValue
	Pattern DataValue
	List    DataValue
}

type Member struct {
	Query
	Member DataValue
	List   DataValue
}

type Concatenate struct {
	Query
	List   DataValue
	Result DataValue
}

type Join struct {
	Query
	List      DataValue
	Separator DataValue
	Result    DataValue
}

type Sum struct {
	Query
	List   DataValue
	Result DataValue
}

type Start struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
	Start    uint
}

func (n Start) GetSubQuery() Querier {
	return n.SubQuery
}

type Limit struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
	Limit    uint
}

func (n Limit) GetSubQuery() Querier {
	return n.SubQuery
}

type Regexp struct {
	Query
	Pattern DataValue
	String  DataValue
	Result  DataValue
}

type True struct {
	Query
}

type OrderBy struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
	Ordering []OrderTemplate
}

func (n OrderBy) GetSubQuery() Querier {
	return n.SubQuery
}

type GroupBy struct {
	Query
	Template Value
	GroupBy  []string
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query,name=query"`
	Grouped  Value
}

func (n GroupBy) GetSubQuery() Querier {
	return n.SubQuery
}

type Length struct {
	Query
	List   DataValue
	Length DataValue
}

type Not struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
}

func (n Not) GetSubQuery() Querier {
	return n.SubQuery
}

type Once struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
}

func (n Once) GetSubQuery() Querier {
	return n.SubQuery
}

type Immediately struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
}

func (n Immediately) GetSubQuery() Querier {
	return n.SubQuery
}

type Count struct {
	Query
	SubQuery Querier `terminusgo:"type=Class,class=Query,name=query"`
	Count    DataValue
}

func (n Count) GetSubQuery() Querier {
	return n.SubQuery
}

type Typecast struct {
	Query
	Value  Value
	Type   NodeValue
	Result Value
}

type Dot struct {
	Query
	Document DataValue
	Field    DataValue
	Value    DataValue
}

type Size struct {
	Query
	Resource string
	Size     DataValue
}

type TripleCount struct {
	Query
	Resource string
	Count    DataValue
}

type TypeOf struct {
	Query
	Value Value
	Type  NodeValue
}

type Eval struct {
	Query
	Expression ArithmeticExpressionType
	Result     ArithmeticValue
}
