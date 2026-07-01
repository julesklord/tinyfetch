package main

import "strings"

type Theme struct {
	Name           string
	Primary        string
	Secondary      string
	Success        string
	Warning        string
	Error          string
	Muted          string
	TreeLines      string
	BannerGradient [2][3]int
	BarColors      [3]string // low, mid, high
}

var themes = map[string]Theme{
	"default": {
		Name:      "default",
		Primary:   "\033[38;2;0;242;254m",    // Electric Cyan
		Secondary: "\033[38;2;79;172;254m",   // Blue
		Success:   "\033[01;32m",              // Green
		Warning:   "\033[01;33m",              // Yellow
		Error:     "\033[01;31m",              // Red
		Muted:     "\033[90m",                 // Gray
		TreeLines: "\033[90m",                 // Gray
		BannerGradient: [2][3]int{
			{255, 94, 98}, // Coral Red
			{0, 242, 254}, // Electric Cyan
		},
		BarColors: [3]string{
			"\033[01;32m", // Green
			"\033[01;33m", // Yellow
			"\033[01;31m", // Red
		},
	},
	"catppuccin": {
		Name:      "catppuccin",
		Primary:   "\033[38;2;245;169;127m", // Peach
		Secondary: "\033[38;2;180;190;254m", // Lavender
		Success:   "\033[38;2;166;227;161m", // Green
		Warning:   "\033[38;2;249;226;175m", // Yellow
		Error:     "\033[38;2;243;139;168m", // Red
		Muted:     "\033[38;2;108;112;134m", // Overlay0
		TreeLines: "\033[38;2;69;71;90m",    // Surface1
		BannerGradient: [2][3]int{
			{245, 169, 127}, // Peach
			{180, 190, 254}, // Lavender
		},
		BarColors: [3]string{
			"\033[38;2;166;227;161m", // Green
			"\033[38;2;249;226;175m", // Yellow
			"\033[38;2;243;139;168m", // Red
		},
	},
	"catppuccin-mocha": {
		Name:      "catppuccin-mocha",
		Primary:   "\033[38;2;245;169;127m", // Peach
		Secondary: "\033[38;2;180;190;254m", // Lavender
		Success:   "\033[38;2;166;227;161m", // Green
		Warning:   "\033[38;2;249;226;175m", // Yellow
		Error:     "\033[38;2;243;139;168m", // Red
		Muted:     "\033[38;2;108;112;134m", // Overlay0
		TreeLines: "\033[38;2;69;71;90m",    // Surface1
		BannerGradient: [2][3]int{
			{245, 169, 127}, // Peach
			{180, 190, 254}, // Lavender
		},
		BarColors: [3]string{
			"\033[38;2;166;227;161m",
			"\033[38;2;249;226;175m",
			"\033[38;2;243;139;168m",
		},
	},
	"catppuccin-latte": {
		Name:      "catppuccin-latte",
		Primary:   "\033[38;2;210;106;64m",  // Peach
		Secondary: "\033[38;2;114;123;224m", // Lavender
		Success:   "\033[38;2;64;160;61m",   // Green
		Warning:   "\033[38;2;202;149;41m",  // Yellow
		Error:     "\033[38;2;210;50;97m",   // Red
		Muted:     "\033[38;2;114;118;139m", // Overlay0
		TreeLines: "\033[38;2;180;184;198m", // Surface1
		BannerGradient: [2][3]int{
			{210, 106, 64},  // Peach
			{114, 123, 224}, // Lavender
		},
		BarColors: [3]string{
			"\033[38;2;64;160;61m",
			"\033[38;2;202;149;41m",
			"\033[38;2;210;50;97m",
		},
	},
	"dracula": {
		Name:      "dracula",
		Primary:   "\033[38;2;189;147;249m", // Purple
		Secondary: "\033[38;2;139;233;253m", // Cyan
		Success:   "\033[38;2;80;250;123m",  // Green
		Warning:   "\033[38;2;255;184;108m", // Orange
		Error:     "\033[38;2;255;85;85m",   // Red
		Muted:     "\033[38;2;98;114;164m",  // Comment
		TreeLines: "\033[38;2;68;71;90m",    // Background
		BannerGradient: [2][3]int{
			{189, 147, 249}, // Purple
			{139, 233, 253}, // Cyan
		},
		BarColors: [3]string{
			"\033[38;2;80;250;123m",
			"\033[38;2;255;184;108m",
			"\033[38;2;255;85;85m",
		},
	},
	"nord": {
		Name:      "nord",
		Primary:   "\033[38;2;136;192;208m", // Frost 2
		Secondary: "\033[38;2;129;161;193m", // Frost 1
		Success:   "\033[38;2;163;190;140m", // Green
		Warning:   "\033[38;2;235;203;139m", // Yellow
		Error:     "\033[38;2;191;97;106m",  // Red
		Muted:     "\033[38;2;76;86;106m",   // Polar Night 3
		TreeLines: "\033[38;2;59;66;82m",    // Polar Night 2
		BannerGradient: [2][3]int{
			{136, 192, 208}, // Frost 2
			{129, 161, 193}, // Frost 1
		},
		BarColors: [3]string{
			"\033[38;2;163;190;140m",
			"\033[38;2;235;203;139m",
			"\033[38;2;191;97;106m",
		},
	},
	"tokyonight": {
		Name:      "tokyonight",
		Primary:   "\033[38;2;125;207;255m", // Blue
		Secondary: "\033[38;2;187;154;247m", // Purple
		Success:   "\033[38;2;158;206;106m", // Green
		Warning:   "\033[38;2;224;175;104m", // Yellow
		Error:     "\033[38;2;247;118;142m", // Red
		Muted:     "\033[38;2;86;95;137m",   // Comment
		TreeLines: "\033[38;2;41;48;80m",    // Dark bg
		BannerGradient: [2][3]int{
			{125, 207, 255}, // Blue
			{187, 154, 247}, // Purple
		},
		BarColors: [3]string{
			"\033[38;2;158;206;106m",
			"\033[38;2;224;175;104m",
			"\033[38;2;247;118;142m",
		},
	},
	"gruvbox": {
		Name:      "gruvbox",
		Primary:   "\033[38;2;251;174;88m",  // Yellow
		Secondary: "\033[38;2;131;165;152m", // Aqua
		Success:   "\033[38;2;184;187;38m",  // Green
		Warning:   "\033[38;2;250;189;47m",  // Yellow
		Error:     "\033[38;2;251;73;52m",   // Red
		Muted:     "\033[38;2;146;131;116m", // Gray
		TreeLines: "\033[38;2;92;102;110m",  // Dark gray
		BannerGradient: [2][3]int{
			{251, 174, 88},  // Yellow
			{131, 165, 152}, // Aqua
		},
		BarColors: [3]string{
			"\033[38;2;184;187;38m",
			"\033[38;2;250;189;47m",
			"\033[38;2;251;73;52m",
		},
	},
	"everforest": {
		Name:      "everforest",
		Primary:   "\033[38;2;164;190;134m", // Green
		Secondary: "\033[38;2;130;187;183m", // Blue
		Success:   "\033[38;2;164;190;134m", // Green
		Warning:   "\033[38;2;223;175;100m", // Yellow
		Error:     "\033[38;2;229;115;115m", // Red
		Muted:     "\033[38;2;127;136;123m", // Gray
		TreeLines: "\033[38;2;51;61;59m",    // Dark bg
		BannerGradient: [2][3]int{
			{164, 190, 134}, // Green
			{130, 187, 183}, // Blue
		},
		BarColors: [3]string{
			"\033[38;2;164;190;134m",
			"\033[38;2;223;175;100m",
			"\033[38;2;229;115;115m",
		},
	},
	"monokai": {
		Name:      "monokai",
		Primary:   "\033[38;2;166;226;46m",  // Green
		Secondary: "\033[38;2;174;129;255m", // Purple
		Success:   "\033[38;2;166;226;46m",  // Green
		Warning:   "\033[38;2;253;151;31m",  // Orange
		Error:     "\033[38;2;249;38;114m",  // Pink
		Muted:     "\033[38;2;117;113;94m",  // Comment
		TreeLines: "\033[38;2;63;62;54m",    // Background
		BannerGradient: [2][3]int{
			{166, 226, 46},  // Green
			{174, 129, 255}, // Purple
		},
		BarColors: [3]string{
			"\033[38;2;166;226;46m",
			"\033[38;2;253;151;31m",
			"\033[38;2;249;38;114m",
		},
	},
	"rose-pine": {
		Name:      "rose-pine",
		Primary:   "\033[38;2;235;188;210m", // Rose
		Secondary: "\033[38;2;174;192;219m", // Foam
		Success:   "\033[38;2;59;178;152m",  // Pine
		Warning:   "\033[38;2;248;196;113m", // Gold
		Error:     "\033[38;2;224;122;122m", // Love
		Muted:     "\033[38;2;115;115;146m", // Muted
		TreeLines: "\033[38;2;39;35;52m",    // Base
		BannerGradient: [2][3]int{
			{235, 188, 210}, // Rose
			{174, 192, 219}, // Foam
		},
		BarColors: [3]string{
			"\033[38;2;59;178;152m",
			"\033[38;2;248;196;113m",
			"\033[38;2;224;122;122m",
		},
	},
	"solarized": {
		Name:      "solarized",
		Primary:   "\033[38;2;38;139;210m",  // Blue
		Secondary: "\033[38;2;211;54;130m",  // Magenta
		Success:   "\033[38;2;133;153;0m",   // Green
		Warning:   "\033[38;2;181;137;0m",   // Yellow
		Error:     "\033[38;2;220;50;47m",   // Red
		Muted:     "\033[38;2;101;123;131m", // Base0
		TreeLines: "\033[38;2;88;110;117m",  // Base01
		BannerGradient: [2][3]int{
			{38, 139, 210}, // Blue
			{211, 54, 130}, // Magenta
		},
		BarColors: [3]string{
			"\033[38;2;133;153;0m",
			"\033[38;2;181;137;0m",
			"\033[38;2;220;50;47m",
		},
	},
}

var currentTheme = themes["default"]

func GetTheme() Theme {
	return currentTheme
}

func SetTheme(name string) bool {
	name = strings.ToLower(name)
	if t, ok := themes[name]; ok {
		currentTheme = t
		return true
	}
	return false
}

func ListThemes() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	return names
}