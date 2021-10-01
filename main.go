package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Account API begin ...")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexProcess)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func indexProcess(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Welcome to Account API")
	fmt.Fprintf(writer, "\nby Software Engineer - Alejandro Rizzuto")
}
