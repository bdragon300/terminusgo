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
	Predicate *string
}

type InversePathPredicate struct {
	PathPattern
	Predicate *string
}

type PathSequence struct {
	PathPattern
	Sequence []PathPattern
}

type PathOr struct {
	PathPattern
	Or []PathPattern
}

type PathPlus struct {
	PathPattern
	Plus PathPattern
}

type PathStar struct {
	PathPattern
	Star PathPattern
}

type PathTimes struct {
	PathPattern
	Times PathPattern
	From  uint
	To    uint
}
