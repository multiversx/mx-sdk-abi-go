package abi

// OptionValue is a wrapper for an option value
type OptionValue struct {
	Value any
}

// Field is a field in a struct, enum etc.
type Field struct {
	Name  string
	Value any
}

// StructValue is a struct (collection of fields)
type StructValue struct {
	Fields []Field
}

// EnumValue is an enum (discriminant and fields)
type EnumValue struct {
	Discriminant uint8
	Fields       []Field
}

// InputListValue is a list of values (used for encoding)
type InputListValue struct {
	Items []any
}

// OutputListValue is a list of values (used for decoding)
type OutputListValue struct {
	Items       []any
	ItemCreator func() any
}
