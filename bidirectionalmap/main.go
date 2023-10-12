package main

import (
	"fmt"

	"mymodulepath/bimap"
)

func main() {
	frToEn := bimap.New()
	frToEn.Store("un", "one")
	frToEn.Store("deux", "two")
	fmt.Println(frToEn)
}
