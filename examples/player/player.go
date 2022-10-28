package main

import (
	"fmt"
	"log"

	"github.com/Louis-Walker/hltv"
)

func main() {
	client := hltv.New()

	player, err := client.GetPlayer(429)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(player)
}
