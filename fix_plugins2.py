with open('plugins/extended/sys_dashboard.sh', 'r') as f:
    content = f.read()

content = content.replace('GREEN="${ESC}[01;32m"\n', '')

with open('plugins/extended/sys_dashboard.sh', 'w') as f:
    f.write(content)


with open('plugins/extended/git_graph.sh', 'r') as f:
    content = f.read()

content = content.replace('GREEN="${ESC}[01;32m"\n', '')
content = content.replace('YELLOW="${ESC}[01;33m"\n', '')

with open('plugins/extended/git_graph.sh', 'w') as f:
    f.write(content)


with open('plugins/git.sh', 'r') as f:
    content = f.read()

content = content.replace('CYAN="${ESC}[01;36m"\n', '')
# Fix literal braces {
content = content.replace('output="${output}${GREEN}⇡${ahead}${RESTORE}"', 'output="${output}${GREEN}⇡${ahead}${RESTORE}"')
content = content.replace('output="${output}${RED}⇣${behind}${RESTORE}"', 'output="${output}${RED}⇣${behind}${RESTORE}"')
content = content.replace('details="${GREEN}●${staged}${RESTORE}"', 'details="${GREEN}●${staged}${RESTORE}"')

with open('plugins/git.sh', 'w') as f:
    f.write(content)


with open('plugins/ip.sh', 'r') as f:
    content = f.read()

content = content.replace('tr -d "\n" | tr -d "\r"', 'tr -d "\n\r"')

with open('plugins/ip.sh', 'w') as f:
    f.write(content)

with open('tests/test.sh', 'r') as f:
    content = f.read()

content = content.replace('''assert_eq() {
  local val1="$1"
  local val2="$2"
  local name="$3"
  if [ "$val1" -eq "$val2" ]; then
    echo "  [PASS] $name"
  else
    echo "  [FAIL] $name (Expected $val2, got $val1)"
    failed=1
  fi
}''', '')

with open('tests/test.sh', 'w') as f:
    f.write(content)
