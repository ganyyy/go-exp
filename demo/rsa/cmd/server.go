package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

var echoHandle http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	log.Printf("receive client %v", r.RemoteAddr)
	fmt.Fprintf(w, "Hi, this is a example of http service in golang!\n")
}

func Server() {
	pool := LoadCA()
	log.SetPrefix("[Server]")

	var s = http.Server{
		Addr:    ":8081",
		Handler: echoHandle,
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	err := s.ListenAndServeTLS(ServerCrt, ServerKey)
	if err != nil {
		panic(err)
	}
}
