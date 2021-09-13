package client

import (
	"encoding/binary"
	"fmt"
	"net"
	"token-validation/internal/pkg/request/oauth2"
	"token-validation/internal/pkg/str"
)

var (
	ByteOrder = binary.LittleEndian
)

// Run TODO: remove panic
func Run() {
	req := oauth2.Create("ab", "x")
	fmt.Printf("Sending the request: %v\n", req)

	conn, err := net.Dial("tcp", "localhost:2000")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	err = req.Write(conn, ByteOrder)
	if err != nil {
		panic(err)
	}

	readStruct, err := str.FromReader(conn, ByteOrder)
	if err != nil {
		panic(err)
	}
	fmt.Println("ReadStruct: ", readStruct)
}
