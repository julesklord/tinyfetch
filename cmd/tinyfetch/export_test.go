package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	func() {
		defer func() {
			w.Close()
			os.Stdout = old
		}()
		f()
	}()
	return <-outC
}

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

func TestPrintJSON(t *testing.T) {
	expected := `{
  "host": "myhost",
  "os": "myos",
  "kernel": "mykernel",
  "uptime": "myuptime",
  "shell": "myshell",
  "cpu": "mycpu",
  "memory": "mymem",
  "disk": "mydisk",
  "plugins": {
    "key1": "val1",
    "key2": "val2"
  }
}
`
	info := SystemInfo{
		Host:   "myhost",
		OSName: "myos",
		Kernel: "mykernel",
		Uptime: "myuptime",
		Shell:  "myshell",
		CPU:    "mycpu",
		Memory: "mymem",
		Disk:   "mydisk",
		Keys:   []string{"key1", "key2"},
		Vals:   []string{"val1", "val2"},
	}
	output := captureStdout(func() {
		printJSON(info)
	})

	if output != expected {
		t.Errorf("printJSON output mismatch.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPrintJSON_NoPlugins(t *testing.T) {
	expected := `{
  "host": "myhost",
  "os": "myos",
  "kernel": "mykernel",
  "uptime": "myuptime",
  "shell": "myshell",
  "cpu": "mycpu",
  "memory": "mymem",
  "disk": "mydisk"
}
`
	info := SystemInfo{
		Host:   "myhost",
		OSName: "myos",
		Kernel: "mykernel",
		Uptime: "myuptime",
		Shell:  "myshell",
		CPU:    "mycpu",
		Memory: "mymem",
		Disk:   "mydisk",
		Keys:   []string{},
		Vals:   []string{},
	}
	output := captureStdout(func() {
		printJSON(info)
	})

	if output != expected {
		t.Errorf("printJSON output mismatch.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPrintXML(t *testing.T) {
	expected := `<tinyfetch>
  <host>myhost</host>
  <os>myos</os>
  <kernel>mykernel</kernel>
  <uptime>myuptime</uptime>
  <shell>myshell</shell>
  <cpu>mycpu</cpu>
  <memory>mymem</memory>
  <disk>mydisk</disk>
  <plugins>
    <key1>val1</key1>
    <key_2>val2</key_2>
  </plugins>
</tinyfetch>
`
	info := SystemInfo{
		Host:   "myhost",
		OSName: "myos",
		Kernel: "mykernel",
		Uptime: "myuptime",
		Shell:  "myshell",
		CPU:    "mycpu",
		Memory: "mymem",
		Disk:   "mydisk",
		Keys:   []string{"key1", "key 2"},
		Vals:   []string{"val1", "val2"},
	}
	output := captureStdout(func() {
		printXML(info)
	})

	if output != expected {
		t.Errorf("printXML output mismatch.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPrintXML_NoPlugins(t *testing.T) {
	expected := `<tinyfetch>
  <host>myhost</host>
  <os>myos</os>
  <kernel>mykernel</kernel>
  <uptime>myuptime</uptime>
  <shell>myshell</shell>
  <cpu>mycpu</cpu>
  <memory>mymem</memory>
  <disk>mydisk</disk>
</tinyfetch>
`
	info := SystemInfo{
		Host:   "myhost",
		OSName: "myos",
		Kernel: "mykernel",
		Uptime: "myuptime",
		Shell:  "myshell",
		CPU:    "mycpu",
		Memory: "mymem",
		Disk:   "mydisk",
		Keys:   []string{},
		Vals:   []string{},
	}
	output := captureStdout(func() {
		printXML(info)
	})

	if output != expected {
		t.Errorf("printXML output mismatch.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPrintTXT(t *testing.T) {
	expected := `Host: myhost
OS: myos
Kernel: mykernel
Uptime: myuptime
Shell: myshell
CPU: mycpu
Memory: mymem
Disk: mydisk
key1: val1
key 2: val2
`
	info := SystemInfo{
		Host:   "myhost",
		OSName: "myos",
		Kernel: "mykernel",
		Uptime: "myuptime",
		Shell:  "myshell",
		CPU:    "mycpu",
		Memory: "mymem",
		Disk:   "mydisk",
		Keys:   []string{"key1", "key 2"},
		Vals:   []string{"val1", "val2"},
	}
	output := captureStdout(func() {
		printTXT(info)
	})

	if output != expected {
		t.Errorf("printTXT output mismatch.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}
