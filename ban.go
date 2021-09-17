package main

import (
	"encoding/binary"
	"fmt"
	"github.com/Postcord/rest"
	"net/http"
	"os"
	"unsafe"
)

func ban() {
	if logFile == "" {
		panic("missing logfile")
	}

	f, err := os.OpenFile(logFile+".dat", os.O_RDONLY, 0777)
	if err != nil {
		if os.IsNotExist(err) {
			panic("logfile not found")
		}

		panic("An unknown error occurred while attempting to open logs: " + err.Error())
	}

	defer func() {
		_ = f.Close()
	}()

	sizeBuf := make([]byte, 4)
	err = binary.Read(f, binary.LittleEndian, sizeBuf)
	if err != nil {
		panic("An error occurred while loading the bans from logs: " + err.Error())
	}

	caseCount := *(*uint32)(unsafe.Pointer(&sizeBuf[0]))
	users := make([]user, caseCount)

	uPtr := 0
	for ; uPtr < len(users); uPtr++ {
		if !users[uPtr].fromBytes(f) {
			panic("an error occurred attempting to load bans from logs")
		}
	}

	rebans := uint32(0)
	for _, u := range users {
		err := client.CreateBan(guild, u.id, &rest.CreateGuildBanParams{Reason: u.reason})
		if err != nil {
			if e, ok := err.(*rest.ErrorREST); ok {
				if e.Status == http.StatusForbidden {
					fmt.Printf("Error: Missing permissions (failed to ban %d): %s\n", u.id, err.Error())
				}
			}
			fmt.Printf("An error occurred re-banning %d: %s\n", u.id, err.Error())
		} else {
			rebans++
		}
	}

	fmt.Printf("Successfully rebanned %d / %d accounts (%d%%)", rebans, caseCount, (100*rebans)/caseCount)
}
