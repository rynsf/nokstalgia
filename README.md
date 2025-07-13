# Nokstalgia
Nokstalgia is an emulator for retro Nokia games.

![space impact gameplay](https://github.com/user-attachments/assets/fceae164-46d8-42c6-b886-98c65531d7e1)

![snake gameplay](https://github.com/user-attachments/assets/a3379aa8-3880-48c2-ac37-286ea1268d5f)

## Getting Started
Current version of nokstalgia works via a CLI. Here's how to use it:
```bash
# list games available in a ROM
$ ./nokstalgia list ./2100.fls
Games available in ROM are:
1. Snake II
2. Space Impact
3. Link5

# run a game with the run command
$ ./nokstalgia run "Space Impact" 2100.fls

# games can also be selected with any uniquely identifiable substring of its name
$ ./nokstalgia run space 2100.fls

# or run snake with
$ ./nokstalgia run snake 2100.fls

# specify different game settings like this
$ ./nokstalgia run snake -level 4 -maze 2 2100.fls

# list available game settings with -help
$ ./nokstalgia run snake -help 2100.fls

```

## Current Status
all games from 2100 5.49 work.

## TODO
- [x] implement a CLI interface
- [ ] emulate audio
- [ ] emulate the menuing system
- [ ] make NokiX games work
- [ ] support games from other phones


## Credits
[NokiX](https://arkadiusz.wahlig.eu/NokiX.html) \
yak \
g3gg0 \
nok5rev
