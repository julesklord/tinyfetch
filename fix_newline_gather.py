with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

content = content.replace('lines := strings.Split(out, "\n")', 'lines := strings.Split(out, "\\n")')

with open('cmd/tinyfetch/main.go', 'w') as f:
    f.write(content)
