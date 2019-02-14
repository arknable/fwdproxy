package main

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/arknable/fwdproxy/config"
	"github.com/arknable/fwdproxy/env"
	plog "github.com/arknable/fwdproxy/log"
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

	// handlerFunc := http.HandlerFunc(handler.HandleRequest)
	// tlssrv := server.NewTLS(handlerFunc)
	// go func() {
	// 	log.Printf("Starting HTTPS Server at %s ...\n", tlssrv.Addr)
	// 	if err := tlssrv.ListenAndServeTLS(config.CertPath, config.KeyPath); err != nil {
	// 		log.Fatal("HTTPS Error: ", err)
	// 	}
	// }()

	// srv := server.New(handlerFunc)
	// log.Printf("Starting HTTP Server at %s ...\n", srv.Addr)
	// if err := srv.ListenAndServe(); err != nil {
	// 	log.Fatal("HTTP Error: ", err)
	// }
}
