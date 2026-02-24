package orderedform

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

// ─── Correctness Tests ───

func TestNewForm(t *testing.T) {
	f := NewForm(5)
	if f == nil {
		t.Fatal("NewForm returned nil")
	}
	got := f.URLEncode()
	if got != "" {
		t.Fatalf("empty form: got %q, want %q", got, "")
	}
}

func TestSetAndURLEncode_Single(t *testing.T) {
	f := NewForm(1)
	f.Set("key", "value")
	got := f.URLEncode()
	want := "key=value"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSetAndURLEncode_Multiple(t *testing.T) {
	f := NewForm(3)
	f.Set("a", "1")
	f.Set("b", "2")
	f.Set("c", "3")
	got := f.URLEncode()
	want := "a=1&b=2&c=3"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestOrderPreserved(t *testing.T) {
	f := NewForm(3)
	f.Set("z", "last")
	f.Set("a", "first")
	f.Set("m", "middle")
	got := f.URLEncode()
	want := "z=last&a=first&m=middle"
	if got != want {
		t.Fatalf("order not preserved: got %q, want %q", got, want)
	}
}

func TestSpecialCharacters(t *testing.T) {
	tests := []struct {
		key, value string
	}{
		{"email", "user@example.com"},
		{"q", "hello world"},
		{"data", "a=b&c=d"},
		{"path", "/foo/bar"},
		{"special", "!@#$%^&*()"},
		{"unicode", "café"},
		{"empty", ""},
		{"", "emptykey"},
		{"spaces", "a b c"},
		{"plus", "a+b"},
		{"percent", "100%"},
		{"brackets", "preferences[theme]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s=%s", tt.key, tt.value), func(t *testing.T) {
			f := NewForm(1)
			f.Set(tt.key, tt.value)
			got := f.URLEncode()
			want := url.QueryEscape(tt.key) + "=" + url.QueryEscape(tt.value)
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

func TestDuplicateKeys(t *testing.T) {
	f := NewForm(3)
	f.Set("k", "v1")
	f.Set("k", "v2")
	f.Set("k", "v3")
	got := f.URLEncode()
	want := "k=v1&k=v2&k=v3"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLargeForm(t *testing.T) {
	n := 1000
	f := NewForm(n)
	var parts []string
	for i := 0; i < n; i++ {
		k := fmt.Sprintf("key%d", i)
		v := fmt.Sprintf("value%d", i)
		f.Set(k, v)
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	got := f.URLEncode()
	want := strings.Join(parts, "&")
	if got != want {
		t.Fatal("large form mismatch")
	}
}

func TestURLEncodeMatchesStdlib(t *testing.T) {
	// Full round-trip: ensure our encoding produces output identical to url.QueryEscape
	f := NewForm(6)
	f.Set("username", "john_doe")
	f.Set("email", "john@example.com")
	f.Set("password", "p@ssw0rd!")
	f.Set("confirm_password", "p@ssw0rd!")
	f.Set("preferences[theme]", "dark")
	f.Set("preferences[notifications]", "email,sms")

	got := f.URLEncode()

	// Build expected using stdlib
	pairs := [][2]string{
		{"username", "john_doe"},
		{"email", "john@example.com"},
		{"password", "p@ssw0rd!"},
		{"confirm_password", "p@ssw0rd!"},
		{"preferences[theme]", "dark"},
		{"preferences[notifications]", "email,sms"},
	}
	var b strings.Builder
	for i, p := range pairs {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(url.QueryEscape(p[0]))
		b.WriteByte('=')
		b.WriteString(url.QueryEscape(p[1]))
	}
	want := b.String()

	if got != want {
		t.Fatalf("README example mismatch:\ngot:  %s\nwant: %s", got, want)
	}
}

func TestZeroCapacity(t *testing.T) {
	f := NewForm(0)
	f.Set("a", "1")
	f.Set("b", "2")
	got := f.URLEncode()
	want := "a=1&b=2"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestReset(t *testing.T) {
	f := NewForm(3)
	f.Set("a", "1")
	f.Set("b", "2")
	f.Reset()

	if got := f.URLEncode(); got != "" {
		t.Fatalf("after Reset: got %q, want %q", got, "")
	}

	// Reuse after reset
	f.Set("x", "10")
	f.Set("y", "20")
	got := f.URLEncode()
	want := "x=10&y=20"
	if got != want {
		t.Fatalf("after reuse: got %q, want %q", got, want)
	}
}

