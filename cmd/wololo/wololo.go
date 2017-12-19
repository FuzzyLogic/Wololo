package main

import (
	"log/syslog"
	"net/http"
	wololo "github.com/FuzzyLogic/Wololo"
)

var globalConfig wololo.WololoConfig
var globalLog *syslog.Writer

func main() {
	// Start logging and defer connection close to end
	globalLog = wololo.SetupLog()
	defer func() {
		if err := globalLog.Close(); err != nil {
			panic(err)
		}
	}()

	// Apply seccomp sandbox to application if activated
	err := wololo.Sandbox()
	if err != nil {
		wololo.WriteToLog(globalLog, "Error sandboxing application")
		panic(err)
	}

	// Read configuration into global variable
	localConfig, err := wololo.ReadConfig("/etc/wololo/wololo.conf")
	if err != nil {
		wololo.WriteToLog(globalLog, "Error reading configuration")
		panic(err)
	}
	globalConfig = *localConfig

	// Start HTTP handler
	wololo.WriteToLog(globalLog, "Starting server")
	http.HandleFunc("/", wolHandler)
	if err := http.ListenAndServe(globalConfig.ListenIP+":"+globalConfig.ListenPort, nil); err != nil {
		wololo.WriteToLog(globalLog, "Unable to start server")
		panic(err)
	}
}
