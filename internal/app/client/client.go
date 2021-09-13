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
	logFatalf = log.Fatalf // might be reassigned during testing
)

func Run() {
	cliArgs := os.Args[1:]
	if len(cliArgs) != 4 {
		logFatalf("invalid cli args; provide \"host port Token Scope\"")
		return
	}
	host, port, token, scope := cliArgs[0], cliArgs[1], cliArgs[2], cliArgs[3]
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		logFatalf("unable to dial the launchMockServer: %v", err)
		return
	}
	defer conn.Close()
	req := oauth2req.Create(token, scope)
	err = req.Write(conn, ByteOrder)
	if err != nil {
		logFatalf("unable to send request: %v", err)
		return
	}
	resp, err := oauth2resp.Read(conn, ByteOrder)
	if err != nil {
		logFatalf("unable to read response: %v", err)
		return
	}
	if req.GetId() != resp.GetId() {
		logFatalf("request id [%d] does not match response id [%d]",
			req.GetId(), resp.GetId())
		return
	}
	var b strings.Builder
	err = resp.Print(&b)
	if err != nil {
		logFatalf("unable to output request: %v", err)
		return
	}
	fmt.Println(b.String())
}
