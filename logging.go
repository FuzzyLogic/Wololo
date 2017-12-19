package wololo

import (
	"log/syslog"
)

// Connect to the syslog daemon if possible
func SetupLog() *syslog.Writer {
	// Establish connection to syslog daemon if possible
	logWriter, _ := syslog.New(syslog.LOG_NOTICE, "wololo")
	return logWriter
}

// Write a message to the syslog daemon if a connection was established
func WriteToLog(log *syslog.Writer, msg string) {
	if log != nil {
		log.Write([]byte(msg))
	}
}
