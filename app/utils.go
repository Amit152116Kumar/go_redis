package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// type Encoding int
//
// const (
// 	StringEncoding Encoding = iota
// 	ListEncoding
// 	SetEncoding
// 	SortedSetEncoding
// 	HashEncoding
// )

const (
	EOF          byte = 0xFF // End of file
	SELECTDB     byte = 0xFE // Select a database
	EXPIRETIME   byte = 0xFD // Expire time in seconds
	EXPIRETIMEMS byte = 0xFC // Expire time in milliseconds
	RESIZEDB     byte = 0xFB // Resize the database
	AUX          byte = 0xFA // Auxiliary data
)

func isValidDir(path string) bool {
	_, err := filepath.Abs(path)
	return err == nil
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

func writeRDBFile(data []byte) error {
	file, err := os.Create(configSettings.dir + configSettings.dbfilename)
	if err != nil {
		return fmt.Errorf("failed to write data :%v", err)
	}
	writer := bufio.NewWriter(file)

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
