package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/arknable/fwdproxy/config"
	"github.com/arknable/fwdproxy/env"
	"github.com/arknable/fwdproxy/handler"
	plog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
)

// Variables to be set at build
var (
	// IsProduction overrides config.IsProduction
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
		config.IsProduction = isProd
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

	srv := server.New(http.HandlerFunc(handler.HandleHTTP))
	log.Infof("Listening at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
