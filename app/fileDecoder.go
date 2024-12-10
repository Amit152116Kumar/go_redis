package main

import (
	"fmt"
	"strconv"
)

func getValueType(flag byte) func() {
	switch flag {
	case 0:
		return StringDecoding
	case 1:
		return ListDecoding
	case 2:
		return SetDecoding
	case 4:
		return HashDecoding
	}
	return nil
}

type SpecialFormat byte

const (
	IntegerAsString SpecialFormat = iota + 1
	CompressedStrings
)

func LengthDecoding(data []byte) (int, SpecialFormat, []byte) {
	bitmask := byte(3 << 6)
	mst := (bitmask & data[0]) >> 6
	var length int
	format := SpecialFormat(0)
	switch mst {
	case 0:
		length = int(data[0])
		data = data[1:]
	case 0b01:
		length = int(data[1])
		length = int(data[0]-mst)<<8 | length
		data = data[2:]
	case 0b10:
		for i := 1; i < 5; i++ {
			length = length<<8 | int(data[i])
		}
		data = data[5:]
	case 0b11:
		value := data[0] ^ (mst << 6)
		switch value {
		case 0:
			length = 8
		case 1:
			length = 16
		case 2:
			length = 32
		}
		data = data[1:]
		format = IntegerAsString
	}
	return length, format, data
}

func StringDecoding() {
	// Implementation for decoding a string
}

func ListDecoding() {
	// Implementation for decoding a list
}

func SetDecoding() {
	// Implementation for decoding a set
}

func HashDecoding() {
	// Implementation for decoding a hash
}

func parseAUX(data []byte) {
	fmt.Println("\n", string(data))
	fmt.Println("\nAUX  -> ")
	fmt.Println(data, "\n")

	if GlobalRDB.metadata == nil {
		GlobalRDB.metadata = make(map[string]string)
	}
	len, _, data := LengthDecoding(data)
	key := string(data[:len])
	data = data[len:]
	len, format, data := LengthDecoding(data)
	var value string
	if format == IntegerAsString {
		value = strconv.FormatUint(convertBytesToINT(data), 10)
	} else {
		value = string(data[:len])
	}
	GlobalRDB.metadata[key] = value
}

func parseDB(data []byte) {
	fmt.Println("\n", string(data))
	fmt.Println("\n DBSelector -> ")
	fmt.Println(data, "\n")

	if GlobalRDB.database == nil {
		GlobalRDB.database = &Database{}
	}
	GlobalRDB.database.dbNumber = int(convertBytesToINT(data))
}

func parseResizeDB(data []byte) {
	fmt.Println("\n", string(data))
	fmt.Println("\nRESIZEDB -> ")
	fmt.Println(data, "\n")

	cursor := 0
	// numKeys := int(data[cursor])
	cursor++

	len, format, data := LengthDecoding(data)
	fmt.Println("Length : ", len)
	var value string
	if format == IntegerAsString {
		value = strconv.FormatUint(convertBytesToINT(data[:len]), 10)
	} else {
		value = string(data[:len])
		data = data[len:]
	}
	fmt.Println("VAlue -> ", value)
}

func parseExpireTime(data []byte) {
	fmt.Println("\n", string(data))
	fmt.Println("\nExpireTime -> ")
	fmt.Println(data, "\n")
}

func parseExpireTimeMS(data []byte) {
	fmt.Println("\n", string(data))
	fmt.Println("\nExpireTimeMS -> ")
	fmt.Println(data, "\n")
}

func parseOpCodeData(data []byte, opCode byte) {
	// fmt.Println(opCode, " -> ", data)
	switch opCode {
	case SELECTDB:
		parseDB(data)
	case AUX:
		parseAUX(data)
	case EOF:
		fmt.Println("\n", string(data))
		fmt.Println("\nEOF -> ")
		fmt.Println(data, "\n")

	case EXPIRETIME:
		parseExpireTime(data)
	case EXPIRETIMEMS:
		parseExpireTimeMS(data)
	case RESIZEDB:
		parseResizeDB(data)
	default:
		header := string(data)
		GlobalRDB.header = header
	}
}

func convertBytesToINT(data []byte) uint64 {
	var u uint64
	for i := 0; i < len(data); i++ {
		u = (u << 8) | uint64(data[i]) // Shift left and add the next byte
	}
	return u
}
