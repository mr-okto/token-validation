package request

import (
	"encoding/binary"
	"io"
)

type Request interface {
	Write(w io.Writer, order binary.ByteOrder) error
	GetId() int32
}
