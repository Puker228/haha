package cmd

import "fmt"

func printHaha(character string) {
	fmt.Print(characterArt(character))
}

func characterArt(character string) string {
	switch character {
	case "joker":
		return joker
	case "trollface":
		return trollface
	default:
		return trollface
	}
}
