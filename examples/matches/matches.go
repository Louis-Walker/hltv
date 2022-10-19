package main

import (
	"fmt"
	"log"

	"github.com/Louis-Walker/hltv"
)

func main() {
	client := hltv.New()

	matches, err := client.GetMatches()
	if err != nil {
		log.Fatal(err.Error())
	}

	if matches != nil {
		fmt.Println("Matches:")
		for _, item := range matches {
			fmt.Println("", item)
		}
	}
}
