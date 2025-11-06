// Copyright © 2024-2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"errors"
	"fmt"
	"io"
	"os"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// 暗号化キー（32バイト = AES-256）
// 本番環境では環境変数や設定ファイルから読み込むことを推奨
//
//	暗号化キー（32バイト = AES-256）を環境変数から64桁の16進数文字列として取得し、バイト列に変換する
var encryptionKey []byte

func init() {
	keyHex := os.Getenv("IP_ENCRYPTION_KEY_HEX")
	if len(keyHex) != 64 {
		panic("IP_ENCRYPTION_KEY_HEX must be a 64-character hexadecimal string")
	}
	encryptionKey = make([]byte, 32)
	for i := 0; i < 32; i++ {
		var byteVal byte
		_, err := fmt.Sscanf(keyHex[i*2:i*2+2], "%02x", &byteVal)
		if err != nil {
			panic("Invalid hexadecimal character in IP_ENCRYPTION_KEY_HEX")
		}
		encryptionKey[i] = byteVal
	}
}

// EncryptIP はIPアドレスを暗号化してBase64エンコードされた文字列として返す
func EncryptIP(ipAddress string) (string, error) {
	if ipAddress == "" {
		return "", nil
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(ipAddress), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptIP は暗号化されたIPアドレスを復号化する
func DecryptIP(encryptedIP string) (string, error) {
	if encryptedIP == "" {
		return "", nil
	}

	ciphertext, err := base64.URLEncoding.DecodeString(encryptedIP)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
