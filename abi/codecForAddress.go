package abi

import (
	"fmt"
	"io"
)

// AddressValue is a wrapper for an address
type AddressValue struct {
	Value []byte
}

func (value *AddressValue) encodeNested(writer io.Writer) error {
	err := value.checkPubKeyLength(value.Value)
	if err != nil {
		return err
	}

	_, err = writer.Write(value.Value)
	return err
}

func (value *AddressValue) encodeTopLevel(writer io.Writer) error {
	return value.encodeNested(writer)
}

func (value *AddressValue) decodeNested(reader io.Reader) error {
	data, err := readBytesExactly(reader, pubKeyLength)
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

func (value *AddressValue) decodeTopLevel(data []byte) error {
	err := value.checkPubKeyLength(data)
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

func (value *AddressValue) checkPubKeyLength(pubkey []byte) error {
	if len(pubkey) != pubKeyLength {
		return fmt.Errorf("public key (address) has invalid length: %d", len(pubkey))
	}

	return nil
}
