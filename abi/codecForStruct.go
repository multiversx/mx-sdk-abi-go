package abi

import (
	"bytes"
	"fmt"
	"io"
)

// StructValue is a struct (collection of fields)
type StructValue struct {
	Fields []Field
}

func (value *StructValue) encodeNested(writer io.Writer) error {
	for _, field := range value.Fields {
		err := field.Value.encodeNested(writer)
		if err != nil {
			return fmt.Errorf("cannot encode field '%s' of struct, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (value *StructValue) encodeTopLevel(writer io.Writer) error {
	return value.encodeNested(writer)
}

func (value *StructValue) decodeNested(reader io.Reader) error {
	for _, field := range value.Fields {
		err := field.Value.decodeNested(reader)
		if err != nil {
			return fmt.Errorf("cannot decode field '%s' of struct, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (value *StructValue) decodeTopLevel(data []byte) error {
	reader := bytes.NewReader(data)
	return value.decodeNested(reader)
}
