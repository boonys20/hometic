package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"active"}`))
}

func main() {
	fmt.Println("Hello hometic : I'm Gopher !!!")

	r := mux.NewRouter()
	r.HandleFunc("/pair-device", PairDeviceHandler).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Printf("server port : %v", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println(" Starting ..... ")
	log.Fatal(server.ListenAndServe())
}
