package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	START_BLOCK byte = 11
	END_BLOCK   byte = 28
	CR          byte = 13
)

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
					fmt.Println("Start start")
					status = "msg"
				} else {
					fmt.Println("No start found")
					return
				}
			case "msg":
				if b == END_BLOCK {
					fmt.Println("Message ends here")
					status = "end"
				} else if b == START_BLOCK {
					fmt.Println("Start block should be in here")
					return
				} else {
					fmt.Println("Normal " + string(b))
					msg = append(msg, b)
				}
			case "end":
				if b == CR {
					fmt.Println("Message really ends here. Processing it")
					handleMessage(msg)
					totalReceived = 0
					status = "start"
					msg = nil
				} else {
					fmt.Println("Expected CR after end")
					return
				}
			}
		}

		if receivedLen < tmpBufSize {
			fmt.Println("No more messages there ought to be")
			break
		} else if receivedLen == tmpBufSize {
			fmt.Println("COuld be, could be not")
			break
		}

	}

}

func handleMessage(msg []byte) {
	fmt.Println(string(msg))
}
