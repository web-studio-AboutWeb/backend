package wcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Encode(s []byte, block cipher.Block) ([]byte, error) {
	plainText, err := PKCS7Pad(s, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("padding error: %w", err)
	}

	// IV needs to be unique, but doesn't have to be secure.
	// It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("crypto.Encode(): %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func Decode(s []byte, block cipher.Block) ([]byte, error) {
	cipherText := make([]byte, len(s))
	copy(cipherText, s)

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext block size is too short")
	}

	// IV needs to be unique, but doesn't have to be secure.
	// It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	res, err := PKCS7Unpad(cipherText, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("unpadding error (probably, invalid AES key): %w", err)
	}

	return res, nil
}

func EncodeToBase64(s []byte, block cipher.Block) (string, error) {
	res, err := Encode(s, block)
	return base64.URLEncoding.EncodeToString(res), err
}

func DecodeFromBase64(s string, block cipher.Block) ([]byte, error) {
	resBytes, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	res, err := Decode(resBytes, block)
	return res, err
}

func DecodeUserPass(username, password string, block cipher.Block) (string, string, error) {
	usernameDecoded, err := DecodeFromBase64(username, block)
	if err != nil {
		return "", "", fmt.Errorf("decoding username: %w", err)
	}

	passwordDecoded, err := DecodeFromBase64(password, block)
	if err != nil {
		return "", "", fmt.Errorf("decoding password: %w", err)
	}

	return string(usernameDecoded), string(passwordDecoded), nil
}

func EncodeUserPass(username, password string, block cipher.Block) (string, string, error) {
	usernameEncoded, err := EncodeToBase64([]byte(username), block)
	if err != nil {
		return "", "", fmt.Errorf("encoding username: %w", err)
	}

	passwordEncoded, err := EncodeToBase64([]byte(password), block)
	if err != nil {
		return "", "", fmt.Errorf("encoding password: %w", err)
	}

	return usernameEncoded, passwordEncoded, nil
}
