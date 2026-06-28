with open('tests/test.sh', 'r') as f:
    content = f.read()
content = content.replace('assert_eq() {', '')
content = content.replace('  local val1="$1"', '')
content = content.replace('  local val2="$2"', '')
content = content.replace('  local name="$3"', '')
content = content.replace('  if [ "$val1" -eq "$val2" ]; then', '')
content = content.replace('    echo "  [PASS] $name"', '')
content = content.replace('  else', '')
content = content.replace('    echo "  [FAIL] $name (Expected $val2, got $val1)"', '')
content = content.replace('    failed=1', '')
content = content.replace('  fi', '')
content = content.replace('}', '}')
content = content.replace('line_count=$(echo "$no_ascii_out" | grep -v "^$" | wc -l)', 'line_count=$(echo "$no_ascii_out" | grep -c -v "^$")')
with open('tests/test.sh', 'w') as f:
    f.write(content)

with open('plugins/extended/git_graph.sh', 'r') as f:
    content = f.read()
content = content.replace('GREEN="${ESC}[01;32m"\n', '')
content = content.replace('YELLOW="${ESC}[01;33m"\n', '')
with open('plugins/extended/git_graph.sh', 'w') as f:
    f.write(content)

with open('plugins/extended/sys_dashboard.sh', 'r') as f:
    content = f.read()
content = content.replace('GREEN="${ESC}[01;32m"\n', '')
with open('plugins/extended/sys_dashboard.sh', 'w') as f:
    f.write(content)

with open('plugins/ip.sh', 'r') as f:
    content = f.read()
content = content.replace('tr -d "\n" | tr -d "\r"', 'tr -d "\n\r"')
with open('plugins/ip.sh', 'w') as f:
    f.write(content)

with open('plugins/git.sh', 'r') as f:
    content = f.read()
content = content.replace('CYAN="${ESC}[01;36m"\n', '')
content = content.replace('${ESC}[01;35m ${branch}${RESTORE}', '${ESC}[01;35m ${branch}${RESTORE}')
with open('plugins/git.sh', 'w') as f:
    f.write(content)
