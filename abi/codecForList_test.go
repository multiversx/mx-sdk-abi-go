package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodecForList(t *testing.T) {
	codec, _ := newCodec(argsNewCodec{
		pubKeyLength: 32,
	})

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec,
			InputListValue{
				Items: []any{
					U16Value{Value: 1},
					U16Value{Value: 2},
					U16Value{Value: 3},
				},
			},
			"00000003000100020003",
		)
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec,
			InputListValue{
				Items: []any{
					U16Value{Value: 1},
					U16Value{Value: 2},
					U16Value{Value: 3},
				},
			},
			"000100020003",
		)
	})

	t.Run("should decode nested", func(t *testing.T) {
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
			},
			destination.Items,
		)
	})

	t.Run("should decode top-level", func(t *testing.T) {
		data, _ := hex.DecodeString("000100020003")

		destination := &OutputListValue{
			ItemCreator: func() any { return &U16Value{} },
			Items:       []any{},
		}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t,
			[]any{
				&U16Value{Value: 1},
				&U16Value{Value: 2},
				&U16Value{Value: 3},
			},
			destination.Items,
		)
	})
}