package stringutils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func UcFirst(s string) string {
	if len(s) == 0 {
		return ""
	}

	t := []rune(s)
	t[0] = unicode.ToUpper(t[0])

	return string(t)
}

func RandString(n int) string {
	//https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func ConvertToUtf8(s string) string {
	return strings.ToValidUTF8(s, "")
}

func SpacePadRight(s string, length int) string {
	return fmt.Sprintf("%-[1]*[2]s", length, s)
}

func SpacePadLeft(s string, length int) string {
	return fmt.Sprintf("%[1]*[2]s", length, s)
}
