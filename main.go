package main

import (
	"os"

	"github.com/DJDNS/djdns/app"
)

func main() {
	app.Main(nil, true, os.Stderr)
}
