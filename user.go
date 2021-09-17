package main

import (
	"encoding/binary"
	"fmt"
	"github.com/Postcord/objects"
	"io"
	"unsafe"
)

type user struct {
	id     objects.Snowflake
	name   string
	reason string
}

func (u *user) String() string {
	return fmt.Sprintf("Unbanned user %d (%s) [ban reason: %s]\n", u.id, u.name, u.reason)
}

func (u *user) toBytes() []byte {
	id := *(*[8]byte)(unsafe.Pointer(&u.id))
	reasonLength := int16(len(u.reason))
	reasonLengthBytes := *(*[2]byte)(unsafe.Pointer(&(reasonLength)))

	return append(
		id[:],
		append(reasonLengthBytes[:], []byte(u.reason)...)...,
	)
}

func (u *user) fromBytes(reader io.Reader) bool {
	input := make([]byte, 8)
	err := binary.Read(reader, binary.LittleEndian, input)
	if err != nil {
		return false
	}

	u.id = *((*objects.Snowflake)(unsafe.Pointer(&input[0])))

	input = make([]byte, 2)
	err = binary.Read(reader, binary.LittleEndian, input)
	if err != nil {
		return false
	}

	reasonLength := *((*uint16)(unsafe.Pointer(&input[0])))

	input = make([]byte, reasonLength)
	err = binary.Read(reader, binary.LittleEndian, input)
	if err != nil {
		return false
	}
	u.reason = string(input)
	return true
}
