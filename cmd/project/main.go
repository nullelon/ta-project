package main

import "ta-project-go/internal/app"

func main() {
	server, err := app.NewServer(":8080")
	if err != nil {
		panic(err)
	}

	panic(server.Start())
}
