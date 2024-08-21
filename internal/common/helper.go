package common

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"
)

// The SplitUserIDAndPasswordFromAuth takes a Basic Authorization token value, separates the UserID and the password and returns them
func SplitUserIDAndPasswordFromAuth(auth string) (userID string, pass string, err error) {
	authSplit := strings.Split(auth, " ")
	if len(authSplit) != 2 {
		err = errors.New("incorrect value for authorization header")
		return
	}
	authBase64 := authSplit[1]
	decodedByte, err := base64.StdEncoding.DecodeString(authBase64)
	if err != nil {
		return
	}
	decodedString := string(decodedByte)
	splitString := strings.Split(decodedString, ":")
	if len(splitString) != 2 {
		err = errors.New("incorrect value for authorization header")
		return
	}
	userID = splitString[0]
	pass = splitString[1]
	log.Printf("%s %s %s %s", authBase64, decodedString, userID, pass)
	return
}
