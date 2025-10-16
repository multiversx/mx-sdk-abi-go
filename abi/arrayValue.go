package abi

import (
	"bytes"
	"errors"
	"io"
)

// ArrayValue is an array of values
type ArrayValue struct {
	Length      uint32
	Items       []SingleValue
	ItemCreator func() SingleValue
}

// EncodeNested encodes the value in the nested form
func (value *ArrayValue) EncodeNested(writer io.Writer) error {
	return value.encodeItems(writer)
}

// EncodeTopLevel encodes the value in the top-level form
func (value *ArrayValue) EncodeTopLevel(writer io.Writer) error {
	return value.encodeItems(writer)
}

func (value *ArrayValue) encodeItems(writer io.Writer) error {
	for _, item := range value.Items {
		err := item.EncodeNested(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

// DecodeNested decodes the value from the nested form
func (value *ArrayValue) DecodeNested(reader io.Reader) error {
	value.Items = make([]SingleValue, 0, value.Length)

	for i := uint32(0); i < value.Length; i++ {
		err := value.decodeItem(reader)
		if err != nil {
			return err
		}
	}

	return nil
}

// DecodeTopLevel decodes the value from the top-level form
func (value *ArrayValue) DecodeTopLevel(data []byte) error {
	reader := bytes.NewReader(data)
	value.Items = make([]SingleValue, 0)

	for reader.Len() > 0 {
		err := value.decodeItem(reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func (value *ArrayValue) decodeItem(reader io.Reader) error {
	if value.ItemCreator == nil {
		return errors.New("cannot decode array: item creator is nil")
	}

	newItem := value.ItemCreator()

	err := newItem.DecodeNested(reader)
	if err != nil {
		return err
	}

	value.Items = append(value.Items, newItem)
	return nil
}
