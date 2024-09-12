package abi

import (
	"bytes"
	"errors"
	"io"
)

var errNilArrayItemCreator = errors.New("array item creator is nil")

// ArrayValue is an array of values
type ArrayValue struct {
	Size        uint32
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

// DecodeNested decodes the value from the nested form
func (value *ArrayValue) DecodeNested(reader io.Reader) error {
	return value.decodeItems(reader)
}

// DecodeTopLevel decodes the value from the top-level form
func (value *ArrayValue) DecodeTopLevel(data []byte) error {
	return value.decodeItems(bytes.NewReader(data))
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

func (value *ArrayValue) decodeItems(reader io.Reader) error {
	items := make([]SingleValue, 0, value.Size)
	for i := uint32(0); i < value.Size; i++ {
		item, err := value.decodeItem(reader)
		if err != nil {
			return err
		}

		items = append(items, item)
	}

	value.Items = items
	return nil
}

func (value *ArrayValue) decodeItem(reader io.Reader) (SingleValue, error) {
	if value.ItemCreator == nil {
		return nil, errNilArrayItemCreator
	}

	item := value.ItemCreator()
	err := item.DecodeNested(reader)
	if err != nil {
		return nil, err
	}

	return item, nil
}
