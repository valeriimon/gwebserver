package appmanager

import (
	"log"
	"net/http"

	"github.com/go-webserver/main-service"

	"github.com/gorilla/mux"
)

// ExecuteServer : execute app-manager server
func ExecuteServer() {
	err := mainservice.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()

	r.Methods("POST").Path("/authenticate").HandlerFunc(authHandler)
	r.Methods("POST").Path("/create-user").HandlerFunc(createUserHandler)
	r.Methods("GET").Path("/data").HandlerFunc(getData)
	log.Fatal(http.ListenAndServe(":8888", r))

}
