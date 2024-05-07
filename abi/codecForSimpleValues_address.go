package abi

import (
	"fmt"
	"io"
)

type codecForAddress struct {
	pubKeyLength int
}

func (c *codecForAddress) encodeNested(writer io.Writer, value AddressValue) error {
	err := c.checkPubKeyLength(value.Value)
	if err != nil {
		return err
	}

	_, err = writer.Write(value.Value)
	return err
}

func (c *codecForAddress) encodeTopLevel(writer io.Writer, value AddressValue) error {
	return c.encodeNested(writer, value)
}

func (c *codecForAddress) decodeNested(reader io.Reader, value *AddressValue) error {
	data, err := readBytesExactly(reader, c.pubKeyLength)
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

func (c *codecForAddress) decodeTopLevel(data []byte, value *AddressValue) error {
	err := c.checkPubKeyLength(data)
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}

func (c *codecForAddress) checkPubKeyLength(pubkey []byte) error {
	if len(pubkey) != c.pubKeyLength {
		return fmt.Errorf("public key (address) has invalid length: %d", len(pubkey))
	}

	return nil
}
