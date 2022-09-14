package helpers

import (
	"math/rand"
	"strings"
	"time"
)

type RandomCodeGenerator interface {
	GenerateRandomDigits(length int) string
}

type randomCodeGeneratorStruct struct{}

var digits = []rune("0123456789")

func NewRandomCodeGenerator() RandomCodeGenerator {
	return &randomCodeGeneratorStruct{}
}

func (r *randomCodeGeneratorStruct) GenerateRandomDigits(length int) string {
	rand.Seed(time.Now().UnixNano())
	digitSize := len(digits)
	var sb strings.Builder

	for i := 0; i < length; i++ {
		char := digits[rand.Intn(digitSize)]
		sb.WriteRune(char)
	}

	s := sb.String()
	return s
}
