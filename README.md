**SETUP**

Build the bot's binary file with:
```
$ go build bot.go
```

In the same location that the bot's binary executable is,
ensure that you create a plaintext file for the bot's
database called database.txt. Example:
```
$ cd BOT_BINARY_DIRECTORY
$ touch database.txt
```

To launch the bot:
```
$ ./bot -t YOUR_BOT_TOKEN
```

**BOT USAGE**

To award rep to a user just begin the message with !rep and
mention every user that you want to give rep to. Example:
```
!rep @username1 @username2 @username3
```

**NOTES**

This bot is currently under development so documentation is
highly likely to be out-of-date. For definitive information
on behaviour of bot, refer to the source code.