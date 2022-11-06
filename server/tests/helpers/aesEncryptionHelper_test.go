package tests_helpers

import (
	"nft-raffle/helpers"
	"testing"
)

var (
	aesEncryptionHelper helpers.IAesEncrptionHelper = helpers.AesEncryptionHelper
)

func TestAesGcmEncryptionDecryption(t *testing.T) {
	s := "myString"
	encryptedString, err := aesEncryptionHelper.AesGCMEncrypt(s)

	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("Encrypted string: %s", encryptedString)

	decryptedString, err := aesEncryptionHelper.AesGCMDecrypt(encryptedString)

	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("Decrypted string: %s", decryptedString)

	if decryptedString != s {
		t.Error("Decrypted string is not matching")
	}

}
