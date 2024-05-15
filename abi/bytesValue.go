package abi

import (
	"io"
)

// BytesValue is a wrapper for a byte slice
type BytesValue struct {
	Value []byte
}

func (value *BytesValue) encodeNested(writer io.Writer) error {
	err := encodeLength(writer, uint32(len(value.Value)))
	if err != nil {
		return err
	}

	_, err = writer.Write(value.Value)
	return err
}

func (value *BytesValue) encodeTopLevel(writer io.Writer) error {
	_, err := writer.Write(value.Value)
	return err
}

func (value *BytesValue) decodeNested(reader io.Reader) error {
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

func (value *BytesValue) decodeTopLevel(data []byte) error {
	value.Value = data
	return nil
}
