package probs

import "math/rand"

func GetPercentProb(rate int) bool {
	if rate < 1 || rate > 100 {
		return false
	}
	num := rand.Intn(100)
	return num <= rate-1
}
