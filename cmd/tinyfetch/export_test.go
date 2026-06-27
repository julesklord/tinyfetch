package main

import (
	"testing"
)

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty string",
			in:   "",
			want: "",
		},
		{
			name: "no special characters",
			in:   "hello world",
			want: "hello world",
		},
		{
			name: "ampersand",
			in:   "a & b",
			want: "a &amp; b",
		},
		{
			name: "less than",
			in:   "a < b",
			want: "a &lt; b",
		},
		{
			name: "greater than",
			in:   "a > b",
			want: "a &gt; b",
		},
		{
			name: "double quote",
			in:   `"hello"`,
			want: "&quot;hello&quot;",
		},
		{
			name: "single quote",
			in:   `'hello'`,
			want: "&apos;hello&apos;",
		},
		{
			name: "mixed characters",
			in:   `<hello class="world" id='1'>&</hello>`,
			want: `&lt;hello class=&quot;world&quot; id=&apos;1&apos;&gt;&amp;&lt;/hello&gt;`,
		},
		{
			name: "consecutive special characters",
			in:   `<<&>>`,
			want: `&lt;&lt;&amp;&gt;&gt;`,
		},
		{
			name: "unicode characters",
			in:   `world 世界 & "quote"`,
			want: `world 世界 &amp; &quot;quote&quot;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeXML(tt.in); got != tt.want {
				t.Errorf("escapeXML() = %v, want %v", got, tt.want)
			}
		})
	}
}
