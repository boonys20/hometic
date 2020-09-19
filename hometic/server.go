package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

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

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	row := db.QueryRow("INSERT INTO pairs (DEVICE_ID, USER_ID) values ($1, $2) RETURNING id", p.DeviceID, p.UserID)
	var id int
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("can't scan id", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	fmt.Printf("pair: %#v\n", p)
	resp, err := json.Marshal(p)

	fmt.Println("success id : ", id)
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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}
	defer db.Close()

	log.Println(" Starting ..... ")
	log.Fatal(server.ListenAndServe())
}
