package test_helpers

import (
	"nft-raffle/helpers"
	"testing"
)

var (
	randomCodeGenerator helpers.RandomCodeGenerator = helpers.NewRandomCodeGenerator()
)

func TestGenerateRandomDigits(t *testing.T) {
	s := randomCodeGenerator.GenerateRandomDigits(6)

	if len(s) != 6 {
		t.Error("Length of generated random 6 digits is not 6")
	}

}
