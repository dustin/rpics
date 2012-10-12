package main

import (
	baselog "log"
	"log/syslog"
	"os"
)

var log *baselog.Logger

func setupLogging(useSyslog bool) {
	if useSyslog {
		lw, err := syslog.New(syslog.LOG_INFO, "rpics")
		if err != nil {
			baselog.Fatalf("Error initializing logging: %v", err)
		}
		log = baselog.New(lw, "", 0)
	} else {
		log = baselog.New(os.Stderr, "", baselog.LstdFlags)
	}
}
