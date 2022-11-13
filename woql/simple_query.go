package woql

import (
	"github.com/bdragon300/terminusgo/woql/bare"
	"github.com/bdragon300/terminusgo/woql/schema"
)

// type QueryMaker interface {
//	RootQuery() schema.QueryType
//}

func NewSimpleQuery() *SimpleQueryBuilder {
	return &SimpleQueryBuilder{}
}

type SimpleQueryBuilder struct {
	bare.QueryBuilder
}

// func (s *SimpleQueryBuilder) And() *SimpleQueryBuilder {
//
// }
// And
// Or
// Not
// Select
// From
// Using
// Distinct
// Into
// OrderBy
// GroupBy
// TripleCount
// Triple
// AddTriple
// DeleteTriple
// AddedTriple
// Quad
// AddQuad
// DeleteQuad
// AddedQuad
// Subsumption
// Equals
// SubString
// UpdateDocument
// InsertDocument
// DeleteDocument
// ReadDocument
// Get
// File
// Once
// Remote
// Post
// Trim
// Upper
// Lower
// Concatenate
// Join
// Regexp
// Eval
// IsA
// Like
// Less
// Greater
// Optional
// HashKey
// LexicalKey
// RandomKey
// Pad
// Split
// Dot
// Member
// Sum
// Start
// Limit
// Length
// Immediately
// Count
// TypeCast
// TypeOf
// True
// Path
// Size

func (q *SimpleQueryBuilder) Triple(subj, pred, obj string) *SimpleQueryBuilder {
	q.QueryBuilder.Triple(
		*ParseVariable(subj, &schema.NodeValue{}),
		*ParseVariable(pred, &schema.NodeValue{}),
		*ParseVariable(obj, &schema.Value{}),
	)
	return q
}

// func (q *SimpleQueryBuilder) Path(subj, pattern, obj string) *SimpleQueryBuilder {
//	p, err := ParseTriplePattern(pattern)
//	if err != nil {
//		q.QueryBuilder.errors = multierror.Append(q.QueryBuilder.errors, err)
//	}
//	q.QueryBuilder.Path(*ParseVariable(subj, &schema.Value{}), p, *ParseVariable(obj, &schema.Value{}))
//	return q
//}

// func (q *SimpleQueryBuilder) And(queries... QueryMaker) *SimpleQueryBuilder {
//	appendQueries(&q.QueryBuilder, queries...)
//	return q
//}

// func appendQueries[T schema.QueryType](q *query.QueryBuilder, items... T) {
//	for _, query := range items {
//		q.agg = append(q.agg, query)
//	}
//}
