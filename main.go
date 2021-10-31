package main

import (
	"flag"
	"fmt"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	token       = os.Getenv("TOKEN")
	reason      string
	guild       objects.Snowflake
	deletedOnly bool
	logFile     string
	mode        string

	client *rest.Client
)

func init() {
	flag.Uint64Var((*uint64)(&guild), "guild", 0, "The guild to manage bans for")
	flag.Uint64("g", 0, "The guild to manage bans for")

	flag.BoolVar(&deletedOnly, "deleted", false, "Whether or not to only delete bans for deleted account")
	flag.Bool("d", false, "Whether or not to only delete bans for deleted account")

	flag.StringVar(&reason, "reason", "[ClearBans]: no reason provided", "The reason to include in the audit log")
	flag.String("r", "[ClearBans]: no reason provided", "The reason to include in the audit log")

	flag.StringVar(&logFile, "logfile", "", "The file to write the logs to (required for bans, but new logs are only available when unbanning)")
	flag.String("l", "", "The file to write the logs to (required for bans, but new logs are only available when unbanning)")

	flag.StringVar(&mode, "mode", "", "Mode of operation: whether to unban users, or re-ban all users in the log file (`ban` or `unban`")
	flag.String("m", "", "Mode of operation: whether to unban users, or re-ban all users in the log file (`ban` or `unban`")

	flag.Parse()

	if guild == 0 {
		f := flag.Lookup("g")
		if f != nil {
			b, err := strconv.ParseUint(f.Value.String(), 10, 64)
			if err == nil {
				guild = objects.Snowflake(b)
			}
		}
	}

	if !deletedOnly {
		f := flag.Lookup("d")
		if f != nil && f.Value.String() == "true" {
			deletedOnly = true
		}
	}

	if logFile == "" {
		f := flag.Lookup("l")
		if f != nil && f.Value.String() != "" {
			logFile = f.Value.String()
		}
	}

	if mode == "" {
		f := flag.Lookup("m")
		if f != nil && f.Value.String() != "" {
			mode = f.Value.String()
		}
	}

	if reason == "[ClearBans]: no reason provided" {
		f := flag.Lookup("r")
		if f != nil && f.Value.String() != "" {
			reason = f.Value.String()
		}
	}

	if token == "" {
		panic("Empty TOKEN env var")
	}

	if guild == 0 {
		panic("guild cannot be empty")
	}
}

func main() {
	client = rest.New(&rest.Config{
		Authorization: "Bot " + token,
		UserAgent:     "DiscordBot (https://github.com/Jack073/clearbans, 1.0)",
		Ratelimiter: rest.NewMemoryRatelimiter(&rest.MemoryConf{
			MaxRetries: 5,
		}),
	})

	switch strings.ToLower(mode) {
	case "unban", "":
		fmt.Println("In unban mode (will begin in 5 seconds)")
		time.Sleep(5 * time.Second)
		unban()
	case "ban", "reban", "re-ban":
		fmt.Println("In ban mode (will begin in 5 seconds)")
		time.Sleep(5 * time.Second)
		ban()
	default:
		fmt.Println("[ERROR] mode must be \"ban\" or \"unban\"")
	}
}
