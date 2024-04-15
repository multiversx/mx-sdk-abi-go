package abi

import (
	"io"
)

func (c *codec) encodeNestedString(writer io.Writer, value StringValue) error {
	data := []byte(value.Value)
	err := encodeLength(writer, uint32(len(data)))
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (c *codec) encodeTopLevelString(writer io.Writer, value StringValue) error {
	_, err := writer.Write([]byte(value.Value))
	return err
}

func (c *codec) decodeNestedString(reader io.Reader, value *StringValue) error {
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = string(data)
	return nil
}

func (c *codec) decodeTopLevelString(data []byte, value *StringValue) error {
	value.Value = string(data)
	return nil
}
