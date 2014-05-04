package main

//import "flag"
import (
	"github.com/campadrenalin/djdns/model"
	"github.com/campadrenalin/djdns/server"
	"log"
)

func main() {
	filename := "./model/demo.json"
	log.Printf("Converting %s", filename)
	json, err := model.GetJSONFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	addr := "127.0.0.1:9953"
	s := server.NewServer()
	s.Root.LoadFrom(json)
	log.Printf("Starting server on %s", addr)
	err = s.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
