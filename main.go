package main

import (
	"flag"
	"os"
)

var (
	listenAddr string
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", "", "server listen address")
	flag.Parse()

	// for heroku
	if listenAddr == "" && os.Getenv("PORT") != "" {
		listenAddr = ":" + os.Getenv("PORT")
	}

	// bootstrap any other packages or components here.
	// things like database connectivity or domain services
	// you can inject this as a dependency on your webserver if required
	// by adding another parameter to the NewServer func.
	// Now I left this for brevity.
	server, err := NewServer(listenAddr)
	if err != nil {
		panic(err)
	}
	server.Start()
}
