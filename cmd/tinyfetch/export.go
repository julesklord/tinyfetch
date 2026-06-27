package main

import (
	"fmt"
	"strings"
)

type SystemInfo struct {
	Host   string
	OSName string
	Kernel string
	Uptime string
	Shell  string
	CPU    string
	Memory string
	Disk   string
	Keys   []string
	Vals   []string
}

func printJSON(info SystemInfo) {
	fmt.Printf("{\n")
	fmt.Printf("  \"host\": %q,\n", info.Host)
	fmt.Printf("  \"os\": %q,\n", info.OSName)
	fmt.Printf("  \"kernel\": %q,\n", info.Kernel)
	fmt.Printf("  \"uptime\": %q,\n", info.Uptime)
	fmt.Printf("  \"shell\": %q,\n", info.Shell)
	fmt.Printf("  \"cpu\": %q,\n", info.CPU)
	fmt.Printf("  \"memory\": %q,\n", info.Memory)
	fmt.Printf("  \"disk\": %q", info.Disk)

	if len(info.Keys) > 0 {
		fmt.Printf(",\n  \"plugins\": {\n")
		for i := 0; i < len(info.Keys); i++ {
			cleanVal := stripANSI(info.Vals[i])
			fmt.Printf("    %q: %q", info.Keys[i], cleanVal)
			if i < len(info.Keys)-1 {
				fmt.Printf(",\n")
			} else {
				fmt.Printf("\n")
			}
		}
		fmt.Printf("  }\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("}\n")
}

func escapeXML(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			sb.WriteString("&amp;")
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		case '"':
			sb.WriteString("&quot;")
		case '\'':
			sb.WriteString("&apos;")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func printXML(info SystemInfo) {
	fmt.Printf("<tinyfetch>\n")
	fmt.Printf("  <host>%s</host>\n", escapeXML(info.Host))
	fmt.Printf("  <os>%s</os>\n", escapeXML(info.OSName))
	fmt.Printf("  <kernel>%s</kernel>\n", escapeXML(info.Kernel))
	fmt.Printf("  <uptime>%s</uptime>\n", escapeXML(info.Uptime))
	fmt.Printf("  <shell>%s</shell>\n", escapeXML(info.Shell))
	fmt.Printf("  <cpu>%s</cpu>\n", escapeXML(info.CPU))
	fmt.Printf("  <memory>%s</memory>\n", escapeXML(info.Memory))
	fmt.Printf("  <disk>%s</disk>\n", escapeXML(info.Disk))
	if len(info.Keys) > 0 {
		fmt.Printf("  <plugins>\n")
		for i := 0; i < len(info.Keys); i++ {
			tag := strings.ToLower(info.Keys[i])
			var sb strings.Builder
			for _, r := range tag {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
					sb.WriteRune(r)
				} else {
					sb.WriteRune('_')
				}
			}
			tagStr := sb.String()
			cleanVal := stripANSI(info.Vals[i])
			fmt.Printf("    <%s>%s</%s>\n", tagStr, escapeXML(cleanVal), tagStr)
		}
		fmt.Printf("  </plugins>\n")
	}
	fmt.Printf("</tinyfetch>\n")
}

func printTXT(info SystemInfo) {
	fmt.Printf("Host: %s\n", info.Host)
	fmt.Printf("OS: %s\n", info.OSName)
	fmt.Printf("Kernel: %s\n", info.Kernel)
	fmt.Printf("Uptime: %s\n", info.Uptime)
	fmt.Printf("Shell: %s\n", info.Shell)
	fmt.Printf("CPU: %s\n", info.CPU)
	fmt.Printf("Memory: %s\n", info.Memory)
	fmt.Printf("Disk: %s\n", info.Disk)
	for i := 0; i < len(info.Keys); i++ {
		fmt.Printf("%s: %s\n", info.Keys[i], stripANSI(info.Vals[i]))
	}
}
