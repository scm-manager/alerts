package main

import (
	"fmt"
	"github.com/scm-manager/alerts/src/alert"
	"github.com/scm-manager/alerts/src/api"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("use alerts directory")
		os.Exit(1)
	}

	directory := os.Args[1]
	log.Printf("Read alerts from directory %s", directory)

	alerts, err := alert.ReadFromDirectory(directory)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.Handle("/ready", api.CreateOkEndpoint())
	router.Handle("/live", api.CreateOkEndpoint())
	router.Handle("/api/v1/alerts", api.CreateAlertsEndpoint(alerts))

	log.Println("start http server on 8080 ...")
	log.Println("")
	log.Println("http://localhost:8080/api/v1/alerts")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to start http server", err)
	}
}
