package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/arknable/fwdproxy/config"
	"github.com/arknable/fwdproxy/handler"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
)

// Variables to be set at build
var (
	// IsProduction overrides config.IsProduction
	IsProduction = "false"

	// TLSAllowedHost overrides config.TLSAllowedHost
	TLSAllowedHost = "localhost"

	// BuiltInUsername overrides user.BuiltInUsername
	BuiltInUsername = "admin"

	// BuiltInUserPwd overrides user.BuiltInUserPwd
	BuiltInUserPwd = "4dm1n"
)

func main() {
	// Check config overrides
	isProd, err := strconv.ParseBool(IsProduction)
	if err != nil {
		log.Printf("Invalid value for IsProduction: %s", IsProduction)
	} else {
		config.IsProduction = isProd
	}
	config.TLSAllowedHost = TLSAllowedHost
	user.BuiltInUsername = BuiltInUsername
	user.BuiltInUserPwd = BuiltInUserPwd

	log.Println("Using following configurations:")
	log.Println("-------------------------------")
	log.Println("IsProduction	: ", IsProduction)
	log.Println("TLSAllowedHost	: ", TLSAllowedHost)
	log.Println()

	// Dump log to standard output and file
	file, err := os.OpenFile("proxy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Can't open log file: %v", err)
	}
	defer file.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	// Make sure repository cleaned
	defer user.Repo().Close()

	handlerFunc := http.HandlerFunc(handler.HandleRequest)
	tlssrv := server.NewTLS(handlerFunc)
	go func() {
		log.Printf("Starting HTTPS Server at %s ...\n", tlssrv.Addr)
		if err := tlssrv.ListenAndServeTLS(config.CertPath, config.KeyPath); err != nil {
			log.Fatal("HTTPS Error: ", err)
		}
	}()

	srv := server.New(handlerFunc)
	log.Printf("Starting HTTP Server at %s ...\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("HTTP Error: ", err)
	}
}
