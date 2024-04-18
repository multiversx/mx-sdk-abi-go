package abi

import (
	"io"
)

func (c *codec) encodeNestedOption(writer io.Writer, value OptionValue) error {
	if value.Value == nil {
		_, err := writer.Write([]byte{0})
		return err
	}

	_, err := writer.Write([]byte{1})
	if err != nil {
		return err
	}

	return c.doEncodeNested(writer, value.Value)
}

func (c *codec) decodeNestedOption(reader io.Reader, value *OptionValue) error {
	bytes, err := readBytesExactly(reader, 1)
	if err != nil {
		return err
	}

	if bytes[0] == 0 {
		value.Value = nil
		return nil
	}

	return c.doDecodeNested(reader, value.Value)
}
