package main

import (
	"fmt"
	"log"

	"github.com/Louis-Walker/hltv"
)

func main() {
	client := hltv.New()

	results, err := client.GetResults()
	if err != nil {
		log.Fatal(err.Error())
	}

	if results != nil {
		fmt.Println("Results:")
		for _, item := range results {
			fmt.Println("", item.MatchID)
		}
	}
}
