package abi

import (
	"io"
	"math/big"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

func (c *codec) encodeNestedBigNumber(writer io.Writer, value *big.Int, withSign bool) error {
	data := bigIntToBytes(value, withSign)
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

func (c *codec) encodeTopLevelBigNumber(writer io.Writer, value *big.Int, withSign bool) error {
	data := bigIntToBytes(value, withSign)
	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *codec) decodeNestedBigNumber(reader io.Reader, withSign bool) (*big.Int, error) {
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

	return bigIntFromBytes(data, withSign), nil
}

func (c *codec) decodeTopLevelBigNumber(data []byte, withSign bool) *big.Int {
	return bigIntFromBytes(data, withSign)
}

func bigIntToBytes(value *big.Int, withSign bool) []byte {
	if withSign {
		return twos.ToBytes(value)
	}

	return value.Bytes()
}

func bigIntFromBytes(data []byte, withSign bool) *big.Int {
	if withSign {
		return twos.FromBytes(data)
	}

	return big.NewInt(0).SetBytes(data)
}
