package abi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type codec struct {
	codecForStruct *codeForStruct
	codecForEnum   *codecForEnum
	codecForOption *codecForOption
	codecForList   *codecForList
}

// argsNewCodec defines the arguments needed for a new codec
type argsNewCodec struct {
	pubKeyLength int
}

// newCodec creates a new default codec which follows the rules of the MultiversX Serialization format:
// https://docs.multiversx.com/developers/data/serialization-overview
func newCodec(args argsNewCodec) (*codec, error) {
	if args.pubKeyLength <= 0 {
		return nil, errors.New("cannot create codec: bad public key length")
	}

	codec := &codec{}

	codec.codecForStruct = &codeForStruct{
		generalCodec: codec,
	}

	codec.codecForEnum = &codecForEnum{
		generalCodec: codec,
	}

	codec.codecForOption = &codecForOption{
		generalCodec: codec,
	}

	codec.codecForList = &codecForList{
		generalCodec: codec,
	}

	return codec, nil
}

// EncodeNested encodes the given value following the nested encoding rules
func (c *codec) EncodeNested(value any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := c.doEncodeNested(buffer, value)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *codec) doEncodeNested(writer io.Writer, value any) error {
	switch value := value.(type) {
	case BoolValue:
		return value.encodeNested(writer)
	case U8Value:
		return value.encodeNested(writer)
	case U16Value:
		return value.encodeNested(writer)
	case U32Value:
		return value.encodeNested(writer)
	case U64Value:
		return value.encodeNested(writer)
	case I8Value:
		return value.encodeNested(writer)
	case I16Value:
		return value.encodeNested(writer)
	case I32Value:
		return value.encodeNested(writer)
	case I64Value:
		return value.encodeNested(writer)
	case BigUIntValue:
		return value.encodeNested(writer)
	case BigIntValue:
		return value.encodeNested(writer)
	case AddressValue:
		return value.encodeNested(writer)
	case StringValue:
		return value.encodeNested(writer)
	case BytesValue:
		return value.encodeNested(writer)
	case StructValue:
		return c.codecForStruct.encodeNested(writer, value)
	case EnumValue:
		return c.codecForEnum.encodeNested(writer, value)
	case OptionValue:
		return c.codecForOption.encodeNested(writer, value)
	case InputListValue:
		return c.codecForList.encodeNested(writer, value)
	default:
		return fmt.Errorf("unsupported type for nested encoding: %T", value)
	}
}

// EncodeTopLevel encodes the given value following the top-level encoding rules
func (c *codec) EncodeTopLevel(value any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := c.doEncodeTopLevel(buffer, value)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *codec) doEncodeTopLevel(writer io.Writer, value any) error {
	switch value := value.(type) {
	case BoolValue:
		return value.encodeTopLevel(writer)
	case U8Value:
		return value.encodeTopLevel(writer)
	case U16Value:
		return value.encodeTopLevel(writer)
	case U32Value:
		return value.encodeTopLevel(writer)
	case U64Value:
		return value.encodeTopLevel(writer)
	case I8Value:
		return value.encodeTopLevel(writer)
	case I16Value:
		return value.encodeTopLevel(writer)
	case I32Value:
		return value.encodeTopLevel(writer)
	case I64Value:
		return value.encodeTopLevel(writer)
	case BigUIntValue:
		return value.encodeTopLevel(writer)
	case BigIntValue:
		return value.encodeTopLevel(writer)
	case AddressValue:
		return value.encodeTopLevel(writer)
	case StringValue:
		return value.encodeTopLevel(writer)
	case BytesValue:
		return value.encodeTopLevel(writer)
	case StructValue:
		return c.codecForStruct.encodeTopLevel(writer, value)
	case EnumValue:
		return c.codecForEnum.encodeTopLevel(writer, value)
	case OptionValue:
		return c.codecForOption.encodeTopLevel(writer, value)
	case InputListValue:
		return c.codecForList.encodeTopLevel(writer, value)
	default:
		return fmt.Errorf("unsupported type for top-level encoding: %T", value)
	}
}

// DecodeNested decodes the given data into the provided object following the nested decoding rules
func (c *codec) DecodeNested(data []byte, value any) error {
	reader := bytes.NewReader(data)
	err := c.doDecodeNested(reader, value)
	if err != nil {
		return fmt.Errorf("cannot decode (nested) %T, because of: %w", value, err)
	}

	return nil
}

func (c *codec) doDecodeNested(reader io.Reader, value any) error {
	switch value := value.(type) {
	case *BoolValue:
		return value.decodeNested(reader)
	case *U8Value:
		return value.decodeNested(reader)
	case *U16Value:
		return value.decodeNested(reader)
	case *U32Value:
		return value.decodeNested(reader)
	case *U64Value:
		return value.decodeNested(reader)
	case *I8Value:
		return value.decodeNested(reader)
	case *I16Value:
		return value.decodeNested(reader)
	case *I32Value:
		return value.decodeNested(reader)
	case *I64Value:
		return value.decodeNested(reader)
	case *BigUIntValue:
		return value.decodeNested(reader)
	case *BigIntValue:
		return value.decodeNested(reader)
	case *AddressValue:
		return value.decodeNested(reader)
	case *StringValue:
		return value.decodeNested(reader)
	case *BytesValue:
		return value.decodeNested(reader)
	case *StructValue:
		return c.codecForStruct.decodeNested(reader, value)
	case *EnumValue:
		return c.codecForEnum.decodeNested(reader, value)
	case *OptionValue:
		return c.codecForOption.decodeNested(reader, value)
	case *OutputListValue:
		return c.codecForList.decodeNested(reader, value)
	default:
		return fmt.Errorf("unsupported type for nested decoding: %T", value)
	}
}

// DecodeTopLevel decodes the given data into the provided object following the top-level decoding rules
func (c *codec) DecodeTopLevel(data []byte, value any) error {
	err := c.doDecodeTopLevel(data, value)
	if err != nil {
		return fmt.Errorf("cannot decode (top-level) %T, because of: %w", value, err)
	}

	return nil
}

func (c *codec) doDecodeTopLevel(data []byte, value any) error {
	switch value := value.(type) {
	case *BoolValue:
		return value.decodeTopLevel(data)
	case *U8Value:
		return value.decodeTopLevel(data)
	case *U16Value:
		return value.decodeTopLevel(data)
	case *U32Value:
		return value.decodeTopLevel(data)
	case *U64Value:
		return value.decodeTopLevel(data)
	case *I8Value:
		return value.decodeTopLevel(data)
	case *I16Value:
		return value.decodeTopLevel(data)
	case *I32Value:
		return value.decodeTopLevel(data)
	case *I64Value:
		return value.decodeTopLevel(data)
	case *BigUIntValue:
		value.decodeTopLevel(data)
	case *BigIntValue:
		value.decodeTopLevel(data)
	case *AddressValue:
		return value.decodeTopLevel(data)
	case *StringValue:
		return value.decodeTopLevel(data)
	case *BytesValue:
		return value.decodeTopLevel(data)
	case *StructValue:
		return c.codecForStruct.decodeTopLevel(data, value)
	case *EnumValue:
		return c.codecForEnum.decodeTopLevel(data, value)
	case *OptionValue:
		return c.codecForOption.decodeTopLevel(data, value)
	case *OutputListValue:
		return c.codecForList.decodeTopLevel(data, value)
	default:
		return fmt.Errorf("unsupported type for top-level decoding: %T", value)
	}

	return nil
}
