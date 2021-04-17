package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMLLPMessageFormat(t *testing.T) {
	actual_msg := "sometestmesage"
	got := createMLLPMessage(actual_msg)
	start := got[0]
	end := got[len(got)-2:]
	msg := got[1 : len(got)-2]
	if start != 11 {
		t.Errorf("Incorrect starting block, got %d", int(got[0]))
	}
	if end[0] != 28 || end[1] != 13 {
		t.Errorf("Incorrect ending block, got %d, %d", end[0], end[1])
	}
	if string(msg) != actual_msg {
		t.Errorf("Incorrect message got %s", msg)
	}
}

func TestServer(t *testing.T) {
	mllp_server_address := "localhost:2575"
	mllp_message := createMLLPMessage("AlmostRandomTestStringItAl_most{{||")
	expected_msg := "MSH|Something|something"
	expected_response := createMLLPMessage(expected_msg)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"message\": \"%s\"}", expected_msg)
	}))
	defer ts.Close()

	go startServer(mllp_server_address, ts.URL)

	sendMsgAndValidate(t, mllp_server_address, mllp_message, expected_response)
}

func sendMsgAndValidate(t *testing.T, mllp_server_address string, mllp_message []byte, expected_response []byte) {
	conn, err := net.Dial("tcp", mllp_server_address)
	if err != nil {
		t.Errorf("Could not connect to the TCP server: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write(mllp_message); err != nil {
		t.Errorf("could not write payload to TCP server: %v", err)
	}

	if out, err := ioutil.ReadAll(conn); err == nil {
		if !bytes.Equal(out, expected_response) {
			t.Errorf("Response did not match expected output, gotmsg %s", out)
		}
	} else {
		t.Errorf("could not read from connection: %v", err)
	}
}
