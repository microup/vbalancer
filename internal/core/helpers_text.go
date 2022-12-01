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

func FmtStringWithDelimiter(delimiter string, values ...interface{}) *string {
	resultStr := ""
	for _, v := range values {
		value := ""

		if v == nil {
			continue
		}

		//nolint:exhaustive
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			value = fmt.Sprintf("%s", v)
		case reflect.Int:
			value = fmt.Sprintf("%d", v.(int))
		case reflect.Slice:
			s := reflect.ValueOf(v)
			for i := 0; i < s.Len(); i++ {
				val := s.Index(i)
				if value == "" {
					value = fmt.Sprintf("%s", val)
				} else {
					value = value + delimiter + fmt.Sprintf("%s", val)
				}
			}
		default:
			value = fmt.Sprintf("%v", v)
		}

		if resultStr == "" {
			resultStr = value
		} else {
			resultStr = fmt.Sprintf("%s%s%s", resultStr, delimiter, value)
		}
	}

	reg := regexp.MustCompile("\r?\n|\r")
	resultStr = reg.ReplaceAllString(resultStr, " ")

	return &resultStr

}
