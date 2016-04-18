// Copyright 2016 NetApp, Inc. All Rights Reserved.

package utils

import (
	"strconv"
	"strings"
)

// lookupTable for suffixes up to power (1024 ^ n)
var lookupTable = make(map[string]int)

func init() {
	// populate the lookup table for suffixes up to power (1024 ^ n)
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
	for k, y := range lookupTable {
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

// GetV takes a map, key, and a defaultValue;  will return the value of the key or defaultValue none set
func GetV(opts map[string]string, key string, defaultValue string) string {
	if value, ok := opts[key]; ok {
		return value
	}
	return defaultValue
}
