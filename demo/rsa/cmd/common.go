package cmd

import (
	"crypto/x509"
	"os"
)

const (
	Root   = "./keys/"
	CAFile = Root + "ca.crt"

	ServerCrt = Root + "server.crt"
	ServerKey = Root + "server.key"

	ClientCrt = Root + "client.crt"
	ClientKey = Root + "client.key"
)

func LoadCA() *x509.CertPool {
	pool := x509.NewCertPool()
	caCrt, err := os.ReadFile(CAFile)
	if err != nil {
		panic(err)
	}
	pool.AppendCertsFromPEM(caCrt)
	return pool
}
