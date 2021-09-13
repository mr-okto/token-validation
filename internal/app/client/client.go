package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	oauth2req "token-validation/internal/pkg/request/oauth2"
	oauth2resp "token-validation/internal/pkg/response/oauth2_reader"
)

var (
	ByteOrder = binary.LittleEndian
)

func Run() {
	req := oauth2req.Create("ab", "x")

	conn, err := net.Dial("tcp", "localhost:2000")
	if err != nil {
		log.Fatalf("unable to dial the server: %v", err)
	}
	defer conn.Close()

	err = req.Write(conn, ByteOrder)
	if err != nil {
		log.Fatalf("unable to send request: %v", err)
	}

	resp, err := oauth2resp.Read(conn, ByteOrder)
	if err != nil {
		log.Fatalf("unable to read response: %v", err)
	}
	if req.GetId() != resp.GetId() {
		log.Fatalf("request id [%d] does not match response id [%d]",
			req.GetId(), resp.GetId())
	}
	var b strings.Builder
	err = resp.Print(&b)
	if err != nil {
		log.Fatalf("unable to output request: %v", err)
	}
	fmt.Println(b.String())
}
