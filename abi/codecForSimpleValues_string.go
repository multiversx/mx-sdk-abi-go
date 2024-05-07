package abi

import (
	"io"
)

type codecForString struct {
}

func (c *codecForString) encodeNested(writer io.Writer, value StringValue) error {
	data := []byte(value.Value)
	err := encodeLength(writer, uint32(len(data)))
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (c *codecForString) encodeTopLevel(writer io.Writer, value StringValue) error {
	_, err := writer.Write([]byte(value.Value))
	return err
}

func (c *codecForString) decodeNested(reader io.Reader, value *StringValue) error {
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

func (c *codecForString) decodeTopLevel(data []byte, value *StringValue) error {
	value.Value = string(data)
	return nil
}
