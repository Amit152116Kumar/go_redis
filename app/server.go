package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var PORT int = 6379

func decodeMsg(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString(byte('\n'))
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection close by client.")
			return nil, err
		}
		fmt.Println("read error", err)
		return nil, err
	}

	line = strings.TrimSpace(line)
	fmt.Println(line[1:])
	switch line[0] {
	case '*':
		size, err := strconv.Atoi(line[1:])
		if err != nil {
			fmt.Println(err)
			break
		}

		parts := make([]string, size)
		for i := 0; i < size; i++ {
			line, err = reader.ReadString(byte('\n'))
			if err != nil {
				fmt.Println("Connection close by client or read error", err)
				break
			}
			sizeBuffer := strings.TrimSpace(line)
			if sizeBuffer[0] != '$' {
				return nil, errors.New("InvalidArgumentError")
			}
			length, _ := strconv.Atoi(sizeBuffer[1:])
			fmt.Println(length)

			data := make([]byte, length)
			_, err = io.ReadFull(reader, data)
			if err != nil {
				fmt.Println("Connection close by client or read error", err)
				break
			}
			_, err = reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			parts[i] = string(data)

		}
		return parts, nil
	default:
		fmt.Printf("Unsupported RESP type: %c \n", line[0])
		return nil, errors.New("InvalidArgumentError")
	}
	return nil, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		println("start")
		messages, err := decodeMsg(reader)
		// message, err := reader.ReadString('\n')
		if err != nil {
			print(err)
			break
		}
		for _, msg := range messages {
			if strings.EqualFold(msg, "PING") {
				conn.Write([]byte("+PONG\r\n"))
			} else if strings.EqualFold(msg, "COMMAND") {
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte("-ERR unknown command\r\n"))
			}
			fmt.Println(msg, "end")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(PORT))
	if err != nil {
		fmt.Println("Failed to bind to port", PORT)
		fmt.Print(err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		println("Listening")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
