package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"token-validation/internal/pkg/id_generator"
	"token-validation/internal/pkg/mock_server"
)

func captureOutput(f func(), t *testing.T) string {
	rescueStdout := os.Stdout
	origLogFatalf := logFatalf
	defer func() {
		logFatalf = origLogFatalf
	}()
	r, w, _ := os.Pipe()
	os.Stdout = w
	logFatalf = func(format string, args ...interface{}) {
		// Make Fatal less fatal for testing purposes :)
		var err error
		if len(args) > 0 {
			_, err = w.WriteString(fmt.Sprintf(format, args))
		} else {
			_, err = w.WriteString(format)
		}
		if err != nil {
			t.Fatalf("unable to write string: %v", err)
		}
	}
	f()
	err := w.Close()
	if err != nil {
		t.Fatalf("unable to close Pipe Writer: %v", err)
	}
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func assertEqual(t *testing.T, expectedResult, actualResult string) {
	if expectedResult != actualResult {
		t.Fatalf("Results do not match;\nactual:  %s\nexpected:\t %s", actualResult, expectedResult)
	}
}

func TestRunSuccess(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "client_id: ab\nclient_type: 1\nexpires_in: 120\nuser_id: 255\nusername: user\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockOkResponse{
		SvcId:       2,
		ReqId:       id_generator.GetInstance().GetLastId() + 1,
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
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunTokenNotFound(t *testing.T) {
	token, scope := "xxxxxxxx", "abc"
	expectedResult := "error: CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND\nmessage: token not found\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'x', 'x', 'x', 'x', 'x', 'x', 'x', 'x'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockErrorResponse{
		SvcId:          2,
		ReqId:          id_generator.GetInstance().GetLastId() + 1,
		RetCode:        0x00000001,
		ErrorStringLen: 15,
		ErrorString: []int8{
			't', 'o', 'k', 'e', 'n', ' ',
			'n', 'o', 't', ' ',
			'f', 'o', 'u', 'n', 'd'},
	}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, expectedReqBody.Encode(), resp.Encode())
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunDbError(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "error: CUBE_OAUTH2_ERR_DB_ERROR\nmessage: db error\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockErrorResponse{
		SvcId:          2,
		ReqId:          id_generator.GetInstance().GetLastId() + 1,
		RetCode:        0x00000002,
		ErrorStringLen: 8,
		ErrorString:    []int8{'d', 'b', ' ', 'e', 'r', 'r', 'o', 'r'},
	}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, expectedReqBody.Encode(), resp.Encode())
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunUnknownMsg(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "error: CUBE_OAUTH2_ERR_UNKNOWN_MSG\nmessage: unknown svc message\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockErrorResponse{
		SvcId:          2,
		ReqId:          id_generator.GetInstance().GetLastId() + 1,
		RetCode:        0x00000003,
		ErrorStringLen: 19,
		ErrorString: []int8{
			'u', 'n', 'k', 'n', 'o', 'w', 'n', ' ',
			's', 'v', 'c', ' ', 'm', 'e', 's', 's', 'a', 'g', 'e',
		},
	}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, expectedReqBody.Encode(), resp.Encode())
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunBadPacket(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "error: CUBE_OAUTH2_ERR_BAD_PACKET\nmessage: bad packet\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockErrorResponse{
		SvcId:          2,
		ReqId:          id_generator.GetInstance().GetLastId() + 1,
		RetCode:        0x00000004,
		ErrorStringLen: 10,
		ErrorString: []int8{
			'b', 'a', 'd', ' ', 'p', 'a', 'c', 'k', 'e', 't',
		},
	}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, expectedReqBody.Encode(), resp.Encode())
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunBadClient(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "error: CUBE_OAUTH2_ERR_BAD_CLIENT\nmessage: bad client\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockErrorResponse{
		SvcId:          2,
		ReqId:          id_generator.GetInstance().GetLastId() + 1,
		RetCode:        0x00000005,
		ErrorStringLen: 10,
		ErrorString: []int8{
			'b', 'a', 'd', ' ', 'c', 'l', 'i', 'e', 'n', 't',
		},
	}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, expectedReqBody.Encode(), resp.Encode())
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunBadScope(t *testing.T) {
	token, scope := "abracada", "xxx"
	expectedResult := "error: CUBE_OAUTH2_ERR_BAD_SCOPE\nmessage: bad scope\n"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'x', 'x', 'x'},
	}
	resp := &mock_server.MockErrorResponse{
		SvcId:          2,
		ReqId:          id_generator.GetInstance().GetLastId() + 1,
		RetCode:        0x00000006,
		ErrorStringLen: 9,
		ErrorString: []int8{
			'b', 'a', 'd', ' ', 's', 'c', 'o', 'p', 'e',
		},
	}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, expectedReqBody.Encode(), resp.Encode())
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunNoServer(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "unable to dial the launchMockServer: [dial tcp [::1]:2000: connect: connection refused]"
	os.Args = append(os.Args[:1], "localhost", "2000", token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunInvalidResponse(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "unable to read response: [unexpected EOF]"
	resp := []byte{0x00, 0x00, 0x00}
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, nil, resp)
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunNoResponse(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "unable to read response: [EOF]"
	servReady := make(chan struct{})
	go mock_server.LaunchMockServer(servReady, nil, nil)
	<-servReady
	os.Args = append(os.Args[:1], "localhost", mock_server.PORT, token, scope)
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}

func TestRunInvalidRequestId(t *testing.T) {
	token, scope := "abracada", "abc"
	expectedResult := "request id does not match response id"
	expectedReqBody := &mock_server.MockRequestBody{
		SvcMsg:   0x00000001,
		TokenLen: 8,
		Token:    [8]int8{'a', 'b', 'r', 'a', 'c', 'a', 'd', 'a'},
		ScopeLen: 3,
		Scope:    [3]int8{'a', 'b', 'c'},
	}
	resp := &mock_server.MockOkResponse{
		SvcId:       2,
		ReqId:       0,
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
	actualResult := captureOutput(func() { Run() }, t)
	assertEqual(t, expectedResult, actualResult)
}
