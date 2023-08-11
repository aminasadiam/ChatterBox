package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckComparePass(password, hashedpass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpass), []byte(password))
	return err == nil
}
