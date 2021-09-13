package oauth2_error

import (
	"encoding/binary"
	"fmt"
	"io"
	hdr "token-validation/internal/pkg/header"
	"token-validation/internal/pkg/response"
	"token-validation/internal/pkg/str"
)

var (
	returnCodes = map[int32]string{
		0x00000001: "CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND",
		0x00000002: "CUBE_OAUTH2_ERR_DB_ERROR",
		0x00000003: "CUBE_OAUTH2_ERR_UNKNOWN_MSG",
		0x00000004: "CUBE_OAUTH2_ERR_BAD_PACKET",
		0x00000005: "CUBE_OAUTH2_ERR_BAD_CLIENT",
		0x00000006: "CUBE_OAUTH2_ERR_BAD_SCOPE",
	}
)

type Body struct {
	ReturnCode  int32
	ErrorString *str.Str
}

type ErrorResponse struct {
	header *hdr.Header
	body   *Body
}

func Create(h *hdr.Header, returnCode int32) response.Response {
	return &ErrorResponse{
		header: h,
		body: &Body{
			ReturnCode: returnCode,
		},
	}
}

func (r *ErrorResponse) GetId() int32 {
	return r.header.RequestId
}

// ReadBody Reads error string
func (r *ErrorResponse) ReadBody(reader io.Reader, order binary.ByteOrder) (err error) {
	body := Body{
		ReturnCode: r.body.ReturnCode,
	}
	body.ErrorString, err = str.FromReader(reader, order)
	if err != nil {
		return err
	}
	// 8 = sizeOf(ReturnCode) + sizeOf(ErrorString.Len)
	bodyLength := 8 + body.ErrorString.Len
	if r.header.BodyLength != bodyLength {
		return fmt.Errorf("invalid response body length: expected: %d;actual: %d",
			r.header.BodyLength, bodyLength)
	}
	r.body = &body
	return nil
}

func (r *ErrorResponse) Print(w io.Writer) (err error) {
	errorCode, ok := returnCodes[r.body.ReturnCode]
	if !ok {
		return fmt.Errorf("invalid return code: %d", r.body.ReturnCode)
	}
	_, err = fmt.Fprintf(w, "error: %s\nmessage: %s",
		errorCode, r.body.ErrorString.ToString())
	return err
}
