package oauth2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"token-validation/internal/pkg/header"
	"token-validation/internal/pkg/str"
)

type Body struct {
	SvcMsg int32
	Token  *str.Str
	Scope  *str.Str
}

type Request struct {
	Header *header.Header
	Body   *Body
}

const (
	SvcMsg = int32(0x00000001)
	SvcId  = int32(0x00000002)
)

func Create(token string, scope string) *Request {
	body := &Body{
		SvcMsg: SvcMsg,
		Token:  str.FromString(token),
		Scope:  str.FromString(scope),
	}
	// 12 = sizeOf(SvcMsg) + sizeOf(Token.Len) + sizeOf(Scope.Len)
	bodyLen := 12 + body.Token.Len + body.Scope.Len
	h := &header.Header{
		SvcId:      SvcId,
		BodyLength: bodyLen,
		RequestId:  0x00000005, //TODO: generate IDs
	}
	return &Request{
		Header: h,
		Body:   body,
	}
}

func (r *Request) Write(w io.Writer, order binary.ByteOrder) error {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, order, r.Header)
	if err != nil {
		return err
	}
	err = binary.Write(buf, order, r.Body.SvcMsg)
	if err != nil {
		return err
	}
	err = r.Body.Token.Write(buf, order)
	if err != nil {
		return err
	}
	err = r.Body.Scope.Write(buf, order)
	if err != nil {
		return err
	}

	written, err := buf.WriteTo(w)
	if err != nil {
		return err
	}
	payloadSize := int64(header.Size) + int64(r.Header.BodyLength)
	if written != payloadSize {
		return fmt.Errorf("unable to write request; "+
			"written bytes: %d; request length: %d", written, payloadSize)
	}
	return nil
}
