package abi

import (
	"bytes"
	"fmt"
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

func (c *codec) encodeTopLevelOption(writer io.Writer, value OptionValue) error {
	if value.Value == nil {
		return nil
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

func (c *codec) decodeTopLevelOption(data []byte, value *OptionValue) error {
	if len(data) == 0 {
		value.Value = nil
		return nil
	}

	firstByte := data[0]
	dataAfterFirstByte := data[1:]

	if firstByte != 0x01 {
		return fmt.Errorf("invalid first byte for top-level encoded option: %d", firstByte)
	}

	reader := bytes.NewReader(dataAfterFirstByte)
	return c.doDecodeNested(reader, value.Value)
}
