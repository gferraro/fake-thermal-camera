package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"os"
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
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/create/{device-name}", createDeviceHandler)
	router.HandleFunc("/triggerEvent/{type}", triggerEventHandler)
	router.HandleFunc("/playCPTVFile", sendCPTVFramesHandler)

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
		logError("'group-name' query parameter is missing", w, http.StatusBadRequest);
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
			logError(fmt.Sprintf("Error was %v", outputString), w, http.StatusInternalServerError);
		} else {
			log.Printf("Device created")
			restartThermalUploader()
			deviceID, err := getDeviceID()
			if err != nil {
				logError(fmt.Sprintf("Could not read device id %v", err), w, http.StatusInternalServerError);
			}
			io.WriteString(w, fmt.Sprintf("%d", deviceID))
		}
	}
}

func restartThermalUploader() {
	log.Printf("restarting thermal uploader")
	cmd := exec.Command("supervisorctl", "restart", "thermal-uploader")
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
		logError(fmt.Sprintf("Could not marshal json %s: %s", eventDetails, err), w, http.StatusInternalServerError);
		return
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		logError(fmt.Sprintf("Could not connect to dbus: %s", err), w, http.StatusInternalServerError);
		return
	}

	obj := conn.Object("org.cacophony.Events", "/org/cacophony/Events")
	call := obj.Call("org.cacophony.Events.Add", 0, string(detailsJSON), eventType, ts.UnixNano())
	if call.Err != nil {
		logError(fmt.Sprintf("Could not record %s event: %s", eventType, call.Err), w, http.StatusInternalServerError);
		return
	}
}

func sendCPTVFramesHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("cptv-file")
	if fileName == "" {
		fileName = "test.cptv"
	}

	fullFileName := "/cptv-files/" + fileName
	_, err := os.Stat(fullFileName)
	if os.IsNotExist(err) {
		logError(fmt.Sprintf("Could not find file with name %s", fileName), w, http.StatusBadRequest);
		return
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		logError(fmt.Sprintf("Could not connect to dbus: %v", err), w, http.StatusInternalServerError);
		return
	}
	obj := conn.Object("org.cacophony.FakeLepton", "/org/cacophony/FakeLepton")
	call := obj.Call("org.cacophony.FakeLepton.SendCPTV", 0, fullFileName)
	if call.Err != nil {
		logError(fmt.Sprintf("Could not send cptv %s: %s", fileName, call.Err), w, http.StatusInternalServerError);
		return
	}
	log.Printf("Sent CPTV Frames")
	io.WriteString(w, "Sent file + filename");
}

func logError(errorString string, w http.ResponseWriter, code int) {
	log.Printf("Error: %s", errorString)
	http.Error(w, fmt.Sprintf(errorString), code)
}

