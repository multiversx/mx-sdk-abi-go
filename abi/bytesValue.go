package abi

import (
	"io"
)

type codecForBytes struct {
}

func (c *codecForBytes) encodeNested(writer io.Writer, value BytesValue) error {
	err := encodeLength(writer, uint32(len(value.Value)))
	if err != nil {
		return err
	}

	_, err = writer.Write(value.Value)
	return err
}

func (c *codecForBytes) encodeTopLevel(writer io.Writer, value BytesValue) error {
	_, err := writer.Write(value.Value)
	return err
}

func (c *codecForBytes) decodeNested(reader io.Reader, value *BytesValue) error {
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

func (c *codecForBytes) decodeTopLevel(data []byte, value *BytesValue) error {
	value.Value = data
	return nil
}
