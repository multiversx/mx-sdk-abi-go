package abi

import (
	"bytes"
	"fmt"
	"io"
)

type codeForStruct struct {
	generalCodec *codec
}

func (c *codeForStruct) encodeNested(writer io.Writer, value StructValue) error {
	for _, field := range value.Fields {
		err := c.generalCodec.doEncodeNested(writer, field.Value)
		if err != nil {
			return fmt.Errorf("cannot encode field '%s' of struct, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (c *codeForStruct) encodeTopLevel(writer io.Writer, value StructValue) error {
	return c.encodeNested(writer, value)
}

func (c *codeForStruct) decodeNested(reader io.Reader, value *StructValue) error {
	for _, field := range value.Fields {
		err := c.generalCodec.doDecodeNested(reader, field.Value)
		if err != nil {
			return fmt.Errorf("cannot decode field '%s' of struct, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (c *codeForStruct) decodeTopLevel(data []byte, value *StructValue) error {
	reader := bytes.NewReader(data)
	return c.decodeNested(reader, value)
}
