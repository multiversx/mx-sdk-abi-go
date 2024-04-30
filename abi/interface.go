package abi

import "io"

// generalCodec is an internal interface that allows "leaf" codecs to rely on the general "composite" codec, if needed.
type generalCodec interface {
	doEncodeNested(writer io.Writer, value any) error
	doDecodeNested(reader io.Reader, value any) error
}
