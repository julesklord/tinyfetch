# arbol — demo.tape
# Records an animated GIF showcasing arbol's main features
Output docs/demo.gif

# === Visual Settings ===
Set FontSize 13
Set FontFamily "CommitMono Nerd Font Mono"
Set Width 1100
Set Height 760
Set Theme "Catppuccin Mocha"
Set Padding 14
Set TypingSpeed 60ms
Set PlaybackSpeed 1.0
Set CursorBlink false

# === Launch shell silently, then run arbol ===
Hide
Type "export TERM=xterm-256color"
Enter
Sleep 200ms
Show

# === Step 1: Run arbol — full tree output ===
Type "./arbol"
Sleep 200ms
Enter
Sleep 5s

# === Step 2: Blank line for visual breathing room ===
Sleep 1.5s

# === Step 3: Show minimal mode ===
Type "./arbol --no-ascii"
Sleep 200ms
Enter
Sleep 4s

Sleep 1.5s

# === Step 4: Structured JSON export ===
Type "./arbol --output=json | head -20"
Sleep 200ms
Enter
Sleep 3.5s

Sleep 1.5s

# === Step 5: Show help ===
Type "./arbol --help"
Sleep 200ms
Enter
Sleep 2.5s

Sleep 1.5s

# === Step 6: Exit cleanly ===
Type "exit"
Enter
Sleep 500ms
