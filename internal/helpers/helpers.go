package helpers

import (
	"math/rand"
	"time"
)

func PaymentErrorImitation() bool {
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(60)
	b := rand.Intn(45)
	return a > b
}
