package mock_server

import (
	"bytes"
	"encoding/binary"
	"log"
)

type MockRequestBody struct {
	// Body
	SvcMsg   int32
	TokenLen int32
	// Array required to Encode whole struct at once
	Token    [8]int8
	ScopeLen int32
	Scope    [3]int8
}

func (req *MockRequestBody) Encode() []byte {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, req)
	if err != nil {
		log.Fatalf("Unable to Encode request: %v", err)
	}
	return buf.Bytes()
}

type MockOkResponse struct {
	// Header
	SvcId   int32
	BodyLen int32
	ReqId   int32
	// Body
	RetCode     int32
	ClientIdLen int32
	ClientId    [2]int8
	ClientType  int32
	UsernameLen int32
	Username    [4]int8
	ExpiresIn   int32
	UserId      int64
}

func (resp *MockOkResponse) Encode() []byte {
	// 28 = 5 * sizeOf(int32) + sizeOf(int64)
	resp.BodyLen = 28 + resp.ClientIdLen + resp.UsernameLen
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, resp)
	if err != nil {
		log.Fatalf("Unable to Encode response: %v", err)
	}
	return buf.Bytes()
}

type MockErrorResponse struct {
	// Header
	SvcId   int32
	BodyLen int32
	ReqId   int32
	// Body
	RetCode        int32
	ErrorStringLen int32
	ErrorString    []int8
}

func (resp *MockErrorResponse) Encode() []byte {
	// 8 = 2 * sizeOf(int32)
	resp.BodyLen = 8 + resp.ErrorStringLen
	buf := &bytes.Buffer{}
	// encoding fields with fixed size
	err := binary.Write(buf, binary.LittleEndian, struct {
		SvcId          int32
		BodyLen        int32
		ReqId          int32
		RetCode        int32
		ErrorStringLen int32
	}{resp.SvcId, resp.BodyLen, resp.ReqId, resp.RetCode, resp.ErrorStringLen})
	if err != nil {
		log.Fatalf("Unable to Encode response: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, resp.ErrorString)
	if err != nil {
		log.Fatalf("Unable to Encode response error string: %v", err)
	}
	return buf.Bytes()
}
