package guess

import (
	"math/rand"
	"time"
)

func GetDifficulty(choice int) (string, int) {
	switch choice {
	case 1:
		return "Easy", 10
	case 3:
		return "Hard", 3
	default:
		return "Medium", 5
	}
}

func GenerateNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100) + 1
}
