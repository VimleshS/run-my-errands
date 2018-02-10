package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VimleshS/run-my-errands/orderer/handlers"
	"github.com/VimleshS/run-my-errands/setup"
	que "github.com/bgentry/que-go"
	"github.com/jackc/pgx"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	qc      *que.Client
	pgxpool *pgx.ConnPool
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		logrus.WithField("PORT", port).Fatal("$PORT must be set")
	}

	var err error
	pgxpool, qc, err = setup.SetUp.PoolAndQueueConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer pgxpool.Close()
	handlers.InjectQc(qc)

	router := mux.NewRouter()
	logrus.Info("Starting the application...")
	router.HandleFunc("/authenticate", handlers.CreateTokenEndpoint).Methods("POST")
	router.HandleFunc("/uploadlist", handlers.ValidateToken(handlers.GroceryUploadList)).Methods("POST")
	logrus.Errorln(http.ListenAndServe(":"+port, router))
}
