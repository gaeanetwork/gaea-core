package user

import "github.com/pkg/errors"

func checkUsernameLen(username string) error {
	len := len(username)
	if len <= 0 || len > 32 {
		return errors.Errorf("Invalid username length, should be (0, 32], now: %d", len)
	}

	return nil
}

func checkPasswordLen(password string) error {
	len := len(password)
	if len <= 0 || len > 64 {
		return errors.Errorf("Invalid password length, should be (0, 64], now: %d", len)
	}

	return nil
}
