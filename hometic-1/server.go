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

type Pair struct {
	DeviceID int64
	UserID   int64
}

type CreatePairDeviceFunc func(p Pair) error

func PairDeviceHandler(createPairDevice CreatePairDeviceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p Pair
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		defer r.Body.Close()
		fmt.Printf("pair: %#v\n", p)

		err = createPairDevice(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Write([]byte(`{"status":"active"}`))
	}
}

func CreatePairDevice(p Pair) error {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}

	_, err = db.Exec("INSERT INTO pairs (DEVICE_ID, USER_ID) values ($1, $2) RETURNING id", p.DeviceID, p.UserID)

	return err

}

func main() {
	fmt.Println("hello hometic : I'm Gopher!!")

	r := mux.NewRouter()
	r.Handle("/pair-device", PairDeviceHandler(CreatePairDevice)).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr:", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("starting...")
	log.Fatal(server.ListenAndServe())
}
