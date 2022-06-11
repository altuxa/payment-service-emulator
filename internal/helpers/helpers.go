package helpers

import (
	"math/rand"
	"net/mail"
	"time"
)

func PaymentErrorImitation() bool {
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(60)
	b := rand.Intn(45)
	return a > b
}

func FailStatusImitation() bool {
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(60)
	b := rand.Intn(30)
	return a > b
}

func ValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
