package main

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/mocks"
	"log"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core"
	_ "gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/docs"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/examples/full/app"
)

// @title Backend IskayPet SDK
// @description Provide an interface to build core APIs.
// @basePath /
// @version v1.
func main() {
	// create a new server
	server := core.NewServer()

	// Decoupled server from application
	server.On(new(app.Application))

	// Optionally, you can register a background worker when the server starts
	server.AddHostedService(new(mocks.DummySQSConsumerService))

	// Start internal server
	server.Start()

	// Thread join
	if err := server.Join(); err != nil {
		log.Fatal(err)
	}
}
