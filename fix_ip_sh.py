with open('plugins/ip.sh', 'r') as f:
    content = f.read()

content = content.replace("tr -d 'addr:'", "sed 's/addr://'")

with open('plugins/ip.sh', 'w') as f:
    f.write(content)
