package main

import (
	"fmt"
	"log"

	"github.com/Louis-Walker/hltv"
)

func main() {
	client := hltv.New()

	players, err := client.GetPlayers()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(players)
}
