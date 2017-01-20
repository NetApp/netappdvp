// Copyright 2016 NetApp, Inc. All Rights Reserved.

package utils

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// lookupTable for suffixes up to power (1024 ^ n)
var lookupTable = make(map[string]int)
var units sizeUnit = sizeUnit{}

type sizeUnit []string

func (s sizeUnit) Len() int {
	return len(s)
}
func (s sizeUnit) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sizeUnit) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func init() {
	// populate the lookup table for suffixes up to power (1024 ^ n)
	lookupTable["b"] = 0
	lookupTable["bytes"] = 0
	lookupTable["k"] = 1
	lookupTable["kb"] = 1
	lookupTable["m"] = 2
	lookupTable["mb"] = 2
	lookupTable["g"] = 3
	lookupTable["gb"] = 3
	lookupTable["t"] = 4
	lookupTable["tb"] = 4
	lookupTable["p"] = 5
	lookupTable["pb"] = 5
	lookupTable["e"] = 6
	lookupTable["eb"] = 6
	lookupTable["z"] = 7
	lookupTable["zb"] = 7
	lookupTable["y"] = 8
	lookupTable["yb"] = 8

	//The slice of units is used to ensure that they are accessed by suffix from the largest length to the shortest
	for k, _ := range lookupTable {
		units = append(units, k)
	}
	sort.Sort(units)
}

// Linux is a constant value for the runtime.GOOS that represents the Linux OS
const Linux = "linux"

// Windows is a constant value for the runtime.GOOS that represents the Windows OS
const Windows = "windows"

// Darwin is a constant value for the runtime.GOOS that represents Apple OSX
const Darwin = "darwin"

// Pow is an integer version of exponentation; existing builtin is float, we needed an int version
func Pow(x, y int) int {
	if y == 0 {
		return 1
	}

	result := x
	for n := 1; n < y; n++ {
		result = result * x
	}
	return result
}

// ConvertSizeToBytes converts size to bytes; see also https://en.wikipedia.org/wiki/Kilobyte
func ConvertSizeToBytes(s string) (string, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	// spin until we find a match, if no match return original string
	for _, k := range units {
		var y int = lookupTable[k]
		if strings.HasSuffix(s, k) {
			s = s[:len(s)-len(k)]
			i, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			i = i * Pow(1024, y)
			s = strconv.Itoa(i)
			return s, nil
		}
	}

	return s, nil
}

// Pow is an integer version of exponentation; existing builtin is float, we needed an int version
func Pow64(x int64, y int) int64 {
	if y == 0 {
		return 1
	}

	result := x
	for n := 1; n < y; n++ {
		result = result * x
	}
	return result
}

// ConvertSizeToBytes64 converts size to bytes; see also https://en.wikipedia.org/wiki/Kilobyte
func ConvertSizeToBytes64(s string) (string, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	// spin until we find a match, if no match return original string
	for _, k := range units {
		var y int = lookupTable[k]
		if strings.HasSuffix(s, k) {
			s = s[:len(s)-len(k)]
			i, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				return "", err
			}
			i = i * Pow64(1024, y)
			s = strconv.FormatInt(i, 10)
			return s, nil
		}
	}

	return s, nil
}

// GetV takes a map, key, and a defaultValue;  will return the value of the key or defaultValue none set
func GetV(opts map[string]string, key string, defaultValue string) string {
	if value, ok := opts[key]; ok {
		return value
	}
	return defaultValue
}

func RandomNumericString(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
