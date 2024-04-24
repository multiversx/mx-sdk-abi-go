package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodec(t *testing.T) {
	t.Run("should create new codec", func(t *testing.T) {
		codec, err := newCodec(argsNewCodec{
			pubKeyLength: 32,
		})

		require.NoError(t, err)
		require.NotNil(t, codec)
	})

	t.Run("should err on creating new codec when public key length is bad", func(t *testing.T) {
		_, err := newCodec(argsNewCodec{
			pubKeyLength: 0,
		})

		require.ErrorContains(t, err, "bad public key length")
	})

	t.Run("should err when encoding or decoding an unknown type", func(t *testing.T) {
		codec, _ := newCodec(argsNewCodec{
			pubKeyLength: 32,
		})

		type dummyType struct {
			foobar string
		}

		encoded, err := codec.EncodeNested(&dummyType{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for nested encoding: *abi.dummy")
		require.Nil(t, encoded)

		encoded, err = codec.EncodeTopLevel(&dummyType{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for top-level encoding: *abi.dummy")
		require.Nil(t, encoded)

		err = codec.DecodeNested([]byte{0x00}, &dummyType{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for nested decoding: *abi.dummy")

		err = codec.DecodeTopLevel([]byte{0x00}, &dummyType{foobar: "hello"})
		require.ErrorContains(t, err, "unsupported type for top-level decoding: *abi.dummy")
	})
}

// testEncodeNested is a helper function to test nested encoding.
func testEncodeNested(t *testing.T, codec *codec, value any, expected string) {
	encoded, err := codec.EncodeNested(value)

	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(encoded))
}

// testEncodeTopLevel is a helper function to test top-level encoding.
func testEncodeTopLevel(t *testing.T, codec *codec, value any, expected string) {
	encoded, err := codec.EncodeTopLevel(value)

	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(encoded))
}

// testDecodeNested is a helper function to test nested decoding.
func testDecodeNested(t *testing.T, codec *codec, encodedData string, destination any, expected any) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeNested(data, destination)

	require.NoError(t, err)
	require.Equal(t, expected, destination)
}

// testDecodeNestedWithError is a helper function to test nested decoding.
func testDecodeNestedWithError(t *testing.T, codec *codec, encodedData string, destination any, expectedError string) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeNested(data, destination)

	require.ErrorContains(t, err, expectedError)
}

// testDecodeTopLevel is a helper function to test top-level decoding.
func testDecodeTopLevel(t *testing.T, codec *codec, encodedData string, destination any, expected any) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeTopLevel(data, destination)

	require.NoError(t, err)
	require.Equal(t, expected, destination)
}

// testDecodeTopLevelWithError is a helper function to test top-level decoding.
func testDecodeTopLevelWithError(t *testing.T, codec *codec, encodedData string, destination any, expectedError string) {
	data, _ := hex.DecodeString(encodedData)
	err := codec.DecodeTopLevel(data, destination)

	require.ErrorContains(t, err, expectedError)
}
