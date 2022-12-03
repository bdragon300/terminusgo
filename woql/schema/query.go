package schema

import "github.com/bdragon300/terminusgo/schema"

// TODO: check if all schema structs are actually used

// TODO: add json marshaler/unmarshaler which is not require json tags everywhere (and applies `omitempty` everywhere)
// TODO: check schema against json_woql.pl

type NamedQuery struct {
	Name  string  `json:"name"`
	Query Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

func (n NamedQuery) GetSubQuery() Querier {
	return n.Query
}

type NamedParametricQuery struct {
	Name       string   `json:"name"`
	Parameters []string `json:"parameters"`
	Query      Querier  `json:"query" terminusgo:"type=Class,class=Query"`
}

func (n NamedParametricQuery) GetSubQuery() Querier {
	return n.Query
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

func (q Query) QueryType() {}

type SubQuerier interface {
	GetSubQuery() Querier
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

func (n Using) GetSubQuery() Querier {
	return n.SubQuery
}

type Select struct {
	Query
	SubQuery  Querier  `json:"query" terminusgo:"type=Class,class=Query"`
	Variables []string `json:"variables"`
}

func (n Select) GetSubQuery() Querier {
	return n.SubQuery
}

type Distinct struct {
	Query
	SubQuery  Querier  `json:"query" terminusgo:"type=Class,class=Query"`
	Variables []string `json:"variables"`
}

func (n Distinct) GetSubQuery() Querier {
	return n.SubQuery
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

func (n From) GetSubQuery() Querier {
	return n.SubQuery
}

type Into struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Graph    string  `json:"graph"`
}

func (n Into) GetSubQuery() Querier {
	return n.SubQuery
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
}

type AddedTriple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
}

type DeleteTriple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     *string   `json:"graph"`
}

type RemovedTriple struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
}

type AddQuad struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     string    `json:"graph"`
}

type AddedQuad struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     string    `json:"graph"`
}

type RemovedQuad struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    Value     `json:"object"`
	Graph     string    `json:"graph"`
}

// TODO: not used
type Link struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

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

type ReadObject struct {
	Query
	Identifier NodeValue `json:"identifier"`
	Document   Value     `json:"document"`
}

type UpdateObject struct {
	Query
	Document Value `json:"document"`
}

type DeleteObject struct {
	Query
	Identifier NodeValue `json:"identifier"`
}

type QueryResource struct {
	// TODO: type is TaggedUnion
	*schema.SubDocumentModel
	Source  Source      `json:"source"`
	Format  FormatType  `json:"format"`
	Options FileOptions `json:"options" terminusgo:"class=xdd:json,optional"`
}

type Get struct {
	Query
	Columns   []Column      `json:"columns"`
	Resource  QueryResource `json:"resource"`
	HasHeader *bool         `json:"has_header"`
}

type AddedData struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    DataValue `json:"object"`
	Graph     *string   `json:"graph"`
}

type AddLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

type AddedLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

type AddData struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    DataValue `json:"object"`
	Graph     *string   `json:"graph"`
}

type DeleteLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

type DeletedLink struct {
	Query
	Subject   NodeValue `json:"subject"`
	Predicate NodeValue `json:"predicate"`
	Object    NodeValue `json:"object"`
	Graph     *string   `json:"graph"`
}

type If struct {
	Query
	Test Querier `json:"test" terminusgo:"type=Class,class=Query"`
	Then Querier `json:"then" terminusgo:"type=Class,class=Query"`
	Else Querier `json:"else" terminusgo:"type=Class,class=Query"`
}

// TODO: implement somehow
// func (n If) GetSubQuery() Querier {
//	return n.SubQuery
//}

type When struct {
	Query
	SubQuery   Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Consequent Querier `json:"consequent" terminusgo:"type=Class,class=Query"`
}

// TODO: implement SubQuerier

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

func (n Optional) GetSubQuery() Querier {
	return n.SubQuery
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

func (n Start) GetSubQuery() Querier {
	return n.SubQuery
}

type Limit struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
	Limit    uint    `json:"limit"`
}

func (n Limit) GetSubQuery() Querier {
	return n.SubQuery
}

type Regexp struct {
	Query
	Pattern DataValue `json:"pattern"`
	String  DataValue `json:"string"`
	Result  DataValue `json:"result"`
}

type True struct {
	Query
}

type OrderBy struct {
	Query
	SubQuery Querier         `json:"query" terminusgo:"type=Class,class=Query"`
	Ordering []OrderTemplate `json:"ordering"`
}

func (n OrderBy) GetSubQuery() Querier {
	return n.SubQuery
}

type GroupBy struct {
	Query
	Template Value    `json:"template"`
	GroupBy  []string `json:"group_by"`
	SubQuery Querier  `json:"query" terminusgo:"type=Class,class=Query"`
	Grouped  Value    `json:"grouped"`
}

func (n GroupBy) GetSubQuery() Querier {
	return n.SubQuery
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

func (n Not) GetSubQuery() Querier {
	return n.SubQuery
}

type Once struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

func (n Once) GetSubQuery() Querier {
	return n.SubQuery
}

type Immediately struct {
	Query
	SubQuery Querier `json:"query" terminusgo:"type=Class,class=Query"`
}

func (n Immediately) GetSubQuery() Querier {
	return n.SubQuery
}

type Count struct {
	Query
	SubQuery Querier   `json:"query" terminusgo:"type=Class,class=Query"`
	Count    DataValue `json:"count"`
}

func (n Count) GetSubQuery() Querier {
	return n.SubQuery
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
	Expression ArithmeticExpressionType `json:"expression"`
	Result     ArithmeticValue          `json:"result"`
}
