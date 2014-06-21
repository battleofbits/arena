// Functions for creating fake data

package main

import (
	"github.com/battleofbits/arena/arena"
	"log"
)

const URL = "http://localhost:5000"

func checkError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() {
	_, redErr := arena.CreatePlayer("Kevin Burke", "kevinburke", URL+"/fourup", URL+"/invite")
	checkError(redErr)
	_, blackErr := arena.CreatePlayer("Kyle Conroy", "kyleconroy", URL+"/fourup", URL+"/invite")
	checkError(blackErr)
}
