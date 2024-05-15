package abi

import (
	"bytes"
	"errors"
	"io"
)

// ListValue is a list of values
type ListValue struct {
	Items       []singleValue
	ItemCreator func() singleValue
}

func (value *ListValue) encodeNested(writer io.Writer) error {
	err := encodeLength(writer, uint32(len(value.Items)))
	if err != nil {
		return err
	}

	return value.encodeItems(writer)
}

func (value *ListValue) encodeTopLevel(writer io.Writer) error {
	return value.encodeItems(writer)
}

func (value *ListValue) encodeItems(writer io.Writer) error {
	for _, item := range value.Items {
		err := item.encodeNested(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (value *ListValue) decodeNested(reader io.Reader) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode list: item creator is nil")
	}

	length, err := decodeLength(reader)
	if err != nil {
		return err
	}

	value.Items = make([]singleValue, 0, length)

	for i := uint32(0); i < length; i++ {
		err := value.decodeItem(reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func (value *ListValue) decodeTopLevel(data []byte) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode list: item creator is nil")
	}

	reader := bytes.NewReader(data)
	value.Items = make([]singleValue, 0)

	for reader.Len() > 0 {
		err := value.decodeItem(reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func (value *ListValue) decodeItem(reader io.Reader) error {
	newItem := value.ItemCreator()

	err := newItem.decodeNested(reader)
	if err != nil {
		return err
	}

	value.Items = append(value.Items, newItem)
	return nil
}
