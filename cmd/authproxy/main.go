// Package main provides a runnable server for authenticating a reversed proxy.
// It takes three environment variables as arguments: PORT, REVERSE_PROXY_URL and SERVICE_ACCOUNT_KEY
// SERVICE_ACCOUNT_KEY must be base64 encoded. PORT defaults to 9090
package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"marwan.io/authproxy"
)

var config struct {
	ProxyURL       string `envconfig:"REVERSE_PROXY_URL" required:"true"`
	ServiceAccount string `envconfig:"SERVICE_ACCOUNT_KEY" required:"true"`
	Port           string `envconfig:"PORT" default:"9090"`
}

func init() {
	envconfig.MustProcess("", &config)
	if config.Port[0] != ':' {
		config.Port = ":" + config.Port
	}
}

func main() {
	log.Println("listening on port " + config.Port)
	handler, err := authproxy.GetProxyHandler(config.ProxyURL, config.ServiceAccount)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(config.Port, handler))
}
