package main

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
