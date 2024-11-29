package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func wrongArguments(cmd string) []byte {
	return []byte(fmt.Sprintf("-ERR wrong number of arguments for %s command\r\n", cmd))
}

func echo(args []string) []byte {
	if len(args) != 1 {
		return wrongArguments("echo")
	}
	return encodeBulkString([]string{args[0]})
}

func command(args []string) []byte {
	if len(args) != 0 {
		return wrongArguments("command")
	}
	return encodeSimpleString("OK")
}

func ping(args []string) []byte {
	if len(args) != 0 {
		return wrongArguments("ping")
	}
	return encodeSimpleString("PONG")
}

func containsExpiry(slice []string) int64 {
	for i, element := range slice {
		if strings.EqualFold(element, "px") {
			exp, _ := strconv.ParseInt(slice[i+1], 10, 64)
			return exp + time.Now().Local().UnixMilli()
		}
		if strings.EqualFold(element, "ex") {
			exp, _ := strconv.ParseInt(slice[i+1], 10, 64)
			return exp*1000 + time.Now().Local().UnixMilli()
		}
	}
	return 0
}

func set(args []string) []byte {
	if len(args)&1 == 1 {
		return wrongArguments("set")
	}
	value := &Value{}
	value.val = args[1]
	value.expiry = containsExpiry(args[2:])
	HashMap[args[0]] = *value
	return encodeSimpleString("OK")
}

func get(args []string) []byte {
	if len(args) != 1 {
		return wrongArguments("get")
	}
	value, ok := HashMap[args[0]]
	if !ok {
		return encodeBulkString(nil)
	}
	if value.expiry == 0 || time.Now().UnixMilli() < value.expiry {
		return encodeBulkString([]string{value.val})
	}
	delete(HashMap, args[0])
	return encodeBulkString(nil)
}

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
	return []byte(fmt.Sprint("+", res, "\r\n"))
}

func encodeSimpleError(err string) []byte {
	return []byte(fmt.Sprint("-Err ", err, "\r\n"))
}
