package main

import (
	"fmt"
)

type anime struct {
	animeName string
	genre     string
	seasons   int
}
type hunters struct {
	hunter1 string
	hunter2 string
	hunter3 string
	anime   anime
}

func main() {
	h := hunters{
		hunter1: "GON",
		hunter2: "KILLUA",
		hunter3: "LEORIO",
		anime: anime{
			animeName: "hunterXhunter",
			genre:     "violence",
			seasons:   6,
		},
	}
	fmt.Println("hunter1:", h.hunter1)
	fmt.Println("hunter2:", h.hunter2)
	fmt.Println("hunter3:", h.hunter3)
	fmt.Println("animeName:", h.anime.animeName)
	fmt.Println("genre:", h.anime.genre)
	fmt.Println("seasons:", h.anime.seasons)
}





-----------------------------------------------------------
PS D:\go\sample1> go run nestedstruct.go
hunter1: GON
hunter2: KILLUA
hunter3: LEORIO
animeName: hunterXhunter
genre: violence
seasons: 6




