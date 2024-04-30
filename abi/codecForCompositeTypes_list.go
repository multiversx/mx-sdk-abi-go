package abi

import (
	"bytes"
	"errors"
	"io"
)

type codecForList struct {
	generalCodec generalCodec
}

func (c *codecForList) encodeNested(writer io.Writer, value InputListValue) error {
	err := encodeLength(writer, uint32(len(value.Items)))
	if err != nil {
		return err
	}

	return c.encodeItems(writer, value)
}

func (c *codecForList) encodeTopLevel(writer io.Writer, value InputListValue) error {
	return c.encodeItems(writer, value)
}

func (c *codecForList) decodeNested(reader io.Reader, value *OutputListValue) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode list: item creator is nil")
	}

	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	value.Items = make([]any, 0, length)

	for i := uint32(0); i < length; i++ {
		err := c.decodeItem(reader, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *codecForList) decodeTopLevel(data []byte, value *OutputListValue) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode list: item creator is nil")
	}

	reader := bytes.NewReader(data)
	value.Items = make([]any, 0)

	for reader.Len() > 0 {
		err := c.decodeItem(reader, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *codecForList) encodeItems(writer io.Writer, value InputListValue) error {
	for _, item := range value.Items {
		err := c.generalCodec.doEncodeNested(writer, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *codecForList) decodeItem(reader io.Reader, value *OutputListValue) error {
	newItem := value.ItemCreator()

	err := c.generalCodec.doDecodeNested(reader, newItem)
	if err != nil {
		return err
	}

	value.Items = append(value.Items, newItem)
	return nil
}
