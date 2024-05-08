package abi

import (
	"bytes"
	"fmt"
	"io"
)

type codecForEnum struct {
	generalCodec generalCodec
}

func (c *codecForEnum) encodeNested(writer io.Writer, value EnumValue) error {
	err := c.generalCodec.doEncodeNested(writer, U8Value{Value: value.Discriminant})
	if err != nil {
		return err
	}

	for _, field := range value.Fields {
		err := c.generalCodec.doEncodeNested(writer, field.Value)
		if err != nil {
			return fmt.Errorf("cannot encode field '%s' of enum, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (c *codecForEnum) encodeTopLevel(writer io.Writer, value EnumValue) error {
	if value.Discriminant == 0 && len(value.Fields) == 0 {
		// Write nothing
		return nil
	}

	return c.encodeNested(writer, value)
}

func (c *codecForEnum) decodeNested(reader io.Reader, value *EnumValue) error {
	discriminant := &U8Value{}
	err := c.generalCodec.doDecodeNested(reader, discriminant)
	if err != nil {
		return err
	}

	value.Discriminant = discriminant.Value

	for _, field := range value.Fields {
		err := c.generalCodec.doDecodeNested(reader, field.Value)
		if err != nil {
			return fmt.Errorf("cannot decode field '%s' of enum, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (c *codecForEnum) decodeTopLevel(data []byte, value *EnumValue) error {
	if len(data) == 0 {
		value.Discriminant = 0
		return nil
	}

	reader := bytes.NewReader(data)
	return c.decodeNested(reader, value)
}
