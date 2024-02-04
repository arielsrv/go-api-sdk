package main

import (
	"log"

	"github.com/arielsrv/go-sdk-api/core"
)

func main() {
	server := core.NewServer()
	server.On(new(Application))
	server.AddHostedService(new(MyBackgroundWorker))
	server.Start()

	if err := server.Join(); err != nil {
		log.Fatal(err)
	}
}
