package abi

import (
	"testing"
)

func TestEnumValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec,
			&EnumValue{
				Discriminant: 0,
			},
			"00",
		)

		testEncodeNested(t, codec,
			&EnumValue{
				Discriminant: 42,
			},
			"2a",
		)

		testEncodeNested(t, codec,
			&EnumValue{
				Discriminant: 42,
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
			"2a014142",
		)
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec,
			&EnumValue{
				Discriminant: 0,
			},
			"",
		)

		testEncodeTopLevel(t, codec,
			&EnumValue{
				Discriminant: 42,
			},
			"2a",
		)

		testEncodeTopLevel(t, codec,
			&EnumValue{
				Discriminant: 42,
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
			"2a014142",
		)
	})

	t.Run("should decode nested", func(t *testing.T) {
		testDecodeNested(t, codec,
			"00",
			&EnumValue{},
			&EnumValue{
				Discriminant: 0x00,
			},
		)

		testDecodeNested(t, codec,
			"2a",
			&EnumValue{},
			&EnumValue{
				Discriminant: 42,
			},
		)

		testDecodeNested(t, codec,
			"01014142",
			&EnumValue{
				Fields: []Field{
					{
						Value: &U8Value{},
					},
					{
						Value: &U16Value{},
					},
				},
			},
			&EnumValue{
				Discriminant: 0x01,
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
		)
	})

	t.Run("should decode top-level", func(t *testing.T) {
		testDecodeTopLevel(t, codec,
			"",
			&EnumValue{},
			&EnumValue{
				Discriminant: 0x00,
			},
		)

		testDecodeTopLevel(t, codec,
			"2a",
			&EnumValue{},
			&EnumValue{
				Discriminant: 42,
			},
		)

		testDecodeTopLevel(t, codec,
			"01014142",
			&EnumValue{
				Fields: []Field{
					{
						Value: &U8Value{},
					},
					{
						Value: &U16Value{},
					},
				},
			},
			&EnumValue{
				Discriminant: 0x01,
				Fields: []Field{
					{
						Value: &U8Value{Value: 0x01},
					},
					{
						Value: &U16Value{Value: 0x4142},
					},
				},
			},
		)
	})
}
