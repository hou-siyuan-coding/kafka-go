package main

import (
	"flag"
	"log"
	"os"

	"github.com/hou-siyuan-coding/kafka-go/server"
	"github.com/hou-siyuan-coding/kafka-go/web"
)

var (
	filename = flag.String("filename", "", "The filename where to put all the data")
	inmem    = flag.Bool("inmem", false, "Whether or not use in-memory storage instead of a disk-based on")
	port     = flag.Uint("port", 0, "Listen to port")
)

func main() {
	flag.Parse()

	var backend web.Storage

	if *port == 0 {
		log.Fatalf("The flag `--port` must be provided")
	}

	if *inmem {
		backend = &server.InMemory{}
		log.Printf("In-Memory\n")
	} else {
		if *filename == "" {
			log.Fatalf("The flag `--filename` must be provided")
		}

		fp, err := os.OpenFile(*filename, os.O_CREATE|os.O_RDWR, 06666)
		if err != nil {
			log.Fatalf("Could not creat file %q: %s", *filename, err)
		}
		defer fp.Close()
		backend = server.NewOnDisk(fp)
		log.Printf("On-Disk\n")
	}

	log.Printf("Listening connections\n")
	web.NewServer(backend, *port).Serve()
}
