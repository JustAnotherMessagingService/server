package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const version = "0.0.1"

var db DBConn

func main() {
	fmt.Println("JAMA Server version ", version)
	fmt.Println("[*] Starting server...")

	StartServer()
}

func StartServer() {
	fmt.Println("[+] Loading BoltDB")
	boltdb, err := BoltDBOpen("my.db")
	if err != nil {
		fmt.Printf("[!] Error opening BoltDB: %s\n", err.Error())
		os.Exit(1)
	}
	db = boltdb
	fmt.Println("[+] BoltDB loaded")

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)            // Return index.html
	r.HandleFunc("/favicon.ico", iconHandler) // Return favicon

	r.HandleFunc("/api/", apiHandler) // Return API reference
	r.HandleFunc("/api/user/{id}", apiUserGetHandler).Methods("GET")

	r.HandleFunc("/api/user", apiUserHandler)       // Handle all user API requests
	r.HandleFunc("/api/message", apiMessageHandler) // Handle all message API requests
	r.HandleFunc("/api/auth", apiAuthHandler)       // Handle all auth API requests

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
	return
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Serving index.html")
	filename := "index.html"
	http.ServeFile(w, r, filename)
}

func iconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Serving favicon.ico")
	filename := "favicon.ico"
	http.ServeFile(w, r, filename)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Serving api.html")
	filename := "api.html"
	http.ServeFile(w, r, filename)
}
