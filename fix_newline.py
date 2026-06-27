with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

content = content.replace('fmt.Printf("Usage: %s [--no-ascii] [--minimal] [--noframe] [--output=json|xml|txt]\n", os.Args[0])', 'fmt.Printf("Usage: %s [--no-ascii] [--minimal] [--noframe] [--output=json|xml|txt]\\n", os.Args[0])')

with open('cmd/tinyfetch/main.go', 'w') as f:
    f.write(content)
