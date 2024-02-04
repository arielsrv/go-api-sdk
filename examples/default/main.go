package main

import (
	"log"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core"
)

func main() {
	server := core.NewServer()
	server.Start()

	if err := server.Join(); err != nil {
		log.Fatal(err)
	}
}
