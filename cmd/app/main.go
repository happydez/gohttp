package main

import (
	"log"

	"github.com/happydez/gohttp.git/cmd"
	"github.com/happydez/gohttp.git/internal/server"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}

	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
