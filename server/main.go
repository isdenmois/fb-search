package main

import (
	"fmt"
	"os"

	deliveryHttp "fb-search/delivery/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	di, err := deliveryHttp.CreateDi()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start app: %v\n", err)
		os.Exit(1)
	}

	app := di.Get(deliveryHttp.HttpServerDef).(*deliveryHttp.HttpServer)

	app.Run()
}
