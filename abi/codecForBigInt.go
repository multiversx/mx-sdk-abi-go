package abi

import (
	"io"
	"math/big"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

// BigUIntValue is a wrapper for a big integer (unsigned)
type BigUIntValue struct {
	Value *big.Int
}

func (value *BigUIntValue) encodeNested(writer io.Writer) error {
	data := value.Value.Bytes()
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

func (value *BigUIntValue) encodeTopLevel(writer io.Writer) error {
	data := value.Value.Bytes()
	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (value *BigUIntValue) decodeNested(reader io.Reader) error {
	// Read the length of the payload
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	// Read the payload
	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = big.NewInt(0).SetBytes(data)
	return nil
}

func (value *BigUIntValue) decodeTopLevel(data []byte) {
	value.Value = big.NewInt(0).SetBytes(data)
}

// BigIntValue is a wrapper for a big integer (signed)
type BigIntValue struct {
	Value *big.Int
}

func (value *BigIntValue) encodeNested(writer io.Writer) error {
	data := twos.ToBytes(value.Value)
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

func (value *BigIntValue) encodeTopLevel(writer io.Writer) error {
	data := twos.ToBytes(value.Value)
	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (value *BigIntValue) decodeNested(reader io.Reader) error {
	// Read the length of the payload
	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	// Read the payload
	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = twos.FromBytes(data)
	return nil
}

func (value *BigIntValue) decodeTopLevel(data []byte) {
	value.Value = twos.FromBytes(data)
}
