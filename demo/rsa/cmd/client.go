package cmd

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

func Client() {
	pool := LoadCA()

	log.SetPrefix("[Client]")
	cliCrt, err := tls.LoadX509KeyPair(ClientCrt, ClientKey)

	if err != nil {
		panic(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost:8081")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("resp:%s", string(body))
}
