package main

import (
	"fmt"
	"strings"
)

type PluginInfo struct {
	Key     string
	Val     string
	Details []string
}

type SystemInfo struct {
	Host      string
	OSName    string
	Kernel    string
	Uptime    string
	Shell     string
	CPU       string
	GPU       string
	DEWM      string
	Terminal  string
	Memory    string
	Swap      string
	Disk      string
	Processes string
	Plugins   []PluginInfo
}

func printJSON(info SystemInfo) {
	fmt.Printf("{\n")
	fmt.Printf("  \"host\": %q,\n", info.Host)
	fmt.Printf("  \"os\": %q,\n", info.OSName)
	fmt.Printf("  \"kernel\": %q,\n", info.Kernel)
	fmt.Printf("  \"uptime\": %q,\n", info.Uptime)
	fmt.Printf("  \"shell\": %q,\n", info.Shell)
	fmt.Printf("  \"cpu\": %q,\n", info.CPU)
	fmt.Printf("  \"gpu\": %q,\n", info.GPU)
	fmt.Printf("  \"de_wm\": %q,\n", info.DEWM)
	fmt.Printf("  \"terminal\": %q,\n", info.Terminal)
	fmt.Printf("  \"memory\": %q,\n", info.Memory)
	fmt.Printf("  \"swap\": %q,\n", info.Swap)
	fmt.Printf("  \"disk\": %q,\n", info.Disk)
	fmt.Printf("  \"processes\": %q", info.Processes)

	if len(info.Plugins) > 0 {
		fmt.Printf(",\n  \"plugins\": {\n")
		for i, plug := range info.Plugins {
			cleanKey := strings.ToLower(plug.Key)
			cleanVal := stripANSI(plug.Val)
			fmt.Printf("    %q: {\n", cleanKey)
			fmt.Printf("      \"value\": %q", cleanVal)
			if len(plug.Details) > 0 {
				fmt.Printf(",\n      \"details\": [\n")
				for j, det := range plug.Details {
					fmt.Printf("        %q", stripANSI(det))
					if j < len(plug.Details)-1 {
						fmt.Printf(",\n")
					} else {
						fmt.Printf("\n")
					}
				}
				fmt.Printf("      ]\n")
			} else {
				fmt.Printf("\n")
			}
			fmt.Printf("    }")
			if i < len(info.Plugins)-1 {
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
	fmt.Printf("  <gpu>%s</gpu>\n", escapeXML(info.GPU))
	fmt.Printf("  <de_wm>%s</de_wm>\n", escapeXML(info.DEWM))
	fmt.Printf("  <terminal>%s</terminal>\n", escapeXML(info.Terminal))
	fmt.Printf("  <memory>%s</memory>\n", escapeXML(info.Memory))
	fmt.Printf("  <swap>%s</swap>\n", escapeXML(info.Swap))
	fmt.Printf("  <disk>%s</disk>\n", escapeXML(info.Disk))
	fmt.Printf("  <processes>%s</processes>\n", escapeXML(info.Processes))

	if len(info.Plugins) > 0 {
		fmt.Printf("  <plugins>\n")
		for _, plug := range info.Plugins {
			tag := strings.ToLower(plug.Key)
			var sb strings.Builder
			for _, r := range tag {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
					sb.WriteRune(r)
				} else {
					sb.WriteRune('_')
				}
			}
			tagStr := sb.String()
			cleanVal := stripANSI(plug.Val)
			fmt.Printf("    <%s>\n", tagStr)
			fmt.Printf("      <value>%s</value>\n", escapeXML(cleanVal))
			for _, det := range plug.Details {
				fmt.Printf("      <detail>%s</detail>\n", escapeXML(stripANSI(det)))
			}
			fmt.Printf("    </%s>\n", tagStr)
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
	fmt.Printf("GPU: %s\n", info.GPU)
	fmt.Printf("DE/WM: %s\n", info.DEWM)
	fmt.Printf("Terminal: %s\n", info.Terminal)
	fmt.Printf("Memory: %s\n", info.Memory)
	fmt.Printf("Swap: %s\n", info.Swap)
	fmt.Printf("Disk: %s\n", info.Disk)
	fmt.Printf("Processes: %s\n", info.Processes)
	for _, plug := range info.Plugins {
		fmt.Printf("%s: %s\n", plug.Key, stripANSI(plug.Val))
		for _, det := range plug.Details {
			fmt.Printf("  %s\n", stripANSI(det))
		}
	}
}
