package orderedform

import (
	"net/url"
	"strings"
)

type Form struct {
	pairs [][2]string
}

func NewForm(capacity int) *Form {
	return &Form{
		pairs: make([][2]string, 0, capacity),
	}
}

func (f *Form) Set(k, v string) {
	f.pairs = append(f.pairs, [2]string{url.QueryEscape(k), url.QueryEscape(v)})
}

func (f *Form) URLEncode() string {
	var b strings.Builder
	for i, p := range f.pairs {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(p[0])
		b.WriteByte('=')
		b.WriteString(p[1])
	}
	return b.String()
}
