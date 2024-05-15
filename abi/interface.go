package abi

import "io"

type singleValue interface {
	encodeNested(writer io.Writer) error
	encodeTopLevel(writer io.Writer) error
	decodeNested(reader io.Reader) error
	decodeTopLevel(data []byte) error
}
