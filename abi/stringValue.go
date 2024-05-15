package abi

import (
	"io"
)

// StringValue is a wrapper for a string
type StringValue struct {
	Value string
}

func (value *StringValue) encodeNested(writer io.Writer) error {
	data := []byte(value.Value)
	err := encodeLength(writer, uint32(len(data)))
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (value *StringValue) encodeTopLevel(writer io.Writer) error {
	_, err := writer.Write([]byte(value.Value))
	return err
}

func (value *StringValue) decodeNested(reader io.Reader) error {
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

func (value *StringValue) decodeTopLevel(data []byte) error {
	value.Value = string(data)
	return nil
}
