with open('plugins/git.sh', 'r') as f:
    content = f.read()

content = content.replace('@{u}', '@\{u\}')

with open('plugins/git.sh', 'w') as f:
    f.write(content)
