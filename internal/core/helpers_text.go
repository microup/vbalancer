package core

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"
)

func TrimLastChar(value string) string {
	for len(value) > 0 {
		_, size := utf8.DecodeLastRuneInString(value)

		return value[:len(value)-size]
	}

	return value
}

//nolint:cyclop
func FmtStringWithDelimiter(delimiter string, values ...interface{}) *string {
	resultStr := ""

	for _, value := range values {
		valueStr := ""

		if value == nil {
			continue
		}

		//nolint:exhaustive
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			valueStr = fmt.Sprintf("%s", value)
		case reflect.Int:
			valueInt, ok := value.(int)
			if !ok {
				continue
			}

			valueStr = fmt.Sprintf("%d", valueInt)
		case reflect.Slice:
			s := reflect.ValueOf(value)
			for i := 0; i < s.Len(); i++ {
				val := s.Index(i)
				if valueStr == "" {
					valueStr = fmt.Sprintf("%s", val)
				} else {
					valueStr = valueStr + delimiter + fmt.Sprintf("%s", val)
				}
			}
		default:
			valueStr = fmt.Sprintf("%v", value)
		}

		if resultStr == "" {
			resultStr = valueStr
		} else {
			resultStr = fmt.Sprintf("%s%s%s", resultStr, delimiter, valueStr)
		}
	}

	reg := regexp.MustCompile("\r?\n|\r")
	resultStr = reg.ReplaceAllString(resultStr, " ")

	return &resultStr
}
