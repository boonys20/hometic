package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boonys20/hometic/logger"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pair struct {
	DeviceID int64
	UserID   int64
}

type Device interface {
	Pair(p Pair) error
}

type CustomHandlerFunc func(CustomRespWriter, *http.Request)

func (handler CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(&JSONRespWriter{w}, r)
}

type CustomRespWriter interface {
	JSON(statusCode int, date interface{})
}

type JSONRespWriter struct {
	http.ResponseWriter
}

func (w *JSONRespWriter) JSON(statusCode int, data interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(data)
}

type CreatePairDeviceFunc func(p Pair) error

func (fn CreatePairDeviceFunc) Pair(p Pair) error {
	return fn(p)
}

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func NewCreatePairDevice(db DB) CreatePairDeviceFunc {
	//func NewCreatePairDevice(db *sql.DB) CreatePairDeviceFunc {
	return func(p Pair) error {
		_, err := db.Exec("INSERT INTO pairs VALUES ($1,$2);", p.DeviceID, p.UserID)
		return err
	}
	//}
}

func PairDeviceHandler(device Device) func(w CustomRespWriter, r *http.Request) {

	return func(w CustomRespWriter, r *http.Request) {

		logger.L(r.Context()).Info("pair-device")

		var p Pair
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.JSON(http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()
		fmt.Printf("pair: %#v\n", p)

		err = device.Pair(p)
		if err != nil {
			w.JSON(http.StatusBadRequest, err.Error())
			return
		}

		w.JSON(http.StatusOK, []byte(`{"status":"active"}`))

	}

}

func run() error {
	fmt.Println("hello hometic : I'm Gopher!!")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	r := mux.NewRouter()
	r.Handle("/pair-device", CustomHandlerFunc(PairDeviceHandler(NewCreatePairDevice(db)))).Methods(http.MethodPost)

	r.Use(logger.MiddleWare)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Println("addr:", addr)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("starting...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal("Can't Start application", err)
	}
}
