package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Postcord/rest"
	"net/http"
	"os"
	"strings"
	"unsafe"
)

func unban() {
	bans, err := client.GetGuildBans(guild)
	if err != nil {
		panic(fmt.Errorf("error when attempting to fetch guild bans: %w", err))
	}

	fmt.Println("Loaded", len(bans), "bans from guild")

	var unbannedUsers []user

	if logFile != "" {
		unbannedUsers = make([]user, 0, len(bans))
		defer func() {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR, 0777)
			if err != nil {
				panic("failed to write to log file: " + err.Error())
			}

			txtBuilder := &strings.Builder{}

			binFile, err := os.OpenFile(logFile+".dat", os.O_CREATE|os.O_RDWR, 0777)
			if err != nil {
				panic("failed to write dat ban file" + err.Error())
			}

			binBuf := &bytes.Buffer{}

			for _, user := range unbannedUsers {
				err = binary.Write(binBuf, binary.LittleEndian, user.toBytes())
				if err != nil {
					panic("error occurred writing dat file: " + err.Error())
				}
				txtBuilder.WriteString(user.String())
			}

			_, _ = file.WriteString(txtBuilder.String())

			binData := binBuf.Bytes()
			caseLength := uint32(len(unbannedUsers))
			sizeData := *(*[4]byte)(unsafe.Pointer(&caseLength))
			err = binary.Write(binFile, binary.LittleEndian, sizeData[:])
			if err != nil {
				panic("error writing dat file: " + err.Error())
			}

			_, err = binFile.Write(binData)
			if err != nil {
				panic("error writing dat file: " + err.Error())
			}

			err = binFile.Close()
			if err != nil {
				panic("error closing dat file: " + err.Error())
			}

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
				unbannedUsers = append(unbannedUsers, user{
					id:     ban.User.ID,
					name:   tag,
					reason: strings.ReplaceAll(ban.Reason, "\n", "\\n"),
				})
			}
			fmt.Printf("%d) Successfully unbanned: %s (%d)\n", i+1, tag, ban.User.ID)
		}
	}
}
