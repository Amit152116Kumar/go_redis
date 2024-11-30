package main

import (
	"strings"
	"time"
)

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
	return encodeSimpleString("ok")
}

func ping(args []string) []byte {
	if len(args) != 0 {
		return wrongArguments("ping")
	}
	return encodeSimpleString("pong")
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

func config(args []string) []byte {
	if len(args) != 2 {
		return wrongArguments("config")
	}
	if !strings.EqualFold("get", args[0]) {
		return encodeSimpleError("`" + strings.ToUpper(args[0]) + "` is wrong method")
	}
	var value string
	switch args[1] {
	case "dir":
		value = configSettings.dir
	case "dbfilename":
		value = configSettings.dbfilename
	default:
		return encodeSimpleError("wrong config parameter")
	}
	return encodeBulkString([]string{args[1], value})
}
