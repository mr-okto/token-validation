package response

import (
	"encoding/binary"
	"io"
)

type Response interface {
	Print(w io.Writer) error
	GetId() int32
	ReadBody(reader io.Reader, order binary.ByteOrder) error
}
