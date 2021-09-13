package str

import (
	"encoding/binary"
	"io"
)

type Str struct {
	Len  int32
	Data []int8
}

func FromString(s string) *Str {
	bs := []byte(s)
	res := &Str{
		Len:  int32(len(bs)),
		Data: make([]int8, len(bs)),
	}
	for i, v := range bs {
		res.Data[i] = int8(v)
	}
	return res
}

func FromReader(r io.Reader, order binary.ByteOrder) (*Str, error) {
	var bodyLen int32
	err := binary.Read(r, order, &bodyLen)
	if err != nil || bodyLen <= 0 {
		return nil, err
	}
	bodyData := make([]int8, bodyLen)
	err = binary.Read(r, order, bodyData)
	if err != nil {
		return nil, err
	}
	return &Str{
		Len:  bodyLen,
		Data: bodyData,
	}, nil
}

func (s *Str) Write(w io.Writer, order binary.ByteOrder) error {
	err := binary.Write(w, order, s.Len)
	if err != nil {
		return err
	}
	err = binary.Write(w, order, s.Data)
	return err
}

func (s *Str) ToString() string {
	bs := make([]byte, len(s.Data))
	for i, v := range s.Data {
		bs[i] = byte(v)
	}
	return string(bs)
}
