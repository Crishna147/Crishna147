package main

import (
	"fmt"
)

var m = map[string]string{
	"t": "Tanjiro",
	"n": "Nezuko",
	"z": "Zenitzu",
}

func main() {
	var anime = map[string]interface{}{
		"ANIME NAME": "Demon Slayer",
		"season":     3,
	}
	fmt.Println("Anime Name:", anime["ANIME NAME"])
	fmt.Println("Seasons:", anime["season"])

	fmt.Println(m["t"], "THE DEMON SLAYER")

	m["H"] = "HASHIRA"
	fmt.Println(m)

	for key, value := range m {
		fmt.Println(key, value)
	}
}
