# Screenshots of a game in progress, one per theme, for the README

Set Shell bash
Set FontSize 18
Set Width 900
Set Height 900
Set Padding 20
Set TypingSpeed 1ms
Set Theme "TokyoNight"

Hide
Type "./termsweeper --no-config"
Enter
Sleep 500ms

# Grow the board to 30x24, with 99 mines
Type "j"
Enter
Type "lllllllllllllljjjjjjjj"
Enter
Type "jlllllllllllllllllllllllllllllllllllllllllllllllllllllllllll"

# Seed it, so the board is the same every run
Type "j"
Enter
Type "demo"
Enter

# Start the game and play a good part of it
Type "kkk"
Enter
Type "12j15l 3khf2hf4hf4jf2lflf4jfhfk5lf5k3lfkf15hf5jflflfl k2l l k 2j 3j2l j k4l k kl l 7kl k2h h h jh 6h h 4j4h j j 3kf3k3lfkl l l 5l k l l j2lf8jlfjh h 2jhf5h k 2khfkh hfj h h 3l jl j j j4l l l 2klfk2l kl k k 6k k3hfj4hfk7h 2jh 2j2h jh j j jf4kfklfklf6klfk2lf5lf3jhf3j hf2l kl 5l k 2k4lf5jlfj4lf2klf6jhf6h kfhf2jl 4j4h j4hf3k8h 16k4hfk 2j4l l j6l j j jh j l k5l klf3j5l jl 3l k 2kl 9j h 2l 6hf5j12hfjhf4hfkfjhfhf2k3hf4k 2k k k "
Type "j15l"
Sleep 1s
Show
Sleep 500ms
Screenshot assets/default_theme.png
Sleep 500ms
Hide

# Back to the shell
Type "qq"
Sleep 500ms
Type "clear"
Enter

Type "./termsweeper --config Themes/tokyo-night.yaml"
Enter
Sleep 500ms

# Grow the board to 30x24, with 99 mines
Type "j"
Enter
Type "lllllllllllllljjjjjjjj"
Enter
Type "jlllllllllllllllllllllllllllllllllllllllllllllllllllllllllll"

# Seed it, so the board is the same every run
Type "j"
Enter
Type "demo"
Enter

# Start the game and play a good part of it
Type "kkk"
Enter
Type "12j15l 3khf2hf4hf4jf2lflf4jfhfk5lf5k3lfkf15hf5jflflfl k2l l k 2j 3j2l j k4l k kl l 7kl k2h h h jh 6h h 4j4h j j 3kf3k3lfkl l l 5l k l l j2lf8jlfjh h 2jhf5h k 2khfkh hfj h h 3l jl j j j4l l l 2klfk2l kl k k 6k k3hfj4hfk7h 2jh 2j2h jh j j jf4kfklfklf6klfk2lf5lf3jhf3j hf2l kl 5l k 2k4lf5jlfj4lf2klf6jhf6h kfhf2jl 4j4h j4hf3k8h 16k4hfk 2j4l l j6l j j jh j l k5l klf3j5l jl 3l k 2kl 9j h 2l 6hf5j12hfjhf4hfkfjhfhf2k3hf4k 2k k k "
Type "j15l"
Sleep 1s
Show
Sleep 500ms
Screenshot assets/tokyo_night_theme.png
Sleep 500ms
