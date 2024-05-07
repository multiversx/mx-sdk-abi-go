package abi

// InputMultiValue is a multi-value (used for encoding)
type InputMultiValue struct {
	Items []any
}

// OutputMultiValue is a multi-value (used for decoding)
type OutputMultiValue struct {
	Items []any
}

// InputVariadicValues holds variadic values (used for encoding)
type InputVariadicValues struct {
	Items []any
}

// OutputVariadicValues holds variadic values (used for decoding)
type OutputVariadicValues struct {
	Items       []any
	ItemCreator func() any
}

// InputOptionalValue holds an optional value (used for encoding)
type InputOptionalValue struct {
	Value any
}

// OutputOptionalValue holds an optional value (used for decoding)
type OutputOptionalValue struct {
	Value    any
	HasValue bool
}
