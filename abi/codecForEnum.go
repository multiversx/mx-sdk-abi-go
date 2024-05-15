package abi

import (
	"bytes"
	"fmt"
	"io"
)

// EnumValue is an enum (discriminant and fields)
type EnumValue struct {
	Discriminant uint8
	Fields       []Field
}

func (value *EnumValue) encodeNested(writer io.Writer) error {
	discriminant := U8Value{Value: value.Discriminant}
	err := discriminant.encodeNested(writer)
	if err != nil {
		return err
	}

	for _, field := range value.Fields {
		err := field.Value.encodeNested(writer)
		if err != nil {
			return fmt.Errorf("cannot encode field '%s' of enum, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (value *EnumValue) encodeTopLevel(writer io.Writer) error {
	if value.Discriminant == 0 && len(value.Fields) == 0 {
		// Write nothing
		return nil
	}

	return value.encodeNested(writer)
}

func (value *EnumValue) decodeNested(reader io.Reader) error {
	discriminant := &U8Value{}
	err := discriminant.decodeNested(reader)
	if err != nil {
		return err
	}

	value.Discriminant = discriminant.Value

	for _, field := range value.Fields {
		err := field.Value.decodeNested(reader)
		if err != nil {
			return fmt.Errorf("cannot decode field '%s' of enum, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (value *EnumValue) decodeTopLevel(data []byte) error {
	if len(data) == 0 {
		value.Discriminant = 0
		return nil
	}

	reader := bytes.NewReader(data)
	return value.decodeNested(reader)
}
