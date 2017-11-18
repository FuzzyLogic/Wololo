package main

import (
	"log/syslog"
)

// Connect to the syslog daemon if possible
func setupLog() *syslog.Writer {
	// Establish connection to syslog daemon if possible
	logWriter, _ := syslog.New(syslog.LOG_NOTICE, "wololo")
	return logWriter
}

// Write a message to the syslog daemon if a connection was established
func writeToLog(log *syslog.Writer, msg string) {
	if log != nil {
		log.Write([]byte(msg))
	}
}
