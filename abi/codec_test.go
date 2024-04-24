package abi

import (
	"encoding/hex"
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

func testDecodeNestedWithError(t *testing.T, codec *codec, encodedData string, destination any, expectedError string) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeNested(data, destination)

	require.ErrorContains(t, err, expectedError)
}

func testDecodeTopLevel(t *testing.T, codec *codec, encodedData string, destination any, expected any) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeTopLevel(data, destination)

	require.NoError(t, err)
	require.Equal(t, expected, destination)
}

func testDecodeTopLevelWithError(t *testing.T, codec *codec, encodedData string, destination any, expectedError string) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeTopLevel(data, destination)

	require.ErrorContains(t, err, expectedError)
}
