package abi

import (
	"errors"
	"io"
)

func (c *codec) encodeNestedList(writer io.Writer, value InputListValue) error {
	err := encodeLength(writer, uint32(len(value.Items)))
	if err != nil {
		return err
	}

	for _, item := range value.Items {
		err := c.doEncodeNested(writer, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *codec) decodeNestedList(reader io.Reader, value *OutputListValue) error {
	if value.ItemCreator == nil {
		return errors.New("cannot deserialize list: item creator is nil")
	}

	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	value.Items = make([]any, 0, length)

	for i := uint32(0); i < length; i++ {
		newItem := value.ItemCreator()

		err := c.doDecodeNested(reader, newItem)
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)
	}

	return nil
}
