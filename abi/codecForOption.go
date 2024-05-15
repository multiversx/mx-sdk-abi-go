package abi

import (
	"bytes"
	"fmt"
	"io"
)

// OptionValue is a wrapper for an option value
type OptionValue struct {
	Value singleValue
}

func (value *OptionValue) encodeNested(writer io.Writer) error {
	if value.Value == nil {
		_, err := writer.Write([]byte{optionMarkerForAbsentValue})
		return err
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return value.Value.encodeNested(writer)
}

func (value *OptionValue) encodeTopLevel(writer io.Writer) error {
	if value.Value == nil {
		return nil
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return value.Value.encodeNested(writer)
}

func (value *OptionValue) decodeNested(reader io.Reader) error {
	data, err := readBytesExactly(reader, 1)
	if err != nil {
		return err
	}

	firstByte := data[0]

	if firstByte == optionMarkerForAbsentValue {
		value.Value = nil
		return nil
	}

	if firstByte == optionMarkerForPresentValue {
		return value.Value.decodeNested(reader)
	}

	return fmt.Errorf("invalid first byte for nested encoded option: %d", firstByte)
}

func (value *OptionValue) decodeTopLevel(data []byte) error {
	if len(data) == 0 {
		value.Value = nil
		return nil
	}

	firstByte := data[0]
	dataAfterFirstByte := data[1:]

	if firstByte != optionMarkerForPresentValue {
		return fmt.Errorf("invalid first byte for top-level encoded option: %d", firstByte)
	}

	reader := bytes.NewReader(dataAfterFirstByte)
	return value.Value.decodeNested(reader)
}
