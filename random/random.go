package random

import (
	"math/rand"
)

func RandomNumber(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
