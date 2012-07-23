package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	webserver "./webserver"
)

func main() {
	fmt.Printf(">> Starting Media Streamer Webserver (Go Edition v. " + webserver.Version + ")\n")

	pwd, err := os.Getwd()

	if err != nil {
		fmt.Printf("FATAL: Could not get current working directory!\n")

		return
	}

	webserver.PublicDirectory = flag.String("d", pwd + "/public", "Public Directory")
	webserver.ListenPort = flag.String("p", "4568", "Listen Port")
	webserver.MediaDirectory = flag.String("m", "/Users/seshbaugh/Music/iTunes/iTunes Media/", "Media Directory")

	flag.Parse()

	fmt.Printf(">> Go application starting on http://0.0.0.0:" + *webserver.ListenPort + "\n")
	fmt.Printf(">> ctrl+c to shutdown server\n")
	fmt.Printf(">> pid=" + strconv.Itoa(os.Getpid()) + "\n")

	http.ListenAndServe(":" + *webserver.ListenPort, http.HandlerFunc(webserver.Handler))
}