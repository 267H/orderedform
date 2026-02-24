package orderedform

import "unsafe"

const upperhex = "0123456789ABCDEF"

// noEscape[c] is true if byte c does not need query escaping.
// Matches net/url.QueryEscape behavior exactly:
//   - a-z, A-Z, 0-9: pass through
//   - '-', '_', '.', '~': pass through (RFC 3986 unreserved)
//   - ' ': encoded as '+' (handled separately)
//   - everything else: percent-encoded as %XX
var noEscape [256]bool

func init() {
	for c := 'a'; c <= 'z'; c++ {
		noEscape[c] = true
	}
	for c := 'A'; c <= 'Z'; c++ {
		noEscape[c] = true
	}
	for c := '0'; c <= '9'; c++ {
		noEscape[c] = true
	}
	noEscape['-'] = true
	noEscape['_'] = true
	noEscape['.'] = true
	noEscape['~'] = true
}

type Form struct {
	pairs [][2]string
}

func NewForm(capacity int) *Form {
	return &Form{
		pairs: make([][2]string, 0, capacity),
	}
}

// Reset clears all pairs while retaining the allocated capacity.
// Allows reuse of the Form to avoid repeated slice allocation in hot paths.
func (f *Form) Reset() {
	f.pairs = f.pairs[:0]
}

func (f *Form) Set(k, v string) {
	f.pairs = append(f.pairs, [2]string{k, v})
}

func (f *Form) URLEncode() string {
	if len(f.pairs) == 0 {
		return ""
	}

	// First pass: calculate exact output size
	n := len(f.pairs) - 1 // '&' separators
	for _, p := range f.pairs {
		n += escapedLen(p[0]) + 1 + escapedLen(p[1])
	}

	// Single allocation, exact size
	buf := make([]byte, n)
	pos := 0
	for i, p := range f.pairs {
		if i > 0 {
			buf[pos] = '&'
			pos++
		}
		pos = writeEscaped(buf, pos, p[0])
		buf[pos] = '='
		pos++
		pos = writeEscaped(buf, pos, p[1])
	}

	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

// escapedLen returns the length of s when query-escaped.
func escapedLen(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if noEscape[c] {
			n++
		} else if c == ' ' {
			n++
		} else {
			n += 3
		}
	}
	return n
}

// writeEscaped writes query-escaped s into buf at pos, returns new pos.
func writeEscaped(buf []byte, pos int, s string) int {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if noEscape[c] {
			buf[pos] = c
			pos++
		} else if c == ' ' {
			buf[pos] = '+'
			pos++
		} else {
			buf[pos] = '%'
			buf[pos+1] = upperhex[c>>4]
			buf[pos+2] = upperhex[c&0x0F]
			pos += 3
		}
	}
	return pos
}
