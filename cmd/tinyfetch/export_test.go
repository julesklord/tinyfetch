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
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"<tag>", "&lt;tag&gt;"},
		{"&amp;", "&amp;amp;"},
		{"\"quotes\"", "&quot;quotes&quot;"},
		{"'apos'", "&apos;apos&apos;"},
		{"&<>'\"", "&amp;&lt;&gt;&apos;&quot;"},
	}

	for _, test := range tests {
		result := escapeXML(test.input)
		if result != test.expected {
			t.Errorf("escapeXML(%q) = %q, expected %q", test.input, result, test.expected)
		}
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
	output := captureStdout(func() {
		printJSON("myhost", "myos", "mykernel", "myuptime", "myshell", "mycpu", "mymem", "mydisk", []string{"key1", "key2"}, []string{"val1", "val2"})
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
	output := captureStdout(func() {
		printJSON("myhost", "myos", "mykernel", "myuptime", "myshell", "mycpu", "mymem", "mydisk", []string{}, []string{})
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
	output := captureStdout(func() {
		printXML("myhost", "myos", "mykernel", "myuptime", "myshell", "mycpu", "mymem", "mydisk", []string{"key1", "key 2"}, []string{"val1", "val2"})
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
	output := captureStdout(func() {
		printXML("myhost", "myos", "mykernel", "myuptime", "myshell", "mycpu", "mymem", "mydisk", []string{}, []string{})
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
	output := captureStdout(func() {
		printTXT("myhost", "myos", "mykernel", "myuptime", "myshell", "mycpu", "mymem", "mydisk", []string{"key1", "key 2"}, []string{"val1", "val2"})
	})

	if output != expected {
		t.Errorf("printTXT output mismatch.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}
