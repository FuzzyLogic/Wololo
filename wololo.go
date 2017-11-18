package main

import (
	"log/syslog"
	"net/http"
)

var globalConfig WololoConfig
var globalLog *syslog.Writer

func main() {
	// Start logging and defer connection close to end
	globalLog = setupLog()
	defer func() {
		if err := globalLog.Close(); err != nil {
			panic(err)
		}
	}()

	// Read configuration into global variable
	localConfig, err := readConfig("/etc/wololo/wololo.config")
	globalConfig = *localConfig
	if err != nil {
		writeToLog(globalLog, "Error reading configuration")
		panic(err)
	}

	// Start HTTP handler
	writeToLog(globalLog, "Starting server")
	http.HandleFunc("/", wolHandler)
	if err := http.ListenAndServe(globalConfig.listenIP+":"+globalConfig.listenPort, nil); err != nil {
		writeToLog(globalLog, "Unable to start server")
		panic(err)
	}
}
