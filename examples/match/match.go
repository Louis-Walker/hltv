package main

import (
	"fmt"
	"log"

	"github.com/Louis-Walker/hltv"
)

func main() {
	client := hltv.New()

	match, err := client.GetMatch(2359847, "websterz-vs-copenhagen-flames-cct-south-europe-series-1")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(match)
}
