package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"active"}`))
}

func main() {
	fmt.Println("Hello hometic : I'm Gopher !!!")

	r := mux.NewRouter()
	r.HandleFunc("/pair-device", PairDeviceHandler).Methods(http.MethodPost)

	server := http.Server{
		Addr:    "127.0.0.1:2009",
		Handler: r,
	}

	log.Println(" Starting ..... ")
	log.Fatal(server.ListenAndServe())
}
