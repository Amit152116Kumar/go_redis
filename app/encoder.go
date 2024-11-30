package main

import (
	"fmt"
	"strings"
)

func encodeBulkString(response []string) []byte {
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

func encodeSimpleString(res string) []byte {
	return []byte(fmt.Sprint("+", strings.ToUpper(res), "\r\n"))
}

func encodeSimpleError(err string) []byte {
	return []byte(fmt.Sprint("-Err ", err, "\r\n"))
}

func wrongArguments(cmd string) []byte {
	return []byte(fmt.Sprintf("-ERR wrong number of arguments for `%s` command\r\n", strings.ToLower(cmd)))
}
