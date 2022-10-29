package test_helpers

import (
	"nft-raffle/helpers"
	"testing"
)

var (
	randomCodeGenerator helpers.IRandomCodeGenerator = helpers.RandomCodeGenerator
)

func TestGenerateRandomDigits(t *testing.T) {
	s := randomCodeGenerator.GenerateRandomDigits(6)

	if len(s) != 6 {
		t.Error("Length of generated random 6 digits is not 6")
	}

}
