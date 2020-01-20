package main

import (
	"log"
	"net/http"
	"io"
	"os/exec"

	"github.com/gorilla/mux"
)

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}

var version = "<not set>"

func runServer() error {
	log.SetFlags(0)

	router := mux.NewRouter()
	router.HandleFunc("/create/{device-name}", createDeviceHandler)
	router.HandleFunc("/", homeHandler)

	log.Fatal(http.ListenAndServe(":2040", router))
	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("device created")
	io.WriteString(w, "This is a Fake thermal camera test server.")
}

func createDeviceHandler(w http.ResponseWriter, r *http.Request) {
	deviceName := mux.Vars(r)["device-name"]
	groupnames, ok := r.URL.Query()["group-name"]
	if !ok {
		log.Printf("'group-name' query parameter is missing")
		http.Error(w, "'group-name' query parameter is missing", http.StatusBadRequest)
	} else {
		apiServers, ok := r.URL.Query()["api-server"]
		apiServer := "https://api-test.cacophony.org.nz"
		if ok {
			apiServer = apiServers[0]
		}
		log.Printf("Creating device " + deviceName + " and group-name " + groupnames[0] + " on server " + apiServer)
		cmd := exec.Command("./device-register",
												"--name",
												deviceName,
												"--group",
												groupnames[0],
												"--password",
												"password_" + deviceName,
												"--api",
												apiServer,
												"--ignore-minion-id",
												"--remove-device-config")

		cmd.Dir = "/code/device-register"

		if output, err := cmd.CombinedOutput(); err != nil {
			outputString := string(output)
			log.Printf("Error was " + outputString)
			http.Error(w, outputString, http.StatusInternalServerError)
		} else {
			log.Printf("device created")
			io.WriteString(w, "Success")
		}
	}
}

func serverError(w *http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
}
