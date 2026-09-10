Output assets/demo.gif

Set Shell bash
Set FontSize 18
Set Width 900
Set Height 800
Set Padding 20
Set TypingSpeed 100ms
Set Theme "TokyoNight"

# Do setup invisibly
Hide
Type "./termsweeper --config Themes/tokyo-night.yaml"
Enter
Sleep 500ms
Show

# Start capturing
Sleep 1s

# Shrink the board to 12x12
Type "j"
Sleep 500ms
Enter
Sleep 500ms
Type@200ms "hhhhkkkk"
Sleep 500ms
Enter
Sleep 500ms

# Down to 12 mines
Type "j"
Sleep 300ms
Type@40ms "hhhhhhhhhhhhhhhhhhhhhhhhhhhh"
Sleep 500ms

# Seed it, so the game is the same every run
Type "j"
Sleep 300ms
Enter
Type@150ms "demo"
Sleep 300ms
Enter
Sleep 500ms

# Start the game
Type "kkk"
Sleep 300ms
Enter
Sleep 800ms

# Play a while
Type "6j"
Type "llllll"
Space
Sleep 400ms
Type "kkk"
Type "f"
Sleep 400ms
Type "h"
Type "f"
Sleep 400ms
Type "k"
Type "f"
Sleep 400ms
Type "2jll"
Type "f"
Sleep 400ms
Right
Type "f"
Sleep 400ms
Type "jll"
Type "f"
Sleep 400ms
Type "3j"
Type "f"
Sleep 400ms
Down
Right
Type "f"
Sleep 400ms
Type "4k7h"
Type "f"
Sleep 400ms
Type "jhhhh"
Type "f"
Sleep 400ms

# Misread a 2 and step on a mine
Sleep 600ms
Type "jj"
Sleep 600ms
Space

# Linger on the loss
Sleep 4s
