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

type Value struct {
	val    string
	expiry int64
}

var (
	PORT     int                              = 6379
	HashMap  map[string]Value                 = make(map[string]Value)
	Commands map[string]func([]string) []byte = make(map[string]func([]string) []byte)
)

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
	fmt.Print(line, " ")
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
			fmt.Print(sizeBuffer, " ")

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
	return nil, errors.New("InvalidArgumentError")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		print("start -> ")
		messages, err := decodeMsg(reader)
		if err != nil {
			print(err)
			conn.Write(encodeSimpleError(err.Error()))
			break
		}

		function, ok := Commands[strings.ToLower(messages[0])]

		if !ok {
			conn.Write(encodeSimpleError("Unknown Command"))
		}

		response := function(messages[1:])
		conn.Write(response)
		fmt.Println(messages, " <- end")
	}
}

func main() {
	Commands["ping"] = ping
	Commands["echo"] = echo
	Commands["command"] = command
	Commands["set"] = set
	Commands["get"] = get
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(PORT))
	if err != nil {
		fmt.Println("Failed to bind to port")
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
