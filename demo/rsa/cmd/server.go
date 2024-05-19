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
	log.SetPrefix("[Server]")

	var s = http.Server{
		Addr:    ":8081",
		Handler: echoHandle,
		TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
			Certificates: []tls.Certificate{func() tls.Certificate {
				crt, err := tls.LoadX509KeyPair(ServerCrt, ServerKey)
				if err != nil {
					panic(err)
				}
				return crt
			}()},
		},
	}

	err := s.ListenAndServeTLS(ServerCrt, ServerKey)
	if err != nil {
		panic(err)
	}
}
