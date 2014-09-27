package main

import (
	"log"

	"github.com/DJDNS/djdns/app"
)

func main() {
	if err := app.Main(nil, true); err != nil {
		log.Fatal(err)
	}
}
