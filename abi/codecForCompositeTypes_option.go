package abi

import (
	"bytes"
	"fmt"
	"io"
)

type codecForOption struct {
	generalCodec *codec
}

func (c *codecForOption) encodeNested(writer io.Writer, value OptionValue) error {
	if value.Value == nil {
		_, err := writer.Write([]byte{optionMarkerForAbsentValue})
		return err
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return c.generalCodec.doEncodeNested(writer, value.Value)
}

func (c *codecForOption) encodeTopLevel(writer io.Writer, value OptionValue) error {
	if value.Value == nil {
		return nil
	}

	_, err := writer.Write([]byte{optionMarkerForPresentValue})
	if err != nil {
		return err
	}

	return c.generalCodec.doEncodeNested(writer, value.Value)
}

func (c *codecForOption) decodeNested(reader io.Reader, value *OptionValue) error {
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
		return c.generalCodec.doDecodeNested(reader, value.Value)
	}

	return fmt.Errorf("invalid first byte for nested encoded option: %d", firstByte)
}

func (c *codecForOption) decodeTopLevel(data []byte, value *OptionValue) error {
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
	return c.generalCodec.doDecodeNested(reader, value.Value)
}
