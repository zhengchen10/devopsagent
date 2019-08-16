package utils

import "strings"

type StringTools struct {
}

func (st *StringTools) HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func (st *StringTools) HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func (st *StringTools) Trim(str string) string {
	ret := strings.Trim(str, " ")
	if len(ret) != len(str) {
		return st.Trim(ret)
	}
	ret = strings.Trim(str, "\n")
	if len(ret) != len(str) {
		return st.Trim(ret)
	}
	ret = strings.Trim(str, "\t")
	if len(ret) != len(str) {
		return st.Trim(ret)
	}
	ret = strings.Trim(str, "\r")
	if len(ret) != len(str) {
		return st.Trim(ret)
	}
	ret = strings.Trim(str, "")
	if len(ret) != len(str) {
		return st.Trim(ret)
	}
	return ret
}

func (st *StringTools) SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}
