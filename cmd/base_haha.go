package cmd

import "fmt"

func printHaha(character string) {
	switch character {
	case "joker":
		fmt.Print(joker)
	case "trollface":
		fmt.Print(trollface)
	default:
		fmt.Print(trollface)
	}
}
