package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	START_BLOCK byte = 11
	END_BLOCK   byte = 28
	CR          byte = 13
)

type ReqData struct {
	Message string
}

type ResData struct {
	Status string
	Data   string
}

func main() {
	fmt.Println("I'm starting to do something now!")
	listen, err := net.Listen("tcp", "localhost:2575")
	if err != nil {
		log.Fatal(err)
	}

	// Make sure to close the conneciton after leaving main
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// CReate a new goroutine for handling the connection
		go handleRequest(conn)

	}
}

func handleRequest(conn net.Conn) {
	// Make sure connection is closed after we are done
	defer conn.Close()
	timeout := 5 * time.Second

	var msg []byte
	var status string = "start"
	tmpBufSize := 256
	totalReceived := 0
	var tmpBuf []byte
	for {
		conn.SetDeadline(time.Now().Add(timeout))
		tmpBuf = make([]byte, tmpBufSize)
		receivedLen, err := conn.Read(tmpBuf)
		totalReceived = totalReceived + receivedLen
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < receivedLen; i++ {
			b := tmpBuf[i]
			switch status {
			case "start":
				if b == START_BLOCK {
					status = "msg"
				} else {
					return
				}
			case "msg":
				if b == END_BLOCK {
					status = "end"
				} else if b == START_BLOCK {
					return
				} else {
					msg = append(msg, b)
				}
			case "end":
				if b == CR {
					handleMessage(conn, msg)
					totalReceived = 0
					status = "start"
					msg = nil
				} else {
					return
				}
			}
		}

		if receivedLen < tmpBufSize {
			fmt.Println("No more messages there ought to be")
			break
		} else if receivedLen == tmpBufSize {
			// TODO
			break
		}

	}

}

func handleMessage(conn net.Conn, msg []byte) {
	fmt.Println(string(msg))
	reqData := ReqData{string(msg)}

	s, _ := json.Marshal(reqData)
	buf := bytes.NewBuffer(s)
	resp, err := http.Post("https://ptsv2.com/t/tz0jp-1616963104/post",
		"application/json", buf)

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var response ResData
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	conn.Write([]byte(response.Data))

}
