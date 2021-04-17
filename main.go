package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
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

func main() {

	help := flag.Bool("help", false, "Show help")
	ip := flag.String("ip", "localhost", "Address to be listened, e.g. localhost or 0.0.0.0.")
	port := flag.Int("port", 2575, "Port to be listened. ")
	targetUrl := flag.String("url", "", "Target URL where data is sent as HTTP POST")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}
	if *targetUrl == "" {
		log.Fatal("Target URL must be given")
	}
	address := fmt.Sprintf("%s:%d", *ip, *port)

	startServer(address, *targetUrl)
}

func startServer(address string, targetUrl string) {
	fmt.Printf("Starting to listen %s and routing messages to %s\n", address, targetUrl)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure to close the conneciton after leaving
	defer listen.Close()

	for {
		acceptAndHandleRequest(listen, targetUrl)
	}
}

func acceptAndHandleRequest(listen net.Listener, targetUrl string) {
	conn, err := listen.Accept()
	if err != nil {
		log.Fatal(err)
	}
	// Create a new goroutine for handling the connection
	go handleRequest(conn, targetUrl)
}

func handleRequest(conn net.Conn, targetUrl string) {
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
		if err := handleMessage(conn, msg[1:len(msg)-1], targetUrl); err != nil {
			conn.Write([]byte(err.Error()))
		}

	}

}

func handleMessage(conn net.Conn, msg []byte, targetUrl string) error {
	fmt.Print("INFO: Processing message\n", strings.ReplaceAll(string(msg), "\r", "\n"), "\n")
	reqData := ReqData{string(msg)}

	s, _ := json.Marshal(reqData)
	buf := bytes.NewBuffer(s)
	resp, err := http.Post(targetUrl,
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
	var response ReqData
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("ERROR: Unmarshaling POST response body: ", err)
		return errors.New("Error")
	}

	conn.Write(createMLLPMessage(response.Message))
	fmt.Println("INFO: Done")
	return nil
}

func createMLLPMessage(msg string) []byte {
	b := []byte{START_BLOCK}
	b = append(b, []byte(msg)...)
	b = append(b, END_BLOCK, CR)
	return b
}
