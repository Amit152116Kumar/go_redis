package main

import "fmt"

func wrongArguments(cmd string) []byte {
	return []byte(fmt.Sprintf("-ERR wrong number of arguments for %s command\r\n", cmd))
}

func echo(args []string) []byte {
	if len(args) != 1 {
		return wrongArguments("echo")
	}
	return encodeResponse([]string{args[0]})
}

func command(args []string) []byte {
	if len(args) != 0 {
		return wrongArguments("command")
	}
	return []byte("+OK\r\n")
}

func ping(args []string) []byte {
	if len(args) != 0 {
		return wrongArguments("ping")
	}
	return []byte("+PONG\r\n")
}

func set(args []string) []byte {
	if len(args) != 2 {
		return wrongArguments("set")
	}
	HashMap[args[0]] = args[1]
	return []byte("+OK\r\n")
}

func get(args []string) []byte {
	if len(args) != 1 {
		return wrongArguments("get")
	}
	value, ok := HashMap[args[0]]
	if !ok {
		return encodeResponse(nil)
	}
	return encodeResponse([]string{value})
}

func encodeResponse(response []string) []byte {
	if response == nil {
		return []byte("$-1\r\n")
	}
	var result string
	if len(response) > 1 {
		result = fmt.Sprintf("*%d\r\n", len(response))
	}
	for _, bulk := range response {
		result += fmt.Sprintf("$%d\r\n", len(bulk))
		result += fmt.Sprint(bulk, "\r\n")
	}
	return []byte(result)
}
