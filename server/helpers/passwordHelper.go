package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// 14 >>> cost
	// Bcrypt uses a cost parameter that specify the number of cycles to use in the algorithm.
	// Increasing this number the algorithm will spend more time to generate the hash output.
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(hashedPasswordBytes), err
}

func VerifyPassword(hashedPassword, userPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))

	return err
}
