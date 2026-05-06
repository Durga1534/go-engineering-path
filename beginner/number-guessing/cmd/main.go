package main

import (
	"fmt"
	"number-guessing/internal/guess"
)

func main() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")

	fmt.Println("\nPlease select the difficulty level:")
	fmt.Println("1. Easy(10 chances)\n2. Medium(5 chances)\n3. Hard(3 chances)")

	var choice int
	fmt.Print("\nEnter your choice:")
	fmt.Scan(&choice)

	diffName, chances := guess.GetDifficulty(choice)
	target := guess.GenerateNumber()

	fmt.Printf("\nGreat! You have selected the %s difficulty level.\n", diffName)
	fmt.Println("Let's start the game!")

	for attempts := 1; attempts <= chances; attempts++ {
		var guess int
		fmt.Printf("\nEnter your guess (Attempt %d/%d): ", attempts, chances)
		fmt.Scan(&guess)

		if guess == target {
			fmt.Printf("Congratulations! You guessed the correct number in %d attempts. \n", attempts)
			return
		} else if guess < target {
			fmt.Println("Incorrect! The number is greater than", guess)
		} else {
			fmt.Println("Incorrect! The number is less than", guess)
		}
	}
	fmt.Printf("\nGame Over! You've run out of chances. The number was %d.\n", target)
}
