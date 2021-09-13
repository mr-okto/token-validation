package oauth2_reader

import (
	"encoding/binary"
	"io"
	"token-validation/internal/pkg/header"
	"token-validation/internal/pkg/response"
	"token-validation/internal/pkg/response/oauth2_error"
	"token-validation/internal/pkg/response/oauth2_ok"
)

func Read(r io.Reader, order binary.ByteOrder) (response.Response, error) {
	h := &header.Header{}
	err := binary.Read(r, order, h)
	if err != nil {
		return nil, err
	}
	var returnCode int32
	err = binary.Read(r, order, &returnCode)
	if err != nil {
		return nil, err
	}
	var resp response.Response
	if returnCode == 0 {
		resp = oauth2_ok.Create(h)
	} else {
		resp = oauth2_error.Create(h, returnCode)
	}
	err = resp.ReadBody(r, order)
	return resp, err
}
