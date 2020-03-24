package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/godbus/dbus"
	"github.com/gorilla/mux"

	config "github.com/TheCacophonyProject/go-config"
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
	router.HandleFunc("/triggerEvent/{type}", triggerEventHandler)
	router.HandleFunc("/sendCPTVFrames", sendCPTVFramesHandler)

	log.Fatal(http.ListenAndServe(":2040", router))
	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
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
			"password_"+deviceName,
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
			startThermalUploader()
			deviceID, err := getDeviceID()
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not read device id %v", err), http.StatusInternalServerError)
			}
			io.WriteString(w, fmt.Sprintf("%d", deviceID))
		}
	}
}

func startThermalUploader() {
	log.Printf("starting thermal uploader")
	cmd := exec.Command("supervisorctl", "start", "thermal-uploader")
	cmd.Start()
}
func getDeviceID() (int, error) {
	configRW, err := config.New(config.DefaultConfigDir)
	if err != nil {
		return 0, err
	}
	var deviceConfig config.Device
	if err := configRW.Unmarshal(config.DeviceKey, &deviceConfig); err != nil {
		return 0, err
	}
	return deviceConfig.ID, nil
}

func serverError(w *http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
}

func triggerEventHandler(w http.ResponseWriter, r *http.Request) {
	eventType := mux.Vars(r)["type"]
	eventDetails := map[string]interface{}{
		"description": map[string]interface{}{
			"type": eventType,
		},
	}
	ts := time.Now()
	detailsJSON, err := json.Marshal(&eventDetails)
	if err != nil {
		log.Printf("Could not record %s event: %s", eventType, err)
		return
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Could not record %s event: %s", eventType, err)
		return
	}

	obj := conn.Object("org.cacophony.Events", "/org/cacophony/Events")
	call := obj.Call("org.cacophony.Events.Add", 0, string(detailsJSON), eventType, ts.UnixNano())
	if call.Err != nil {
		log.Printf("Could not record %s event: %s", eventType, call.Err)
		return
	}
}

func newRecordingHandler(w http.ResponseWriter, r *http.Request) {
	eventType := mux.Vars(r)["type"]
	eventDetails := map[string]interface{}{
		"description": map[string]interface{}{
			"type": eventType,
		},
	}
	ts := time.Now()
	detailsJSON, err := json.Marshal(&eventDetails)
	if err != nil {
		log.Printf("Could not record %s event: %s", eventType, err)
		return
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Could not record %s event: %s", eventType, err)
		return
	}

	obj := conn.Object("org.cacophony.Events", "/org/cacophony/Events")
	call := obj.Call("org.cacophony.Events.Add", 0, string(detailsJSON), eventType, ts.UnixNano())
	if call.Err != nil {
		log.Printf("Could not record %s event: %s", eventType, call.Err)
		return
	}
}

func sendCPTVFramesHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("cptv-file")
	var extraCmds []string
	if fileName != "" {
		extraCmds = []string{"--cptv", fileName}
	}
	cmd := exec.Command("./fake-lepton", extraCmds...)
	cmd.Dir = "/server/cmd/fake-lepton"

	if output, err := cmd.CombinedOutput(); err != nil {
		outputString := string(output)
		log.Printf("Error was " + outputString)
		http.Error(w, outputString, http.StatusInternalServerError)
	} else {
		log.Printf("Sent CPTV Frames")
		io.WriteString(w, "Success")
	}
}
