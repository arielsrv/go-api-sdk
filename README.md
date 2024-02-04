[![pipeline status](https://gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/badges/main/pipeline.svg)](https://gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/-/commits/main)
[![coverage report](https://gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/badges/main/coverage.svg)](https://gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/-/commits/main)
[![release](https://gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/-/badges/release.svg)](https://gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/-/releases)

> This SDK provides a framework to build APIs

The intent of the project is to provide a lightweight api sdk, based on Golang.

The main goal is to provide a modular framework with high-level abstractions to expose RESTful http endpoints.

## Table of contents

* [Download](#download)
* [Example](#example)
    * [Default](#default)
    * [Config](#config)
* [How it works](#how-it-works)

# Download
```shell
go get gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2
```

# Example

full
example [here](https://gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/-/tree/feature/views/examples/full?ref_type=heads)

## Default

routes.go
```go
package main

import (
	"net/http"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core"
)

type Routes struct {
	core.APIRoutes
}

func (r *Routes) Register() {
	r.AddRoute(http.MethodGet, "/message", core.Provide[IMessageController]().GetMessage)
}
```

application.go
```go
package main

import "gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core"

type Application struct {
	core.APIApplication
}

func (r *Application) Init() {
	r.UseMetrics()
	r.UseSwagger()
	r.RegisterDependencyInjectionModule(new(ApplicationModule))
	r.RegisterRoutes(new(Routes))
	r.Build()
}
```

application_module.go
```go
package main

import (

"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/container"
"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/examples/full/app/controllers"
"go.uber.org/dig"
)


type ApplicationModule struct {
	container.DependencyInjectionModule
}

func (r *ApplicationModule) Configure() {
	r.Bind(controllers.NewMessageController, dig.As(new(controllers.IMessageController)))
}

```

controller.go
```go
package main
import "gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core/routing"


type IMessageController interface {
	GetMessage(ctx *routing.HTTPContext) error
}

type MessageController struct {
}

func NewMessageController() *MessageController {
	return &MessageController{}
}

func (r *MessageController) GetMessage(ctx *routing.HTTPContext) error {	
    return ctx.SendString("Hello World")
}

```

main.go
```go
package main

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/examples/full/app"
"log"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/backend-api-sdk/v2/core"
)

func main() {
	server := core.NewServer()
	server.On(new(app.Application))	
	server.Start()

	if err := server.Join(); err != nil {
		log.Fatal(err)
	}
}
```


## Config

Based on Netflix Archaius (default folder recommended)

* src/config/config.yaml
* src/config/local/config.yaml
* src/config/remote/[env].config.yaml

# How it works

Features
* Dashboard & Prometheus metrics (optional)
* Configurable readiness warmupper for /ping (200 or 503) (optional)
* Automatic config discovery (src/config/config.yaml ... see sites-api on gitlab) (optional)
* Dependency Injection Module from Uber (optional)
* Swagger (optional)
* HTML Views (optional)
