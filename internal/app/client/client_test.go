package client

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"token-validation/internal/pkg/mock_server"
)

func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	err := w.Close()
	if err != nil {
		log.Fatalf("unable to close Pipe Writer: %v", err)
	}
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func assertEqual(t *testing.T, expectedResult, actualResult string) {
	if expectedResult != actualResult {
		t.Fatalf("Results do not match;\nactual:\t %s\nexpected:\t %s", actualResult, expectedResult)
	}
}

func TestRunSuccess(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "client_id: ab\nclient_type: 1\nexpires_in: 120\nuser_id: 255\nusername: user\n\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockOkResponse{
		SvcId:       2,
		ReqId:       1,
		ClientIdLen: 2,
		ClientId:    [2]int8{'a', 'b'},
		ClientType:  1,
		UsernameLen: 4,
		Username:    [4]int8{'u', 's', 'e', 'r'},
		ExpiresIn:   120,
		UserId:      255,
	}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, expectedReqBody.Encode(), resp.Encode())
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() })
	assertEqual(t, expectedResult, actualResult)
}
