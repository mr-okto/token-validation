package mock_server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

const (
	PORT = "2000"
)

func handleConn(ln net.Listener, expectedReqBody []byte, responseMock []byte) {
	conn, err := ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	headerData := struct {
		SvcId      int32
		BodyLength int32
		RequestId  int32
	}{}
	err = binary.Read(conn, binary.LittleEndian, &headerData)
	if err != nil {
		log.Fatalf("Unable to read header from Conn: %v", err)
	}
	actualReqBody := make([]byte, headerData.BodyLength)
	_, err = conn.Read(actualReqBody)
	if err != nil {
		log.Fatalf("Unable to read body from Conn: %v", err)
	}
	if expectedReqBody != nil && !bytes.Equal(actualReqBody, expectedReqBody) {
		log.Fatalf("Invalid request body;\nactual:\t%v\nexpected:\t%v", actualReqBody, expectedReqBody)
	}
	if responseMock == nil {
		return // Server should not send anything back
	}
	_, err = conn.Write(responseMock)
	if err != nil {
		log.Fatalf("Unable to write response: %v", err)
	}
	return
}

// LaunchMockServer expectedReqBody should be nil for it not to be compared with actual request
// LaunchMockServer resp should be nil for server not to send anything back
func LaunchMockServer(ready chan struct{}, expectedReqBody []byte, resp []byte) {
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("Unable to listen: %v", err)
	}
	defer ln.Close()
	close(ready)
	handleConn(ln, expectedReqBody, resp)
}
