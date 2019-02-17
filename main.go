package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/arknable/fwdproxy/env"
	"github.com/arknable/fwdproxy/handler"
	plog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
)

// Variables to be set at build
var (
	// IsProduction overrides env.IsProduction
	IsProduction = "false"

	// BuiltInUsername overrides user.BuiltInUsername
	BuiltInUsername = "admin"

	// BuiltInUserPwd overrides user.BuiltInUserPwd
	BuiltInUserPwd = "4dm1n"
)

func main() {
	// // Dump log to standard output and file
	// file, err := os.OpenFile("proxy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("Can't open log file: %v", err)
	// }
	// defer file.Close()
	// log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetReportCaller(true)
	log.SetFormatter(&plog.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	// Init environment
	if err := env.Initialize(); err != nil {
		log.Fatal(err)
	}

	// Check config overrides
	isProd, err := strconv.ParseBool(IsProduction)
	if err != nil {
		log.Error(err)
	} else {
		env.IsProduction = isProd
	}
	user.BuiltInUsername = BuiltInUsername
	user.BuiltInUserPwd = BuiltInUserPwd
	log.WithFields(log.Fields{
		"IsProduction":    isProd,
		"BuiltInUsername": BuiltInUsername,
		"BuiltInUserPwd":  BuiltInUserPwd,
	}).Info("Configuration overrides")

	// Init repository
	repo, err := user.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	err = server.Initialize(http.HandlerFunc(handler.HandleRequest))
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Listening at %s", server.Port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
