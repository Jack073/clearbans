package main

import (
	"flag"
	"fmt"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
	"net/http"
	"os"
	"strings"
)

var (
	token       = os.Getenv("TOKEN")
	reason      string
	guild       objects.Snowflake
	deletedOnly bool
	logFile     string

	client *rest.Client
)

func init() {
	flag.Uint64Var((*uint64)(&guild), "guild", 0, "The guild to manage bans for")
	flag.BoolVar(&deletedOnly, "deleted", false, "Whether or not to only delete bans for deleted account")
	flag.StringVar(&reason, "reason", "[clearbans]: no reason provided", "The reason to include in the audit log")
	flag.StringVar(&logFile, "logfile", "", "The file to write the logs to (optional)")

	flag.Parse()

	if token == "" {
		panic("Empty TOKEN env var")
	}

	if guild == 0 {
		panic("guild cannot be empty")
	}
}

var deletedUser = []byte("Deleted User ")

func isDeletedUser(user *objects.User) bool {
	name := []byte(user.Username)

	if len(name) < len(deletedUser) {
		return false
	}

	for a, b := range deletedUser {
		if name[a] != b {
			return false
		}
	}

	for i := len(deletedUser); i < len(name); i++ {
		switch {
		case 48 <= name[i] && name[i] <= 57:
		case 65 <= name[i] && name[i] <= 70:
		case 97 <= name[i] && name[i] <= 102:
			break

		default:
			return false
		}
	}

	return true
}

type user struct {
	id     objects.Snowflake
	name   string
	reason string
}

func main() {
	client = rest.New(&rest.Config{
		Authorization: "Bot " + token,
		UserAgent:     "DiscordBot (https://github.com/Jack073/clearbans, 1.0)",
		Ratelimiter: rest.NewMemoryRatelimiter(&rest.MemoryConf{
			MaxRetries: 5,
		}),
	})

	bans, err := client.GetGuildBans(guild)
	if err != nil {
		panic(fmt.Errorf("error when attempting to fetch guild bans: %w", err))
	}

	fmt.Println("Loaded", len(bans), "bans from guild")

	var banIDs []user

	if logFile != "" {
		banIDs = make([]user, 0, len(bans))
		defer func() {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR, 0777)
			if err != nil {
				panic("failed to write to log file: " + err.Error())
			}

			txtBuilder := &strings.Builder{}
			for _, ban := range banIDs {
				txtBuilder.WriteString(fmt.Sprintf("Unbanned user %d (%s) [ban reason: %s]\n", ban.id, ban.name, ban.reason))
			}

			_, _ = file.WriteString(txtBuilder.String())

			err = file.Close()
			if err != nil {
				panic("failed to write close log file: " + err.Error())
			}
		}()
	}

	for i, ban := range bans {
		if deletedOnly && !isDeletedUser(ban.User) {
			continue
		}

		err := client.RemoveGuildBan(guild, ban.User.ID, reason)

		if err != nil {
			if e, ok := err.(*rest.ErrorREST); ok {
				if e.Status == http.StatusForbidden {
					fmt.Println("Error: Missing permissions, exiting\n", err.Error())
					return
				}
			}
			fmt.Printf("An error occurred unbanning %s (%d): %s\n", ban.User.Username, ban.User.ID, err.Error())
		} else {
			tag := fmt.Sprintf("%s#%s", ban.User.Username, ban.User.Discriminator)
			if logFile != "" {
				banIDs = append(banIDs, user{
					id:     ban.User.ID,
					name:   tag,
					reason: strings.ReplaceAll(ban.Reason, "\n", "\\n"),
				})
			}
			fmt.Printf("%d) Successfully unbanned: %s (%d)\n", i+1, tag, ban.User.ID)
		}
	}
}
