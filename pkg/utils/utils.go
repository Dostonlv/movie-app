package utils

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"strings"
)

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func GenerateOTP() string {
	// generate 6 digit otp
	var OTP string
	for i := 1; i <= 5; i++ {
		OTP += strconv.Itoa(rand.Intn(10))
	}
	return OTP
}

func CompareHash(hashedpassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
}
