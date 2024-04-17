package abi

import (
	"io"
	"math/big"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

func (c *codec) encodeNestedBigNumber(writer io.Writer, value *big.Int) error {
	data := twos.ToBytes(value)
	dataLength := len(data)

	// Write the length of the payload
	err := encodeLength(writer, uint32(dataLength))
	if err != nil {
		return err
	}

	// Write the payload
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *codec) encodeTopLevelBigNumber(writer io.Writer, value *big.Int) error {
	data := twos.ToBytes(value)
	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *codec) decodeNestedBigNumber(reader io.Reader) (*big.Int, error) {
	// Read the length of the payload
	length, err := decodeLength(reader)
	if err != nil {
		return nil, err
	}

	// Read the payload
	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return nil, err
	}

	return twos.FromBytes(data), nil
}

func (c *codec) decodeTopLevelBigNumber(data []byte) *big.Int {
	return twos.FromBytes(data)
}
