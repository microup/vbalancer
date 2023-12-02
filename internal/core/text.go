package core

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"
)

// TrimLastChar removes the last character from a string.
func TrimLastChar(value string) string {
	if len(value) > 0 {
		_, size := utf8.DecodeLastRuneInString(value)
		if size == len(value) {
			return ""
		}

		return value[:len(value)-size]
	}

	return ""
}

// FmtStringWithDelimiter formats a string with a given delimiter.
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

// getValueStr returns a string representation of a value.
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
