package main

import (
	baselog "log"
	"log/syslog"
	"os"
)

var log *baselog.Logger

func setupLogging(useSyslog bool) {
	if useSyslog {
		var err error
		log, err = syslog.NewLogger(syslog.LOG_INFO, 0)
		if err != nil {
			baselog.Fatalf("Error initializing logging: %v", err)
		}
	} else {
		log = baselog.New(os.Stderr, "", baselog.LstdFlags)
	}
}
