package main

import (
	"fmt"
	"log"

	"github.com/Louis-Walker/hltv"
)

func main() {
	client := hltv.New()

	match, err := client.GetMatch(2359958)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(match)
}
