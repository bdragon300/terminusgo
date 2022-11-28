package schema

import "github.com/bdragon300/terminusgo/schema"

// TODO: not used any of this

type PathPatternType interface {
	PathPattern()
}

type PathPattern struct {
	*schema.AbstractModel
}

func (p PathPattern) PathPattern() {}

type PathPredicate struct {
	PathPattern
	Predicate *string `json:"predicate"`
}

type InversePathPredicate struct {
	PathPattern
	Predicate *string `json:"predicate"`
}

type PathSequence struct {
	PathPattern
	Sequence []PathPattern `json:"sequence"`
}

type PathOr struct {
	PathPattern
	Or []PathPattern `json:"or"`
}

type PathPlus struct {
	PathPattern
	Plus PathPattern `json:"plus"`
}

type PathStar struct {
	PathPattern
	Star PathPattern `json:"star"`
}

type PathTimes struct {
	PathPattern
	Times PathPattern `json:"times"`
	From  uint        `json:"from"`
	To    uint        `json:"to"`
}
