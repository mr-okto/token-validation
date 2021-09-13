package oauth2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	hdr "token-validation/internal/pkg/header"
	"token-validation/internal/pkg/request"
	"token-validation/internal/pkg/str"
)

type Body struct {
	SvcMsg int32
	Token  *str.Str
	Scope  *str.Str
}

type Request struct {
	header *hdr.Header
	body   *Body
}

const (
	SvcMsg = int32(0x00000001)
	SvcId  = int32(0x00000002)
)

func Create(token string, scope string) request.Request {
	body := &Body{
		SvcMsg: SvcMsg,
		Token:  str.FromString(token),
		Scope:  str.FromString(scope),
	}
	// 12 = sizeOf(SvcMsg) + sizeOf(Token.Len) + sizeOf(Scope.Len)
	bodyLen := 12 + body.Token.Len + body.Scope.Len
	h := &hdr.Header{
		SvcId:      SvcId,
		BodyLength: bodyLen,
		RequestId:  0x00000005, //TODO: generate IDs
	}
	return &Request{
		header: h,
		body:   body,
	}
}

func (r *Request) GetId() int32 {
	return r.header.RequestId
}

func (r *Request) Write(w io.Writer, order binary.ByteOrder) error {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, order, r.header)
	if err != nil {
		return err
	}
	err = binary.Write(buf, order, r.body.SvcMsg)
	if err != nil {
		return err
	}
	err = r.body.Token.Write(buf, order)
	if err != nil {
		return err
	}
	err = r.body.Scope.Write(buf, order)
	if err != nil {
		return err
	}

	written, err := buf.WriteTo(w)
	if err != nil {
		return err
	}
	payloadSize := int64(hdr.Size) + int64(r.header.BodyLength)
	if written != payloadSize {
		return fmt.Errorf("unable to write request; "+
			"written bytes: %d; request length: %d", written, payloadSize)
	}
	return nil
}
