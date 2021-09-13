package oauth2_ok

import (
	"encoding/binary"
	"fmt"
	"io"
	hdr "token-validation/internal/pkg/header"
	"token-validation/internal/pkg/response"
	"token-validation/internal/pkg/str"
)

type Body struct {
	ReturnCode int32
	ClientId   *str.Str
	ClientType int32
	Username   *str.Str
	ExpiresIn  int32
	UserId     int64
}

type OkResponse struct {
	header *hdr.Header
	body   *Body
}

func Create(h *hdr.Header) response.Response {
	return &OkResponse{
		header: h,
		body:   &Body{},
	}
}

func (r *OkResponse) GetId() int32 {
	return r.header.RequestId
}

// ReadBody Reads all fields except for ReturnCode
func (r *OkResponse) ReadBody(reader io.Reader, order binary.ByteOrder) (err error) {
	body := Body{}
	body.ClientId, err = str.FromReader(reader, order)
	if err != nil {
		return err
	}
	err = binary.Read(reader, order, &body.ClientType)
	if err != nil {
		return err
	}
	body.Username, err = str.FromReader(reader, order)
	if err != nil {
		return err
	}
	err = binary.Read(reader, order, &body.ExpiresIn)
	if err != nil {
		return err
	}
	err = binary.Read(reader, order, &body.UserId)
	if err != nil {
		return err
	}
	// 28 = sizeOf(ReturnCode) + sizeOf(ClientId.Len) + sizeOf(ClientType) + sizeOf(Username.Len) +
	//     + sizeOf(ExpiresIn) + sizeOf(UserId)
	bodyLength := 28 + body.ClientId.Len + body.Username.Len
	if r.header.BodyLength != bodyLength {
		return fmt.Errorf("invalid response body length: expected: %d;actual: %d",
			r.header.BodyLength, bodyLength)
	}
	r.body = &body
	return nil
}

func (r *OkResponse) Print(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, "client_id: %s\n"+
		"client_type: %d\n"+
		"expires_in: %d\n"+
		"user_id: %d\n"+
		"username: %s",
		r.body.ClientId.ToString(),
		r.body.ClientType,
		r.body.ExpiresIn,
		r.body.UserId,
		r.body.Username.ToString())
	return err
}
