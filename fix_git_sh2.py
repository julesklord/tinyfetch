with open('plugins/git.sh', 'r') as f:
    content = f.read()

lines = content.split('\n')
print(f"Line 48: {lines[47]}")
print(f"Line 49: {lines[48]}")
print(f"Line 50: {lines[49]}")
print(f"Line 51: {lines[50]}")
