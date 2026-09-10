Output assets/demo.gif

Require go

Set Shell bash
Set FontSize 18
Set Width 900
Set Height 800
Set Padding 20
Set TypingSpeed 120ms
Set Theme "TokyoNight"

# Build quietly so the recording starts on a clean prompt.
Hide
Type `export PS1="$ "` Enter
Type "clear" Enter
Type "go build -o /tmp/termsweeper ." Enter
Wait@120s /^\$\s*$/
Type "clear" Enter
Show

Type "termsweeper" Sleep 500ms Enter
Sleep 2s

# --- Menu: tune the board before starting -------------------------------
Down Sleep 400ms          # Board Rows
Left@150ms 6              # 16 -> 10
Sleep 700ms
Down Sleep 400ms          # Board Cols
Left@150ms 4              # 16 -> 12
Sleep 700ms
Down Sleep 400ms          # Mines
Left@80ms 25              # 40 -> 15
Sleep 1s
Up@200ms 3                # back to Start Game
Sleep 800ms
Enter
Sleep 1.5s

# --- Playing: vim motions with count prefixes ---------------------------
Type "5j" Sleep 800ms
Type "6l" Sleep 800ms
Space Sleep 1.5s

Type "3k" Sleep 600ms
Type "4h" Sleep 600ms
Space Sleep 1.5s

# Flag a suspicious cell
Type "2j" Sleep 500ms
Type "l" Sleep 500ms
Type "F" Sleep 1.5s

# Poke around some more
Type "3l" Sleep 600ms
Space Sleep 1.5s
Type "2j" Sleep 500ms
Space Sleep 2s

Type "4h" Sleep 500ms
Type "F" Sleep 1.5s

Type "2k" Sleep 500ms
Space Sleep 2.5s

# --- Back to the menu and out ------------------------------------------
Type "Q" Sleep 1.5s
Down@200ms 5
Sleep 700ms
Enter
Sleep 1.5s
