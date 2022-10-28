package schema

type TerminusSchema struct {
	Context Context
}

type ItemType string

const (
	ItemTypeClass       ItemType = "Class"
	ItemTypeEnum        ItemType = "Enum"
	ItemTypeTaggedUnion ItemType = "TaggedUnion"
	ItemTypeUnit        ItemType = "Unit"
)
