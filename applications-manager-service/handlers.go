package appmanager

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-webserver/main-service"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userInst mainservice.User

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&userInst); err != nil {
		log.Fatal(err.Error())
	}

	err := userInst.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(userInst)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	email, pass, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Basic Auth Failed", http.StatusForbidden)
		return
	}

	user := &mainservice.User{
		Email:    email,
		Password: pass,
	}

	user, err := user.Authenticate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(&user)
}

// getData : temp method for debug bucket data
func getData(w http.ResponseWriter, r *http.Request) {
	db := mainservice.Database{}

	res, err := db.GetData("users", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}
