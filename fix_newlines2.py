with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

content = content.replace('lines := strings.Split(rawOut, "\n")', 'lines := strings.Split(rawOut, "\\n")')
content = content.replace('fmt.Printf("%s│%s %s%s %s│\n", borderCol, restore, rLine, rPadding, borderCol)', 'fmt.Printf("%s│%s %s%s %s│\\n", borderCol, restore, rLine, rPadding, borderCol)')
content = content.replace('fmt.Printf("%s│%s %s%s %s│%s %s%s %s│\n",', 'fmt.Printf("%s│%s %s%s %s│%s %s%s %s│\\n",')
content = content.replace('fmt.Printf("%s│%s %s%s %s│%s %s%s %s│%s %s%s %s│\n",', 'fmt.Printf("%s│%s %s%s %s│%s %s%s %s│%s %s%s %s│\\n",')
content = content.replace('fmt.Printf("%s│%s %s%s %s│\n", borderCol, restore, printLine, padStr, borderCol)', 'fmt.Printf("%s│%s %s%s %s│\\n", borderCol, restore, printLine, padStr, borderCol)')

with open('cmd/tinyfetch/main.go', 'w') as f:
    f.write(content)
