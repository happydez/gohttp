package main

import (
	"log"

	"github.com/happydez/gohttp.git/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
