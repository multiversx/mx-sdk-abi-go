package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodec_Option(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, OptionValue{
			Value: nil,
		}, "00")

		testEncodeNested(t, codec, OptionValue{
			Value: U16Value{Value: 0x08},
		}, "010008")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, OptionValue{
			Value: nil,
		}, "")

		testEncodeTopLevel(t, codec, OptionValue{
			Value: U16Value{Value: 0x08},
		}, "010008")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00",
			&OptionValue{},
			&OptionValue{
				Value: nil,
			},
		)

		testDecodeNested(t, codec, "010008",
			&OptionValue{
				Value: &U16Value{},
			},
			&OptionValue{
				Value: &U16Value{Value: 0x08},
			},
		)
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "",
			&OptionValue{},
			&OptionValue{
				Value: nil,
			},
		)

		testDecodeTopLevel(t, codec, "010008",
			&OptionValue{
				Value: &U16Value{},
			},
			&OptionValue{
				Value: &U16Value{Value: 0x08},
			},
		)
	})

	t.Run("should err on decode top-level (bad marker for value presence)", func(t *testing.T) {
		data, _ := hex.DecodeString("002a")
		err := codec.DecodeTopLevel(data, &OptionValue{})
		require.ErrorContains(t, err, "invalid first byte for top-level encoded option: 0")

		data, _ = hex.DecodeString("072a")
		err = codec.DecodeTopLevel(data, &OptionValue{})
		require.ErrorContains(t, err, "invalid first byte for top-level encoded option: 7")
	})
}
