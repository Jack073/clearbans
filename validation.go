package main

import "github.com/Postcord/objects"

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
