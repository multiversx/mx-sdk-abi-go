package abi

import (
	"math/big"
	"testing"
)

func TestCodec_Numerical(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec, U8Value{Value: 0x00}, "00")
		testEncodeNested(t, codec, U8Value{Value: 0x01}, "01")
		testEncodeNested(t, codec, U8Value{Value: 0x42}, "42")
		testEncodeNested(t, codec, U8Value{Value: 0xff}, "ff")

		testEncodeNested(t, codec, I8Value{Value: 0x00}, "00")
		testEncodeNested(t, codec, I8Value{Value: 0x01}, "01")
		testEncodeNested(t, codec, I8Value{Value: -1}, "ff")
		testEncodeNested(t, codec, I8Value{Value: -128}, "80")
		testEncodeNested(t, codec, I8Value{Value: 127}, "7f")

		testEncodeNested(t, codec, U16Value{Value: 0x0000}, "0000")
		testEncodeNested(t, codec, U16Value{Value: 0x0011}, "0011")
		testEncodeNested(t, codec, U16Value{Value: 0x1234}, "1234")
		testEncodeNested(t, codec, U16Value{Value: 0xffff}, "ffff")

		testEncodeNested(t, codec, I16Value{Value: 0x0000}, "0000")
		testEncodeNested(t, codec, I16Value{Value: 0x0011}, "0011")
		testEncodeNested(t, codec, I16Value{Value: -1}, "ffff")
		testEncodeNested(t, codec, I16Value{Value: -32768}, "8000")

		testEncodeNested(t, codec, U32Value{Value: 0x00000000}, "00000000")
		testEncodeNested(t, codec, U32Value{Value: 0x00000011}, "00000011")
		testEncodeNested(t, codec, U32Value{Value: 0x00001122}, "00001122")
		testEncodeNested(t, codec, U32Value{Value: 0x00112233}, "00112233")
		testEncodeNested(t, codec, U32Value{Value: 0x11223344}, "11223344")
		testEncodeNested(t, codec, U32Value{Value: 0xffffffff}, "ffffffff")

		testEncodeNested(t, codec, I32Value{Value: 0x00000000}, "00000000")
		testEncodeNested(t, codec, I32Value{Value: 0x00000011}, "00000011")
		testEncodeNested(t, codec, I32Value{Value: -1}, "ffffffff")
		testEncodeNested(t, codec, I32Value{Value: -2147483648}, "80000000")

		testEncodeNested(t, codec, U64Value{Value: 0x0000000000000000}, "0000000000000000")
		testEncodeNested(t, codec, U64Value{Value: 0x0000000000000011}, "0000000000000011")
		testEncodeNested(t, codec, U64Value{Value: 0x0011223344556677}, "0011223344556677")
		testEncodeNested(t, codec, U64Value{Value: 0xffffffffffffffff}, "ffffffffffffffff")

		testEncodeNested(t, codec, I64Value{Value: 0x0000000000000000}, "0000000000000000")
		testEncodeNested(t, codec, I64Value{Value: 0x0000000000000011}, "0000000000000011")
		testEncodeNested(t, codec, I64Value{Value: -1}, "ffffffffffffffff")
		testEncodeNested(t, codec, I64Value{Value: -9223372036854775808}, "8000000000000000")

		testEncodeNested(t, codec, BigIntValue{Value: big.NewInt(0)}, "00000000")
		testEncodeNested(t, codec, BigIntValue{Value: big.NewInt(1)}, "0000000101")
		testEncodeNested(t, codec, BigIntValue{Value: big.NewInt(-1)}, "00000001ff")
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec, U8Value{Value: 0x00}, "")
		testEncodeTopLevel(t, codec, U8Value{Value: 0x01}, "01")

		testEncodeTopLevel(t, codec, I8Value{Value: 0x00}, "")
		testEncodeTopLevel(t, codec, I8Value{Value: 0x01}, "01")
		testEncodeTopLevel(t, codec, I8Value{Value: -1}, "ff")

		testEncodeTopLevel(t, codec, U16Value{Value: 0x0000}, "")
		testEncodeTopLevel(t, codec, U16Value{Value: 0x0011}, "11")

		testEncodeTopLevel(t, codec, I16Value{Value: 0x0000}, "")
		testEncodeTopLevel(t, codec, I16Value{Value: 0x0011}, "11")
		testEncodeTopLevel(t, codec, I16Value{Value: -1}, "ff")

		testEncodeTopLevel(t, codec, U32Value{Value: 0x00004242}, "4242")

		testEncodeTopLevel(t, codec, I32Value{Value: 0x00000000}, "")
		testEncodeTopLevel(t, codec, I32Value{Value: 0x00000011}, "11")
		testEncodeTopLevel(t, codec, I32Value{Value: -1}, "ff")

		testEncodeTopLevel(t, codec, U64Value{Value: 0x0042434445464748}, "42434445464748")

		testEncodeTopLevel(t, codec, I64Value{Value: 0x0000000000000000}, "")
		testEncodeTopLevel(t, codec, I64Value{Value: 0x0000000000000011}, "11")
		testEncodeTopLevel(t, codec, I64Value{Value: -1}, "ff")

		testEncodeTopLevel(t, codec, BigIntValue{Value: big.NewInt(0)}, "")
		testEncodeTopLevel(t, codec, BigIntValue{Value: big.NewInt(1)}, "01")
		testEncodeTopLevel(t, codec, BigIntValue{Value: big.NewInt(-1)}, "ff")
	})

	t.Run("should decode nested", func(t *testing.T) {
	})

	t.Run("should decode top-level", func(t *testing.T) {
	})
}
