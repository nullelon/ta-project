package main

import "ta-project-go/internal/app"

func main() {
	server, err := app.NewServer(":8087")
	if err != nil {
		panic(err)
	}

	panic(server.Start())
}
