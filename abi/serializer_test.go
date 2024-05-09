package abi

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializer_Serialize(t *testing.T) {
	serializer, err := NewSerializer(ArgsNewSerializer{
		PartsSeparator: "@",
		PubKeyLength:   32,
	})
	require.NoError(t, err)

	t.Run("u8", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			U8Value{Value: 0x42},
		})

		require.NoError(t, err)
		require.Equal(t, "42", data)
	})

	t.Run("u16", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			U16Value{Value: 0x4243},
		})

		require.NoError(t, err)
		require.Equal(t, "4243", data)
	})

	t.Run("u8, u16", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			U8Value{Value: 0x42},
			U16Value{Value: 0x4243},
		})

		require.NoError(t, err)
		require.Equal(t, "42@4243", data)
	})

	t.Run("optional (missing)", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			U8Value{Value: 0x42},
			InputOptionalValue{},
		})

		require.NoError(t, err)
		require.Equal(t, "42", data)
	})

	t.Run("optional (provided)", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			U8Value{Value: 0x42},
			InputOptionalValue{Value: U8Value{Value: 0x43}},
		})

		require.NoError(t, err)
		require.Equal(t, "42@43", data)
	})

	t.Run("optional: should err because optional must be last", func(t *testing.T) {
		_, err := serializer.Serialize([]any{
			InputOptionalValue{Value: 0x42},
			U8Value{Value: 0x43},
		})

		require.ErrorContains(t, err, "an optional value must be last among input values")
	})

	t.Run("multi<u8, u16, u32>", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			InputMultiValue{
				Items: []any{
					U8Value{Value: 0x42},
					U16Value{Value: 0x4243},
					U32Value{Value: 0x42434445},
				},
			},
		})

		require.NoError(t, err)
		require.Equal(t, "42@4243@42434445", data)
	})

	t.Run("u8, multi<u8, u16, u32>", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			U8Value{Value: 0x42},
			InputMultiValue{
				Items: []any{
					U8Value{Value: 0x42},
					U16Value{Value: 0x4243},
					U32Value{Value: 0x42434445},
				},
			},
		})

		require.NoError(t, err)
		require.Equal(t, "42@42@4243@42434445", data)
	})

	t.Run("multi<multi<u8, u16>, multi<u8, u16>>", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			InputMultiValue{
				Items: []any{
					InputMultiValue{
						Items: []any{
							U8Value{Value: 0x42},
							U16Value{Value: 0x4243},
						},
					},
					InputMultiValue{
						Items: []any{
							U8Value{Value: 0x44},
							U16Value{Value: 0x4445},
						},
					},
				},
			},
		})

		require.NoError(t, err)
		require.Equal(t, "42@4243@44@4445", data)
	})

	t.Run("variadic, of different types", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			InputVariadicValues{
				Items: []any{
					U8Value{Value: 0x42},
					U16Value{Value: 0x4243},
				},
			},
		})

		// For now, the serializer does not perform such a strict type check.
		// Although doable, it would be slightly complex and, if done, might be even dropped in the future
		// (with respect to the decoder that is embedded in Rust-based smart contracts).
		require.Nil(t, err)
		require.Equal(t, "42@4243", data)
	})

	t.Run("variadic<u8>, u8: should err because variadic must be last", func(t *testing.T) {
		_, err := serializer.Serialize([]any{
			InputVariadicValues{
				Items: []any{
					U8Value{Value: 0x42},
					U8Value{Value: 0x43},
				},
			},
			U8Value{Value: 0x44},
		})

		require.ErrorContains(t, err, "variadic values must be last among input values")
	})

	t.Run("u8, variadic<u8>", func(t *testing.T) {
		data, err := serializer.Serialize([]any{
			U8Value{Value: 0x41},
			InputVariadicValues{
				Items: []any{
					U8Value{Value: 0x42},
					U8Value{Value: 0x43},
				},
			},
		})

		require.Nil(t, err)
		require.Equal(t, "41@42@43", data)
	})
}

func TestSerializer_Deserialize(t *testing.T) {
	serializer, err := NewSerializer(ArgsNewSerializer{
		PartsSeparator: "@",
		PubKeyLength:   32,
	})
	require.NoError(t, err)

	t.Run("nil destination", func(t *testing.T) {
		err := serializer.Deserialize("", []any{nil})
		require.ErrorContains(t, err, "cannot deserialize into nil value")
	})

	t.Run("u8", func(t *testing.T) {
		outputValues := []any{
			&U8Value{},
		}

		err := serializer.Deserialize("42", outputValues)

		require.Nil(t, err)
		require.Equal(t, []any{
			&U8Value{Value: 0x42},
		}, outputValues)
	})

	t.Run("u16", func(t *testing.T) {
		outputValues := []any{
			&U16Value{},
		}

		err := serializer.Deserialize("4243", outputValues)

		require.Nil(t, err)
		require.Equal(t, []any{
			&U16Value{Value: 0x4243},
		}, outputValues)
	})

	t.Run("u8, u16", func(t *testing.T) {
		outputValues := []any{
			&U8Value{},
			&U16Value{},
		}

		err := serializer.Deserialize("42@4243", outputValues)

		require.Nil(t, err)
		require.Equal(t, []any{
			&U8Value{Value: 0x42},
			&U16Value{Value: 0x4243},
		}, outputValues)
	})

	t.Run("optional (missing)", func(t *testing.T) {
		outputValues := []any{
			&U8Value{},
			&OutputOptionalValue{},
		}

		err := serializer.Deserialize("42", outputValues)

		require.Nil(t, err)
		require.Equal(t, []any{
			&U8Value{Value: 0x42},
			&OutputOptionalValue{},
		}, outputValues)
	})

	t.Run("optional (provided)", func(t *testing.T) {
		outputValues := []any{
			&U8Value{},
			&OutputOptionalValue{Value: &U8Value{}},
		}

		err := serializer.Deserialize("42@43", outputValues)

		require.Nil(t, err)
		require.Equal(t, []any{
			&U8Value{Value: 0x42},
			&OutputOptionalValue{Value: &U8Value{Value: 0x43}},
		}, outputValues)
	})

	t.Run("optional: should err because optional must be last", func(t *testing.T) {
		outputValues := []any{
			&OutputOptionalValue{Value: &U8Value{}},
			&U8Value{},
		}

		err := serializer.Deserialize("43@42", outputValues)
		require.ErrorContains(t, err, "an optional value must be last among output values")
	})

	t.Run("multi<u8, u16, u32>", func(t *testing.T) {
		outputValues := []any{
			&OutputMultiValue{
				Items: []any{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		}

		err := serializer.Deserialize("42@4243@42434445", outputValues)

		require.Nil(t, err)
		require.Equal(t, []any{
			&OutputMultiValue{
				Items: []any{
					&U8Value{Value: 0x42},
					&U16Value{Value: 0x4243},
					&U32Value{Value: 0x42434445},
				},
			},
		}, outputValues)
	})

	t.Run("u8, multi<u8, u16, u32>", func(t *testing.T) {
		outputValues := []any{
			&U8Value{},
			&OutputMultiValue{
				Items: []any{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		}

		err := serializer.Deserialize("42@42@4243@42434445", outputValues)

		require.Nil(t, err)
		require.Equal(t, []any{
			&U8Value{Value: 0x42},
			&OutputMultiValue{
				Items: []any{
					&U8Value{Value: 0x42},
					&U16Value{Value: 0x4243},
					&U32Value{Value: 0x42434445},
				},
			},
		}, outputValues)
	})

	t.Run("variadic, should err because of nil item creator", func(t *testing.T) {
		destination := &OutputVariadicValues{
			Items: []any{},
		}

		err := serializer.Deserialize("", []any{destination})
		require.ErrorContains(t, err, "cannot deserialize variadic values: item creator is nil")
	})

	t.Run("empty: u8", func(t *testing.T) {
		destination := &OutputVariadicValues{
			Items:       []any{},
			ItemCreator: func() any { return &U8Value{} },
		}

		err := serializer.Deserialize("", []any{destination})
		require.NoError(t, err)
		require.Equal(t, []any{&U8Value{Value: 0}}, destination.Items)
	})

	t.Run("variadic<u8>", func(t *testing.T) {
		destination := &OutputVariadicValues{
			Items:       []any{},
			ItemCreator: func() any { return &U8Value{} },
		}

		err := serializer.Deserialize("2A@2B@2C", []any{destination})
		require.NoError(t, err)

		require.Equal(t, []any{
			&U8Value{Value: 42},
			&U8Value{Value: 43},
			&U8Value{Value: 44},
		}, destination.Items)
	})

	t.Run("varidic<u8>, with empty items", func(t *testing.T) {
		destination := &OutputVariadicValues{
			Items:       []any{},
			ItemCreator: func() any { return &U8Value{} },
		}

		err := serializer.Deserialize("@01@00@", []any{destination})
		require.NoError(t, err)

		require.Equal(t, []any{
			&U8Value{Value: 0},
			&U8Value{Value: 1},
			&U8Value{Value: 0},
			&U8Value{Value: 0},
		}, destination.Items)
	})

	t.Run("varidic<u32>", func(t *testing.T) {
		destination := &OutputVariadicValues{
			Items:       []any{},
			ItemCreator: func() any { return &U32Value{} },
		}

		err := serializer.Deserialize("AABBCCDD@DDCCBBAA", []any{destination})
		require.NoError(t, err)

		require.Equal(t, []any{
			&U32Value{Value: 0xAABBCCDD},
			&U32Value{Value: 0xDDCCBBAA},
		}, destination.Items)
	})

	t.Run("varidic<u8>, should err because decoded value is too large", func(t *testing.T) {
		destination := &OutputVariadicValues{
			Items:       []any{},
			ItemCreator: func() any { return &U8Value{} },
		}

		err := serializer.Deserialize("0100", []any{destination})
		require.ErrorContains(t, err, "cannot decode (top-level) *abi.U8Value, because of: decoded value is too large: 256 > 255")
	})
}

func TestSerializer_InRealWorldScenarios(t *testing.T) {
	serializer, err := NewSerializer(ArgsNewSerializer{
		PartsSeparator: "@",
		PubKeyLength:   32,
	})
	require.NoError(t, err)

	alicePubKeyHex := "0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1"
	alicePubKey, _ := hex.DecodeString(alicePubKeyHex)
	bobPubKeyHex := "8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8"
	bobPubKey, _ := hex.DecodeString(bobPubKeyHex)

	oneQuintillion := big.NewInt(0).SetUint64(1_000_000_000_000_000_000)

	t.Run("real-world (1): serialize input of multisig.proposeBatch(variadic<Action>), ", func(t *testing.T) {
		createEsdtTokenPayment := func(tokenIdentifier string, tokenNonce uint64, amount *big.Int) StructValue {
			return StructValue{
				Fields: []Field{
					{
						Name:  "token_identifier",
						Value: StringValue{Value: tokenIdentifier},
					},
					{
						Name:  "token_nonce",
						Value: U64Value{Value: tokenNonce},
					},
					{
						Name:  "amount",
						Value: BigUIntValue{Value: amount},
					},
				},
			}
		}

		// First action: SendTransferExecuteEgld
		firstAction := EnumValue{
			Discriminant: 5,
			// CallActionData
			Fields: []Field{
				{
					Name:  "to",
					Value: AddressValue{Value: alicePubKey},
				},
				{
					Name:  "egld_amount",
					Value: BigUIntValue{Value: oneQuintillion},
				},
				{
					Name: "opt_gas_limit",
					Value: OptionValue{
						Value: U64Value{Value: 15000000},
					},
				},
				{
					Name:  "endpoint_name",
					Value: BytesValue{Value: []byte("example")},
				},
				{
					Name: "arguments",
					Value: InputListValue{
						Items: []any{
							BytesValue{Value: []byte{0x03, 0x42}},
							BytesValue{Value: []byte{0x07, 0x43}},
						},
					},
				},
			},
		}

		// Second action: SendTransferExecuteEsdt
		secondAction := EnumValue{
			Discriminant: 6,
			// EsdtTransferExecuteData
			Fields: []Field{
				{
					Name:  "to",
					Value: AddressValue{Value: alicePubKey},
				},
				{
					Name: "tokens",
					Value: InputListValue{
						Items: []any{
							createEsdtTokenPayment("beer", 0, oneQuintillion),
							createEsdtTokenPayment("chocolate", 0, oneQuintillion),
						},
					},
				},
				{
					Name: "opt_gas_limit",
					Value: OptionValue{
						Value: U64Value{Value: 15000000},
					},
				},
				{
					Name:  "endpoint_name",
					Value: BytesValue{Value: []byte("example")},
				},
				{
					Name: "arguments",
					Value: InputListValue{
						Items: []any{
							BytesValue{Value: []byte{0x03, 0x42}},
							BytesValue{Value: []byte{0x07, 0x43}},
						},
					},
				},
			},
		}

		data, err := serializer.Serialize([]any{
			firstAction,
			secondAction,
		})

		encodedExpected := strings.Join(
			[]string{
				"05|0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1|000000080de0b6b3a7640000|010000000000e4e1c0|000000076578616d706c65|00000002000000020342000000020743",
				"06|0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1|00000002|0000000462656572|0000000000000000|000000080de0b6b3a7640000|0000000963686f636f6c617465|0000000000000000|000000080de0b6b3a7640000|010000000000e4e1c0|000000076578616d706c65|00000002000000020342000000020743",
			},
			"@",
		)

		// Drop the delimiters (were added for readability)
		encodedExpected = strings.Replace(encodedExpected, "|", "", -1)

		require.NoError(t, err)
		require.Equal(t, encodedExpected, data)
	})

	t.Run("real-world (2): deserialize output of multisig.getPendingActionFullInfo() -> variadic<ActionFullInfo>, ", func(t *testing.T) {
		dataHex := strings.Join([]string{
			"0000002A",
			"0000002A",
			"05|0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1|000000080de0b6b3a7640000|010000000000e4e1c0|000000076578616d706c65|00000002000000020342000000020743",
			"00000002|0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1|8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8",
		}, "")
		// Drop the delimiters (were added for readability)
		data := strings.Replace(dataHex, "|", "", -1)

		actionId := U32Value{}
		groupId := U32Value{}

		actionTo := AddressValue{}
		actionEgldAmount := BigUIntValue{}
		actionGasLimit := U64Value{}
		actionEndpointName := BytesValue{}
		actionArguments := OutputListValue{
			ItemCreator: func() any {
				return &BytesValue{}
			},
		}

		action := EnumValue{
			Fields: []Field{
				{
					Name:  "to",
					Value: &actionTo,
				},
				{
					Name:  "egld_amount",
					Value: &actionEgldAmount,
				},
				{
					Name: "opt_gas_limit",
					Value: &OptionValue{
						Value: &actionGasLimit,
					},
				},
				{
					Name:  "endpoint_name",
					Value: &actionEndpointName,
				},
				{
					Name:  "arguments",
					Value: &actionArguments,
				},
			},
		}

		signers := OutputListValue{
			ItemCreator: func() any {
				return &AddressValue{}
			},
		}

		destination := &OutputVariadicValues{
			ItemCreator: func() any {
				return &StructValue{
					Fields: []Field{
						{
							Name:  "action_id",
							Value: &actionId,
						},
						{
							Name:  "group_id",
							Value: &groupId,
						},
						{
							Name:  "action_data",
							Value: &action,
						},
						{
							Name:  "signers",
							Value: &signers,
						},
					},
				}
			},
		}

		err := serializer.Deserialize(data, []any{destination})
		require.NoError(t, err)
		require.Len(t, destination.Items, 1)

		// result[0].action_id and result[0].group_id
		require.Equal(t, uint32(42), actionId.Value)
		require.Equal(t, uint32(42), groupId.Value)

		// result[0].action_data
		require.Equal(t, uint8(5), action.Discriminant)
		require.Equal(t, alicePubKey, actionTo.Value)
		require.Equal(t, oneQuintillion, actionEgldAmount.Value)
		require.Equal(t, uint64(15000000), actionGasLimit.Value)
		require.Equal(t, []byte("example"), actionEndpointName.Value)
		require.Len(t, actionArguments.Items, 2)
		require.Equal(t, []byte{0x03, 0x42}, actionArguments.Items[0].(*BytesValue).Value)
		require.Equal(t, []byte{0x07, 0x43}, actionArguments.Items[1].(*BytesValue).Value)

		// result[0].signers
		require.Equal(t, alicePubKey, signers.Items[0].(*AddressValue).Value)
		require.Equal(t, bobPubKey, signers.Items[1].(*AddressValue).Value)
	})
}
