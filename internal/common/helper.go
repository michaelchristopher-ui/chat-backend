package common

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/google/uuid"
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

// The CheckNilFields function takes an interface that has fields and returns a formatted error string listing all the empty fields if there are any
func CheckNilFields(s interface{}) error {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	emptyFields := []string{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		if field.Kind() == reflect.Ptr && field.IsNil() {
			emptyFields = append(emptyFields, fieldName)
		}
	}
	if len(emptyFields) > 0 {
		return fmt.Errorf("[CheckNilFields] These fields are empty: %s", strings.Join(emptyFields[:], ", "))
	}
	return nil
}

// The GenerateUUID function generates a universally unique id. Usually used for generating FlowID.
func GenerateUUID() string {
	return uuid.New().String()
}
