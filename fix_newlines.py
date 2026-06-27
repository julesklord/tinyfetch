with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

content = content.replace('fmt.Fprintf(os.Stderr, "Unknown output format: %s\n", outputFmt)', 'fmt.Fprintf(os.Stderr, "Unknown output format: %s\\n", outputFmt)')
content = content.replace('lcyan + "     /     \\" + restore', 'lcyan + "     /     \\\\" + restore')
content = content.replace('lcyan + "     \\__   /" + restore', 'lcyan + "     \\\\__   /" + restore')
content = content.replace('lcyan + "    /   `-\' \\" + restore', 'lcyan + "    /   `-\' \\\\" + restore')
content = content.replace('lcyan + "    \\       /" + restore', 'lcyan + "    \\\\       /" + restore')
content = content.replace('lyellow + "    /     \\" + restore', 'lyellow + "    /     \\\\" + restore')
content = content.replace('lblue + "    \\ " + restore', 'lblue + "    \\\\ " + restore')
content = content.replace('lyellow + "    /  \\-/ \\" + restore', 'lyellow + "    /  \\\\-/ \\\\" + restore')
content = content.replace('lyellow + "   / /     \\ \\" + restore', 'lyellow + "   / /     \\\\ \\\\" + restore')

with open('cmd/tinyfetch/main.go', 'w') as f:
    f.write(content)
