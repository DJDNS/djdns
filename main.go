package main

//import "flag"
import "log"

func main() {
	filename := "demo.json"
	log.Printf("Converting %s", filename)
	json, err := GetJSONFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(json)
}
