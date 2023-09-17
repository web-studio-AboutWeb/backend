package wcrypto

import (
	"bytes"
	"crypto/aes"
	"testing"
)

func TestPadding1(t *testing.T) {
	str1 := []byte("12345678")
	padded1, err := PKCS7Pad(str1, aes.BlockSize)
	if err != nil {
		t.Fatal(err)
	}

	expected1 := append(str1, []byte{8, 8, 8, 8, 8, 8, 8, 8}...)
	if !bytes.Equal(padded1, expected1) {
		t.Fatalf("12345678 padding failed: expected %v, got %v", expected1, padded1)
	}

	unpadded1, err := PKCS7Unpad(padded1, aes.BlockSize)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(unpadded1, str1) {
		t.Fatalf("12345678 unpadding failed: expected %v, got %v", str1, unpadded1)
	}
}

func TestPadding2(t *testing.T) {
	str1 := []byte("")
	padded1, err := PKCS7Pad(str1, aes.BlockSize)
	if err != nil {
		t.Fatal(err)
	}

	expected1 := []byte{16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16}
	if !bytes.Equal(padded1, expected1) {
		t.Fatalf("empty string padding failed: expected %v, got %v", expected1, padded1)
	}

	unpadded1, err := PKCS7Unpad(padded1, aes.BlockSize)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(unpadded1, str1) {
		t.Fatalf("empty string unpadding failed: expected %v, got %v", str1, unpadded1)
	}
}
