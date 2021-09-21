package main

import (
	"log"
	"os"
)

var stdLog, errLog *log.Logger

func init() {
	stdLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errLog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}
