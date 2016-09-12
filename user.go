package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type User struct {
	Id       int
	Username string
	Password string
}

func (user *User) Save(db DBConn) error {
	return db.SaveUser(user)
}

func MarshalUser(u *User) []byte {
	e, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return e
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}