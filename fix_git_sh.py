import re
with open('plugins/git.sh', 'r') as f:
    content = f.read()

content = re.sub(r'output="\$\{output\}\$\{GREEN\}⇡\$\{ahead\}\$\{RESTORE\}"', r'output="${output}${GREEN}\u21E1${ahead}${RESTORE}"', content)
content = re.sub(r'output="\$\{output\}\$\{RED\}⇣\$\{behind\}\$\{RESTORE\}"', r'output="${output}${RED}\u21E3${behind}${RESTORE}"', content)
# Check what's going on on line 49 and 50
lines = content.split('\n')
print(f"Line 49: {lines[48]}")
print(f"Line 50: {lines[49]}")

# Fix literal braces warning
content = content.replace('output="${output}${GREEN}⇡${ahead}${RESTORE}"', 'output="${output}${GREEN}\\⇡${ahead}${RESTORE}"')

with open('plugins/git.sh', 'w') as f:
    f.write(content)
