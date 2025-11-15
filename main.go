package main

import (
	"fb-search/views"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	di, err := views.CreateDi()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start app: %v\n", err)
		os.Exit(1)
	}

	app := di.Get(views.HttpServerDef).(*views.HttpServer)

	app.Run()
}
