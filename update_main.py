import re

with open('cmd/tinyfetch/main.go', 'r') as f:
    content = f.read()

parse_flags_func = """
func parseFlags() (bool, bool, bool, string) {
	noASCII := false
	minimal := false
	noFrame := false
	outputFmt := ""

	for _, arg := range os.Args[1:] {
		if arg == "--no-ascii" {
			noASCII = true
		} else if arg == "--minimal" {
			minimal = true
		} else if arg == "--noframe" {
			noFrame = true
		} else if strings.HasPrefix(arg, "--output=") {
			outputFmt = strings.TrimPrefix(arg, "--output=")
		} else if arg == "--help" || arg == "-h" {
			fmt.Printf("Usage: %s [--no-ascii] [--minimal] [--noframe] [--output=json|xml|txt]\\n", os.Args[0])
			os.Exit(0)
		}
	}
	return noASCII, minimal, noFrame, outputFmt
}

"""

# Insert parseFlags function before main function
content = re.sub(r'func main\(\) \{', parse_flags_func + 'func main() {', content)

with open('cmd/tinyfetch/main.go', 'w') as f:
    f.write(content)
