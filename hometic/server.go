package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Pair struct {
	DeviceID int64
	UserID   int64
}

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var p Pair
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	defer r.Body.Close()
	fmt.Printf("pair: %#v\n", p)
	resp, err := json.Marshal(p)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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
