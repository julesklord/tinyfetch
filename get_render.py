import re

with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

# Find the start of the render logic in main
# We want to grab everything from "// Intercept output format flag early" to the end of main
match = re.search(r'(\t// Intercept output format flag early.*?)\n}\n', content, re.DOTALL)
if match:
    render_logic = match.group(1)

    # We need to replace variables with infoObj.* equivalents
    replacements = {
        'hostname': 'infoObj.Hostname',
        'osName': 'infoObj.OSName',
        'kernel': 'infoObj.Kernel',
        'uptimeVal': 'infoObj.UptimeVal',
        'shellVal': 'infoObj.ShellVal',
        'cpuVal': 'infoObj.CPUVal',
        'memRaw': 'infoObj.MemRaw',
        'diskRaw': 'infoObj.DiskRaw',
        'pluginKeys': 'infoObj.PluginKeys',
        'pluginVals': 'infoObj.PluginVals'
    }

    for old, new in replacements.items():
        # Replace only as whole words
        render_logic = re.sub(r'\b' + old + r'\b', new, render_logic)

    render_func = """
func renderOutput(noASCII, minimal, noFrame bool, outputFmt string, infoObj SystemInfo) {
""" + render_logic + """
}
"""

    # Insert renderOutput function before main function
    content = re.sub(r'func main\(\) \{', render_func + '\nfunc main() {', content)

    with open('cmd/tinyfetch/main.go', 'w') as f:
        f.write(content)
else:
    print("Could not find the target code block in main")
