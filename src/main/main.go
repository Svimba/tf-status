package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"mirantis.com/tungsten-operator/tf-status/src/status"
)

func handle(w http.ResponseWriter, r *http.Request) {
	tfStatus := status.TFStatus{}
	tfStatus.GetContrailStatus()       // get data from contrail-status
	fmt.Printf(tfStatus.PlainText)     // put data to console
	fmt.Fprintf(w, tfStatus.PlainText) // put data to http response
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	tfStatus := status.TFStatus{}

	hostname, err := os.Hostname() // get hostname
	if err != nil {
		log.Fatalf("Get hostname failed with %s\n", err)
	}

	tfStatus.PodName = hostname
	tfStatus.GetContrailStatus()
	tfStatus.ParseToJSON()

	out, err := json.Marshal(tfStatus)
	if err != nil {
		log.Fatalf("Marshal failed with %s\n", err)
	}
	fmt.Printf(tfStatus.PlainText) // put data to console
	fmt.Fprintf(w, string(out))    // put json to http response
}

func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/", handle)
	http.HandleFunc("/json", handleJSON)
	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = "80"
	}
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
