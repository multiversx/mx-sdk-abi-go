package abi

// U8Value is a wrapper for uint8
type U8Value struct {
	Value uint8
}

// U16Value is a wrapper for uint16
type U16Value struct {
	Value uint16
}

// U32Value is a wrapper for uint32
type U32Value struct {
	Value uint32
}

// U64Value is a wrapper for uint64
type U64Value struct {
	Value uint64
}

// I8Value is a wrapper for int8
type I8Value struct {
	Value int8
}

// I16Value is a wrapper for int16
type I16Value struct {
	Value int16
}

// I32Value is a wrapper for int32
type I32Value struct {
	Value int32
}

// I64Value is a wrapper for int64
type I64Value struct {
	Value int64
}

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
