package abi

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCodec(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		codec, err := newCodec(argsNewCodec{
			pubKeyLength: 32,
		})

		require.NoError(t, err)
		require.NotNil(t, codec)
	})

	t.Run("should err if bad public key length", func(t *testing.T) {
		_, err := newCodec(argsNewCodec{
			pubKeyLength: 0,
		})

		require.ErrorContains(t, err, "bad public key length")
	})
}

func TestCodec_EncodeNested(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	doTest := func(t *testing.T, value any, expected string) {
		encoded, err := codec.EncodeNested(value)
		require.NoError(t, err)
		require.Equal(t, expected, hex.EncodeToString(encoded))
	}

	t.Run("address", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1")
		doTest(t, AddressValue{Value: data}, "0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1")
	})

	t.Run("address (bad)", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d42")
		_, err := codec.EncodeNested(AddressValue{Value: data})
		require.ErrorContains(t, err, "public key (address) has invalid length")
	})

	t.Run("struct", func(t *testing.T) {
		fooStruct := StructValue{
			Fields: []Field{
				{
					Value: U8Value{Value: 0x01},
				},
				{
					Value: U16Value{Value: 0x4142},
				},
			},
		}

		doTest(t, fooStruct, "014142")
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 0,
		}

		doTest(t, fooEnum, "00")
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 42,
		}

		doTest(t, fooEnum, "2a")
	})

	t.Run("enum with Fields", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 42,
			Fields: []Field{
				{
					Value: U8Value{Value: 0x01},
				},
				{
					Value: U16Value{Value: 0x4142},
				},
			},
		}

		doTest(t, fooEnum, "2a014142")
	})

	t.Run("option with value", func(t *testing.T) {
		fooOption := OptionValue{
			Value: U16Value{Value: 0x08},
		}

		doTest(t, fooOption, "010008")
	})

	t.Run("option without value", func(t *testing.T) {
		fooOption := OptionValue{
			Value: nil,
		}

		doTest(t, fooOption, "00")
	})

	t.Run("list", func(t *testing.T) {
		fooList := InputListValue{
			Items: []any{
				U16Value{Value: 1},
				U16Value{Value: 2},
				U16Value{Value: 3},
			},
		}

		doTest(t, fooList, "00000003000100020003")
	})

	t.Run("should err when unknown type", func(t *testing.T) {
		type dummy struct {
			foobar string
		}

		encoded, err := codec.EncodeNested(&dummy{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for nested encoding: *abi.dummy")
		require.Nil(t, encoded)
	})
}

func TestCodec_EncodeTopLevel(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	doTest := func(t *testing.T, value any, expected string) {
		encoded, err := codec.EncodeTopLevel(value)
		require.NoError(t, err)
		require.Equal(t, expected, hex.EncodeToString(encoded))
	}

	t.Run("address", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1")
		doTest(t, AddressValue{Value: data}, "0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1")
	})

	t.Run("address (bad)", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d42")
		_, err := codec.EncodeTopLevel(AddressValue{Value: data})
		require.ErrorContains(t, err, "public key (address) has invalid length")
	})

	t.Run("struct", func(t *testing.T) {
		fooStruct := StructValue{
			Fields: []Field{
				{
					Value: U8Value{Value: 0x01},
				},
				{
					Value: U16Value{Value: 0x4142},
				},
			},
		}

		doTest(t, fooStruct, "014142")
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 0,
		}

		doTest(t, fooEnum, "")
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 42,
		}

		doTest(t, fooEnum, "2a")
	})

	t.Run("enum with Fields", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 42,
			Fields: []Field{
				{
					Value: U8Value{Value: 0x01},
				},
				{
					Value: U16Value{Value: 0x4142},
				},
			},
		}

		doTest(t, fooEnum, "2a014142")
	})

	t.Run("should err when unknown type", func(t *testing.T) {
		type dummy struct {
			foobar string
		}

		encoded, err := codec.EncodeTopLevel(&dummy{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for top-level encoding: *abi.dummy")
		require.Nil(t, encoded)
	})
}

func TestCodec_DecodeNested(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	t.Run("u8", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &U8Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U8Value{Value: 0x01}, destination)
	})

	t.Run("i8", func(t *testing.T) {
		data, _ := hex.DecodeString("ff")
		destination := &I8Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I8Value{Value: -1}, destination)
	})

	t.Run("u16", func(t *testing.T) {
		data, _ := hex.DecodeString("4142")
		destination := &U16Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U16Value{Value: 0x4142}, destination)
	})

	t.Run("i16", func(t *testing.T) {
		data, _ := hex.DecodeString("ffff")
		destination := &I16Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I16Value{Value: -1}, destination)
	})

	t.Run("u32", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344")
		destination := &U32Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U32Value{Value: 0x41424344}, destination)
	})

	t.Run("i32", func(t *testing.T) {
		data, _ := hex.DecodeString("ffffffff")
		destination := &I32Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I32Value{Value: -1}, destination)
	})

	t.Run("u64", func(t *testing.T) {
		data, _ := hex.DecodeString("4142434445464748")
		destination := &U64Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U64Value{Value: 0x4142434445464748}, destination)
	})

	t.Run("i64", func(t *testing.T) {
		data, _ := hex.DecodeString("ffffffffffffffff")
		destination := &I64Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I64Value{Value: -1}, destination)
	})

	t.Run("u16, should err because it cannot read 2 bytes", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &U16Value{}

		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot read exactly 2 bytes")
	})

	t.Run("u32, should err because it cannot read 4 bytes", func(t *testing.T) {
		data, _ := hex.DecodeString("4142")
		destination := &U32Value{}

		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot read exactly 4 bytes")
	})

	t.Run("u64, should err because it cannot read 8 bytes", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344")
		destination := &U64Value{}

		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot read exactly 8 bytes")
	})

	t.Run("bigInt", func(t *testing.T) {
		data, _ := hex.DecodeString("00000000")
		destination := &BigIntValue{}
		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &BigIntValue{Value: big.NewInt(0)}, destination)

		data, _ = hex.DecodeString("0000000101")
		destination = &BigIntValue{}
		err = codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &BigIntValue{Value: big.NewInt(1)}, destination)

		data, _ = hex.DecodeString("00000001ff")
		destination = &BigIntValue{}
		err = codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &BigIntValue{Value: big.NewInt(-1)}, destination)
	})

	t.Run("bigInt: should err when bad data", func(t *testing.T) {
		data, _ := hex.DecodeString("0000000301")
		destination := &BigIntValue{}
		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot decode (nested) *abi.BigIntValue, because of: cannot read exactly 3 bytes")
	})

	t.Run("address", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1")

		destination := &AddressValue{}
		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &AddressValue{Value: data}, destination)
	})

	t.Run("address (bad)", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d42")

		destination := &AddressValue{}
		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot read exactly 32 bytes")
	})

	t.Run("struct", func(t *testing.T) {
		data, _ := hex.DecodeString("014142")

		destination := &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		data, _ := hex.DecodeString("00")
		destination := &EnumValue{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x00,
		}, destination)
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &EnumValue{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
		}, destination)
	})

	t.Run("enum with Fields", func(t *testing.T) {
		data, _ := hex.DecodeString("01014142")

		destination := &EnumValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})

	t.Run("option with value", func(t *testing.T) {
		data, _ := hex.DecodeString("010008")

		destination := &OptionValue{
			Value: &U16Value{},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &OptionValue{
			Value: &U16Value{Value: 8},
		}, destination)
	})

	t.Run("option without value", func(t *testing.T) {
		data, _ := hex.DecodeString("00")

		destination := &OptionValue{
			Value: &U16Value{},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &OptionValue{
			Value: nil,
		}, destination)
	})

	t.Run("list", func(t *testing.T) {
		data, _ := hex.DecodeString("00000003000100020003")

		destination := &OutputListValue{
			ItemCreator: func() any { return &U16Value{} },
			Items:       []any{},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t,
			[]any{
				&U16Value{Value: 1},
				&U16Value{Value: 2},
				&U16Value{Value: 3},
			}, destination.Items)
	})

	t.Run("should err when unknown type", func(t *testing.T) {
		type dummy struct {
			foobar string
		}

		err := codec.DecodeNested([]byte{0x00}, &dummy{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for nested decoding: *abi.dummy")
	})
}

func TestCodec_DecodeTopLevel(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	t.Run("u8", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &U8Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U8Value{Value: 0x01}, destination)
	})

	t.Run("i8", func(t *testing.T) {
		data, _ := hex.DecodeString("ff")
		destination := &I8Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I8Value{Value: -1}, destination)
	})

	t.Run("u16", func(t *testing.T) {
		data, _ := hex.DecodeString("02")
		destination := &U16Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U16Value{Value: 0x0002}, destination)
	})

	t.Run("i16", func(t *testing.T) {
		data, _ := hex.DecodeString("ffff")
		destination := &I16Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I16Value{Value: -1}, destination)
	})

	t.Run("u32", func(t *testing.T) {
		data, _ := hex.DecodeString("03")
		destination := &U32Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U32Value{Value: 0x00000003}, destination)
	})

	t.Run("i32", func(t *testing.T) {
		data, _ := hex.DecodeString("ffffffff")
		destination := &I32Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I32Value{Value: -1}, destination)
	})

	t.Run("u64", func(t *testing.T) {
		data, _ := hex.DecodeString("04")
		destination := &U64Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U64Value{Value: 0x0000000000000004}, destination)
	})

	t.Run("i64", func(t *testing.T) {
		data, _ := hex.DecodeString("ffffffffffffffff")
		destination := &I64Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &I64Value{Value: -1}, destination)
	})

	t.Run("u8, i8: should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("4142")

		err := codec.DecodeTopLevel(data, &U8Value{})
		require.ErrorContains(t, err, "decoded value is too large")

		err = codec.DecodeTopLevel(data, &I8Value{})
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u16, i16: should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344")

		err := codec.DecodeTopLevel(data, &U16Value{})
		require.ErrorContains(t, err, "decoded value is too large")

		err = codec.DecodeTopLevel(data, &I16Value{})
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u32, i32: should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("4142434445464748")

		err := codec.DecodeTopLevel(data, &U32Value{})
		require.ErrorContains(t, err, "decoded value is too large")

		err = codec.DecodeTopLevel(data, &I32Value{})
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u64, i64: should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344454647489876")

		err := codec.DecodeTopLevel(data, &U64Value{})
		require.ErrorContains(t, err, "decoded value is too large")

		err = codec.DecodeTopLevel(data, &I64Value{})
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("bigInt", func(t *testing.T) {
		data, _ := hex.DecodeString("")
		destination := &BigIntValue{}
		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &BigIntValue{Value: big.NewInt(0)}, destination)

		data, _ = hex.DecodeString("01")
		destination = &BigIntValue{}
		err = codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &BigIntValue{Value: big.NewInt(1)}, destination)

		data, _ = hex.DecodeString("ff")
		destination = &BigIntValue{}
		err = codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &BigIntValue{Value: big.NewInt(-1)}, destination)
	})

	t.Run("address", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1")

		destination := &AddressValue{}
		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &AddressValue{Value: data}, destination)
	})

	t.Run("address (bad)", func(t *testing.T) {
		data, _ := hex.DecodeString("0139472eff6886771a982f3083da5d42")

		destination := &AddressValue{}
		err := codec.DecodeTopLevel(data, destination)
		require.ErrorContains(t, err, "public key (address) has invalid length")
	})

	t.Run("struct", func(t *testing.T) {
		data, _ := hex.DecodeString("014142")

		destination := &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		data, _ := hex.DecodeString("")
		destination := &EnumValue{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x00,
		}, destination)
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &EnumValue{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
		}, destination)
	})

	t.Run("enum with Fields", func(t *testing.T) {
		data, _ := hex.DecodeString("01014142")

		destination := &EnumValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})

	t.Run("should err when unknown type", func(t *testing.T) {
		type dummy struct {
			foobar string
		}

		err := codec.DecodeTopLevel([]byte{0x00}, &dummy{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for top-level decoding: *abi.dummy")
	})
}

func testEncodeNested(t *testing.T, codec *codec, value any, expected string) {
	encoded, err := codec.EncodeNested(value)

	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(encoded))
}

func testEncodeTopLevel(t *testing.T, codec *codec, value any, expected string) {
	encoded, err := codec.EncodeTopLevel(value)

	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(encoded))
}

func testDecodeNested(t *testing.T, codec *codec, encodedData string, destination any, expected any) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeNested(data, destination)

	require.NoError(t, err)
	require.Equal(t, expected, destination)
}

func testDecodeTopLevel(t *testing.T, codec *codec, encodedData string, destination any, expected any) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeTopLevel(data, destination)

	require.NoError(t, err)
	require.Equal(t, expected, destination)
}
