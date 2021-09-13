package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	oauth2req "token-validation/internal/pkg/request/oauth2"
	oauth2resp "token-validation/internal/pkg/response/oauth2_reader"
)

var (
	ByteOrder = binary.LittleEndian
)

func Run() {
	cliArgs := os.Args[1:]
	if len(cliArgs) != 4 {
		log.Fatalf("invalid cli args; provide \"host port token scope\"")
	}
	host, port, token, scope := cliArgs[0], cliArgs[1], cliArgs[2], cliArgs[3]

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("unable to dial the server: %v", err)
	}
	defer conn.Close()
	req := oauth2req.Create(token, scope)
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
