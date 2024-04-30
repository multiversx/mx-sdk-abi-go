package abi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

type codecForSmallInt struct {
}

func (c *codecForSmallInt) encodeNested(writer io.Writer, value any, numBytes int) error {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, value)
	if err != nil {
		return err
	}

	data := buffer.Bytes()
	if len(data) != numBytes {
		return fmt.Errorf("unexpected number of bytes: %d != %d", len(data), numBytes)
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *codecForSmallInt) encodeTopLevelUnsigned(writer io.Writer, value uint64) error {
	b := big.NewInt(0).SetUint64(value)
	data := b.Bytes()
	_, err := writer.Write(data)
	return err
}

func (c *codecForSmallInt) encodeTopLevelSigned(writer io.Writer, value int64) error {
	data := twos.ToBytes(big.NewInt(value))
	_, err := writer.Write(data)
	return err
}

func (c *codecForSmallInt) decodeNested(reader io.Reader, value any, numBytes int) error {
	data, err := readBytesExactly(reader, numBytes)
	if err != nil {
		return err
	}

	buffer := bytes.NewReader(data)
	err = binary.Read(buffer, binary.BigEndian, value)
	if err != nil {
		return err
	}

	return nil
}

func (c *codecForSmallInt) decodeTopLevelUnsigned(data []byte, maxValue uint64) (uint64, error) {
	b := big.NewInt(0).SetBytes(data)
	if !b.IsUint64() {
		return 0, fmt.Errorf("decoded value is too large or invalid: %s", b)
	}

	n := b.Uint64()
	if n > maxValue {
		return 0, fmt.Errorf("decoded value is too large: %d > %d", n, maxValue)
	}

	return n, nil
}

func (c *codecForSmallInt) decodeTopLevelSigned(data []byte, maxValue int64) (int64, error) {
	b := twos.FromBytes(data)

	if !b.IsInt64() {
		return 0, fmt.Errorf("decoded value is too large or invalid: %s", b)
	}

	n := b.Int64()
	if n > maxValue {
		return 0, fmt.Errorf("decoded value is too large: %d > %d", n, maxValue)
	}

	return n, nil
}
