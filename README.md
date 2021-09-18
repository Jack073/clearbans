# ClearBans

ClearBans is a simple tool developed as a small CLI, which can be used to simply clear the ban list of a Discord guild (and optionally use the log file to restore it).

## Installing ClearBans
1) You'll first need to install [Go](https://golang.org/doc/install) to be able to compile it.
2) `git clone` / download this repository.
3) cd into it the ClearBans folder.
4) Run `go build` - for those of you unfamiliar with Go, if this has no visible output in your terminal, that's a good sign, it just means it was able to generate the executable without error. Also worth noting - it may take longer to compile the first time as it will also need to install any dependency packages (this doesn't require anything else from you, it just means it might take a minute or two longer).

After compilation, you should now have an executable in your current working directory called clearbans (an exe on windows etc), you can simply run this with a `-help` flag appended and it will give you a basic help menu, though I'll explain the different functions below.

```
clearbans
	--help
		Shows a summary of the different flags available.

	--deleted
		Unban suspected deleted accounts only (Deleted User XXXXXXX#1234).

	--guild <id>
		The ID of the guild of the ban list (bot must be in this guild with manage bans permissions).

	--mode [ban|unban]
		The mode of operation - whether to unban accounts, or reban all those in the given logfile. Defaults to unban.

	--logfile [file]
		The name of the file to write logs to (optional for unban).
		Note (unbans only): if set, this will also generate another file "<file>.dat" which contains an easily machine parsable record of all the bans and the associated reasons.
		Note (bans only): this is required for rebanning users, and you must supply the file exactly as you did when you ran it in unban mode - it will automatically add the .dat suffix and choose the correct file.

	--reason [string]
		The reason to use when unbanning accounts.
```

## Authentication
To be able to use this program, you'll also need to have the token of a bot in the guild you want to edit the ban list of (the bot must have the manage bans permission). Then simply add that token as an environment variable called `TOKEN`.