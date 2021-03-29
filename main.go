package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
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

	reader := bufio.NewReader(conn)
	for {
		conn.SetDeadline(time.Now().Add(timeout))
		msg, err := reader.ReadBytes(END_BLOCK)
		if err != nil {
			if err == io.EOF {
				fmt.Println("INFO: EOF reached, closing connection")
				return
			}
			fmt.Print("ERROR: Reading connection: ", err)
			return
		}
		if msg[0] != START_BLOCK {
			fmt.Println("ERROR: Message first character should be byte 11")
			return
		}

		if b, err := reader.ReadByte(); b != CR || err != nil {
			fmt.Println("ERROR: End block should be followed by CR")
			return
		}
		if err := handleMessage(conn, msg[1:len(msg)-1]); err != nil {
			conn.Write([]byte(err.Error()))
		}

	}

}

func handleMessage(conn net.Conn, msg []byte) error {
	fmt.Print("INFO: Processing message\n", strings.ReplaceAll(string(msg), "\r", "\n"), "\n")
	reqData := ReqData{string(msg)}

	s, _ := json.Marshal(reqData)
	buf := bytes.NewBuffer(s)
	resp, err := http.Post("https://ptsv2.com/t/tz0jp-1616963104/post",
		"application/json", buf)

	if err != nil {
		fmt.Println("ERROR: Post request error: ", err)
		return errors.New("Error")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Reading POST response body: ", err)
		return errors.New("Error")
	}
	var response ResData
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("ERROR: Unmarshaling POST response body: ", err)
		return errors.New("Error")
	}

	conn.Write([]byte(response.Data))
	return nil
}
