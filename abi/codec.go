package abi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
)

type codec struct {
	codecForSmallInt *codecForSmallInt
	codecForStruct   *codeForStruct
	codecForEnum     *codecForEnum
	codecForOption   *codecForOption
	codecForList     *codecForList
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

	codec := &codec{
		codecForSmallInt: &codecForSmallInt{},
	}

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
		return c.codecForSmallInt.encodeNested(writer, value.Value, 1)
	case U16Value:
		return c.codecForSmallInt.encodeNested(writer, value.Value, 2)
	case U32Value:
		return c.codecForSmallInt.encodeNested(writer, value.Value, 4)
	case U64Value:
		return c.codecForSmallInt.encodeNested(writer, value.Value, 8)
	case I8Value:
		return c.codecForSmallInt.encodeNested(writer, value.Value, 1)
	case I16Value:
		return c.codecForSmallInt.encodeNested(writer, value.Value, 2)
	case I32Value:
		return c.codecForSmallInt.encodeNested(writer, value.Value, 4)
	case I64Value:
		return c.codecForSmallInt.encodeNested(writer, value.Value, 8)
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
		return c.codecForSmallInt.encodeTopLevelUnsigned(writer, uint64(value.Value))
	case U16Value:
		return c.codecForSmallInt.encodeTopLevelUnsigned(writer, uint64(value.Value))
	case U32Value:
		return c.codecForSmallInt.encodeTopLevelUnsigned(writer, uint64(value.Value))
	case U64Value:
		return c.codecForSmallInt.encodeTopLevelUnsigned(writer, value.Value)
	case I8Value:
		return c.codecForSmallInt.encodeTopLevelSigned(writer, int64(value.Value))
	case I16Value:
		return c.codecForSmallInt.encodeTopLevelSigned(writer, int64(value.Value))
	case I32Value:
		return c.codecForSmallInt.encodeTopLevelSigned(writer, int64(value.Value))
	case I64Value:
		return c.codecForSmallInt.encodeTopLevelSigned(writer, value.Value)
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
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 1)
	case *U16Value:
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 2)
	case *U32Value:
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 4)
	case *U64Value:
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 8)
	case *I8Value:
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 1)
	case *I16Value:
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 2)
	case *I32Value:
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 4)
	case *I64Value:
		return c.codecForSmallInt.decodeNested(reader, &value.Value, 8)
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
		n, err := c.codecForSmallInt.decodeTopLevelUnsigned(data, math.MaxUint8)
		if err != nil {
			return err
		}

		value.Value = uint8(n)
	case *U16Value:
		n, err := c.codecForSmallInt.decodeTopLevelUnsigned(data, math.MaxUint16)
		if err != nil {
			return err
		}

		value.Value = uint16(n)
	case *U32Value:
		n, err := c.codecForSmallInt.decodeTopLevelUnsigned(data, math.MaxUint32)
		if err != nil {
			return err
		}

		value.Value = uint32(n)
	case *U64Value:
		n, err := c.codecForSmallInt.decodeTopLevelUnsigned(data, math.MaxUint64)
		if err != nil {
			return err
		}

		value.Value = uint64(n)
	case *I8Value:
		n, err := c.codecForSmallInt.decodeTopLevelSigned(data, math.MaxInt8)
		if err != nil {
			return err
		}

		value.Value = int8(n)
	case *I16Value:
		n, err := c.codecForSmallInt.decodeTopLevelSigned(data, math.MaxInt16)
		if err != nil {
			return err
		}

		value.Value = int16(n)
	case *I32Value:
		n, err := c.codecForSmallInt.decodeTopLevelSigned(data, math.MaxInt32)
		if err != nil {
			return err
		}

		value.Value = int32(n)

	case *I64Value:
		n, err := c.codecForSmallInt.decodeTopLevelSigned(data, math.MaxInt64)
		if err != nil {
			return err
		}

		value.Value = int64(n)
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
