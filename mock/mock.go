// Functions for creating fake data

package main

import (
	"github.com/battleofbits/arena/arena"
	"log"
)

const URL = "http://localhost:5000/fourup"

func checkError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() {
	_, redErr := arena.CreatePlayer("Kevin Burke", "kevinburke", URL, URL)
	checkError(redErr)
	_, blackErr := arena.CreatePlayer("Kyle Conroy", "kyleconroy", URL, URL)
	checkError(blackErr)
}
