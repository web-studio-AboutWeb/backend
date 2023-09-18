package strhelp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
)

func GenerateRandomString(size int) (string, error) {
	res := make([]byte, size)

	if _, err := io.ReadFull(rand.Reader, res); err != nil {
		return "", fmt.Errorf("string generation error: %w", err)
	}

	return base64.StdEncoding.EncodeToString(res), nil
}

const emailRegexp = "^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$"

func ValidateEmail(email string) bool {
	return regexp.MustCompile(emailRegexp).MatchString(email)
}
