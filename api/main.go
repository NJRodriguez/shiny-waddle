package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NJRodriguez/shiny-waddle/api/setup"
	"github.com/gorilla/mux"
)

var server *http.Server = buildHTTPServer()

func main() {
	log.Println("Starting server...")
	tableName := os.Getenv("TABLE_NAME")
	region := os.Getenv("AWS_REGION")

	server := &setup.Server{
		Router: mux.NewRouter(),
	}
	err := server.Initialize(tableName, region)
	if err != nil {
		log.Fatal("Error when trying to start server!")
		panic(err)
	}
	server.Run(":80")
	log.Println("Server started successfully!")
}

func buildHTTPServer() *http.Server {
	router := mux.NewRouter()
	return &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
