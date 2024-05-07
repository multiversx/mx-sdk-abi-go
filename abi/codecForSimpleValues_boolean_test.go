package abi

import (
	"testing"
)

func TestCodecForBool(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, BoolValue{Value: false}, "00")
		testEncodeNested(t, codec, BoolValue{Value: true}, "01")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, BoolValue{Value: false}, "")
		testEncodeTopLevel(t, codec, BoolValue{Value: true}, "01")
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec, "00", &BoolValue{}, &BoolValue{Value: false})
		testDecodeNested(t, codec, "01", &BoolValue{}, &BoolValue{Value: true})
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec, "", &BoolValue{}, &BoolValue{Value: false})
		testDecodeTopLevel(t, codec, "01", &BoolValue{}, &BoolValue{Value: true})
	})
}
