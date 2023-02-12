package core

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"
)

// TrimLastChar removes the last character from a string.
func TrimLastChar(value string) string {
	// If the string is not empty.
	if len(value) > 0 {
		// Get the last rune and its size.
		_, size := utf8.DecodeLastRuneInString(value)
		// If the size is the same as the length of the string.
		if size == len(value) {
			// Return an empty string.
			return ""
		}
		// Return the string without the last rune.
		return value[:len(value)-size]
	}

	// Return an empty string.
	return ""
}

// FmtStringWithDelimiter formats a string with a delimiter.
func FmtStringWithDelimiter(delimiter string, values ...interface{}) string {
	resultStr := ""

	for iLen, value := range values {
		if value == nil {
			continue
		}

		appendLastChar := true

		if iLen+1 == len(values) {
			valueType := reflect.TypeOf(value).Kind()
			if valueType == reflect.Slice {
				s := reflect.ValueOf(value)
				if s.Len() == 0 {
					appendLastChar = false
				}
			}
		}

		if appendLastChar {
			valueStr := getValueStr(delimiter, value)
			if resultStr == "" {
				resultStr = valueStr
			} else {
				resultStr = fmt.Sprintf("%s%s%s", resultStr, delimiter, valueStr)
			}
		}
	}

	reg := regexp.MustCompile("\r?\n|\r")
	resultStr = reg.ReplaceAllString(resultStr, " ")

	return resultStr
}

// getValueStr - convert a value of any type into a string representation, depending on its type.
func getValueStr(delimiter string, value interface{}) string {
	valueStr := ""

	//nolint:exhaustive
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		valueStr = fmt.Sprintf("%s", value)
	case reflect.Int:
		valueStr = fmt.Sprintf("%d", value)
	case reflect.Slice:
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			val := s.Index(i)
			if valueStr == "" {
				valueStr = fmt.Sprintf("%v", val.Interface())
			} else {
				valueStr = valueStr + delimiter + fmt.Sprintf("%v", val.Interface())
			}
		}
	default:
		valueStr = fmt.Sprintf("%v", value)
	}

	return valueStr
}
