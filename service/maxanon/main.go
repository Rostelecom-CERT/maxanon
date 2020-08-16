package main

import (
	"flag"
	"log"
	"os"

	"github.com/Rostelecom-CERT/maxanon"
	"github.com/peterbourgon/ff"
)

func main() {
	fs := flag.NewFlagSet("maxanon", flag.ExitOnError)
	var (
		apiPortPtr = fs.String("listen", ":8000", "Listen address")
		dbType     = fs.String("db", "redis", "Storage type")
		dbURL      = fs.String("db-url", "localhost:6379", "Storage URL")
		fileCSVPtr = fs.String("file", "", "File with database anonymous IP")
	)

	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("MAXANON"))
	if err != nil {
		log.Fatal(err)
	}
	app, err := maxanon.New(*fileCSVPtr, *dbType, *dbURL)
	if err != nil {
		panic(err)
	}
	err = app.Run(*apiPortPtr)
	if err != nil {
		log.Fatal(err)
	}
}
