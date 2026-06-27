import re

with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

# Replace the current main with the simplified one
main_func = """func main() {
	noASCII, minimal, noFrame, outputFmt := parseFlags()
	infoObj := gatherInfo()
	renderOutput(noASCII, minimal, noFrame, outputFmt, infoObj)
}
"""

# The current main is at the bottom of the file
# We will just replace everything from "func main() {" to the end of the file
content = re.sub(r'func main\(\) \{.*', main_func, content, flags=re.DOTALL)

with open('cmd/tinyfetch/main.go', 'w') as f:
    f.write(content)
