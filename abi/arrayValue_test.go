package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayValue(t *testing.T) {
	codec := &codec{}

	t.Run("should encode nested", func(t *testing.T) {
		testEncodeNested(t, codec,
			&ArrayValue{
				Length: 3,
				Items: []SingleValue{
					&U16Value{Value: 1},
					&U16Value{Value: 2},
					&U16Value{Value: 3},
				},
			},
			"000100020003",
		)
	})

	t.Run("should encode top-level", func(t *testing.T) {
		testEncodeTopLevel(t, codec,
			&ArrayValue{
				Length: 3,
				Items: []SingleValue{
					&U16Value{Value: 1},
					&U16Value{Value: 2},
					&U16Value{Value: 3},
				},
			},
			"000100020003",
		)
	})

	t.Run("should decode nested", func(t *testing.T) {
		data, _ := hex.DecodeString("000100020003")

		destination := &ArrayValue{
			Length:      3,
			ItemCreator: func() SingleValue { return &U16Value{} },
			Items:       []SingleValue{},
		}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t,
			[]SingleValue{
				&U16Value{Value: 1},
				&U16Value{Value: 2},
				&U16Value{Value: 3},
			},
			destination.Items,
		)
	})

	t.Run("should decode top-level", func(t *testing.T) {
		data, _ := hex.DecodeString("000100020003")

		destination := &ArrayValue{
			Length:      3,
			ItemCreator: func() SingleValue { return &U16Value{} },
			Items:       []SingleValue{},
		}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t,
			[]SingleValue{
				&U16Value{Value: 1},
				&U16Value{Value: 2},
				&U16Value{Value: 3},
			},
			destination.Items,
		)
	})

	t.Run("item creator is nil, should error", func(t *testing.T) {
		data, _ := hex.DecodeString("000100020003")

		destination := &ArrayValue{
			Length: 3,
			Items:  []SingleValue{},
		}

		err := codec.DecodeTopLevel(data, destination)
		require.ErrorContains(t, err, "item creator is nil")
	})
}
