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

	t.Run("bigInt: should err when bad data", func(t *testing.T) {
		data, _ := hex.DecodeString("0000000301")
		destination := &BigIntValue{}
		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot decode (nested) *abi.BigIntValue, because of: cannot read exactly 3 bytes")
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
