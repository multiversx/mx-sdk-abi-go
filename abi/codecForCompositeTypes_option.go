package abi

import (
	"bytes"
	"fmt"
	"io"
)

func (c *codec) encodeNestedOption(writer io.Writer, value OptionValue) error {
	if value.Value == nil {
		_, err := writer.Write([]byte{optionMarkerForAbsentValue})
		return err
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return c.doEncodeNested(writer, value.Value)
}

func (c *codec) encodeTopLevelOption(writer io.Writer, value OptionValue) error {
	if value.Value == nil {
		return nil
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return c.doEncodeNested(writer, value.Value)
}

func (c *codec) decodeNestedOption(reader io.Reader, value *OptionValue) error {
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
		return c.doDecodeNested(reader, value.Value)
	}

	return fmt.Errorf("invalid first byte for nested encoded option: %d", firstByte)
}

func (c *codec) decodeTopLevelOption(data []byte, value *OptionValue) error {
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
	return c.doDecodeNested(reader, value.Value)
}
