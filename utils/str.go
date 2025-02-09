package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var ansiEscapeRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func RemoveNotValidUtf8InString(s string) string {
	ret := s
	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for i, r := range s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		ret = string(v)
	}
	return ret
}

func RemoveANSI(input string) string {
	return ansiEscapeRegex.ReplaceAllString(input, "")
}

func RandString(n int) (ret string) {
	allString := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789"
	ret = ""
	for i := 0; i < n; i++ {
		r := rand.Intn(len(allString))
		ret = ret + allString[r:r+1]
	}
	return
}

type KVStr struct {
	m map[string]any
}

func NewKVStr() *KVStr {
	return &KVStr{
		map[string]any{},
	}
}

func (k *KVStr) Add(key string, value any) *KVStr {
	k.m[key] = value
	return k
}

func (k *KVStr) Build() string {
	if len(k.m) == 0 {
		return "-"
	}
	keys := make([]string, 0, len(k.m))
	for key := range k.m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	strList := make([]string, 0, len(k.m))
	for _, key := range keys {
		if k.m[key] == "" {
			continue
		}
		strList = append(strList, fmt.Sprintf("%s:%v", key, k.m[key]))
	}
	return strings.Join(strList, " , ")
}

func JsonStrToStruct[T any](str string) T {
	var data T
	json.Unmarshal([]byte(str), &data)
	return data
}
