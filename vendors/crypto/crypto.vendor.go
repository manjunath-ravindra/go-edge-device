package CryptoVendors

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// DecryptResponse decrypts the hex-encoded encryptedText using the hex-encoded iv and the secret key (hex string, AES-256-CBC), logs, and returns the decrypted string.

func DecryptResponse(transactionId, encryptedTextHex, ivHex, keyHex string) (string, error) {
	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode IV hex: %w", err)
	}
	ciphertext, err := hex.DecodeString(encryptedTextHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted text hex: %w", err)
	}
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode key hex: %w", err)
	}
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("secret key must be 32 bytes for AES-256-CBC after decoding hex")
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	// Remove PKCS7 padding
	paddingLen := int(plaintext[len(plaintext)-1])
	if paddingLen > len(plaintext) {
		return "", fmt.Errorf("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-paddingLen]
	decrypted := string(plaintext)
	return decrypted, nil
}
