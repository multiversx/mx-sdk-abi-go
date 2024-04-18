package abi

import (
	"testing"
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
	})
}
