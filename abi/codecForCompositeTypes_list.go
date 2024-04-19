package abi

import (
	"bytes"
	"errors"
	"io"
)

func (c *codec) encodeNestedList(writer io.Writer, value InputListValue) error {
	err := encodeLength(writer, uint32(len(value.Items)))
	if err != nil {
		return err
	}

	return c.encodeListItems(writer, value)
}

func (c *codec) encodeTopLevelList(writer io.Writer, value InputListValue) error {
	return c.encodeListItems(writer, value)
}

func (c *codec) decodeNestedList(reader io.Reader, value *OutputListValue) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode list: item creator is nil")
	}

	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	value.Items = make([]any, 0, length)

	for i := uint32(0); i < length; i++ {
		err := c.decodeListItem(reader, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *codec) decodeTopLevelList(data []byte, value *OutputListValue) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode list: item creator is nil")
	}

	reader := bytes.NewReader(data)
	value.Items = make([]any, 0)

	for reader.Len() > 0 {
		err := c.decodeListItem(reader, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *codec) encodeListItems(writer io.Writer, value InputListValue) error {
	for _, item := range value.Items {
		err := c.doEncodeNested(writer, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *codec) decodeListItem(reader io.Reader, value *OutputListValue) error {
	newItem := value.ItemCreator()

	err := c.doDecodeNested(reader, newItem)
	if err != nil {
		return err
	}

	value.Items = append(value.Items, newItem)
	return nil
}
