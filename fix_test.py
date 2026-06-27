with open('tests/test.sh', 'r') as f:
    content = f.read()

content = content.replace('line_count=$(echo "$no_ascii_out" | grep -v "^$" | wc -l)', 'line_count=$(echo "$no_ascii_out" | grep -c -v "^$" || true)')

with open('tests/test.sh', 'w') as f:
    f.write(content)
