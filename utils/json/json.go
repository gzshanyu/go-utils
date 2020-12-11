package json

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"unicode"
)

// 小驼峰转下划线
func Hump2Underline(in interface{}, out interface{}) error {
	var (
		err              error
		byteResp         []byte
		keyMatchRegex    *regexp.Regexp
		wordBarrierRegex *regexp.Regexp
		converted        []byte
	)

	keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)
	if byteResp, err = json.Marshal(in); err != nil {
		log.Print("Underline2Hump err1:", err)
		return err
	}

	converted = keyMatchRegex.ReplaceAllFunc(
		byteResp,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)

	if err = json.Unmarshal(converted, &out); err != nil {
		log.Print("Underline2Hump err2:", err)
		return err
	}

	return nil
}

// 下划线转小驼峰
func Underline2Hump(in interface{}, out interface{}) error {
	var (
		err           error
		byteResp      []byte
		keyMatchRegex *regexp.Regexp
		converted     []byte
	)

	keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	if byteResp, err = json.Marshal(in); err != nil {
		return err
	}
	converted = keyMatchRegex.ReplaceAllFunc(
		byteResp,
		func(match []byte) []byte {
			matchStr := string(match)
			key := matchStr[1 : len(matchStr)-2]
			resKey := LcFirst(Case2Camel(key))
			return []byte(`"` + resKey + `":`)
		},
	)

	if err = json.Unmarshal(converted, &out); err != nil {
		return err
	}

	return nil
}

// 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 首字母小写
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// json转字符串
func JsonEncode(o interface{}) string {
	var data, _ = json.Marshal(o)
	return string(data)
}

// 字符串转json
func JsonDecode(s string) interface{} {
	var out interface{}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}
