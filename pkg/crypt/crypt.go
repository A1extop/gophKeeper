package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gophKeeper/internal/client/services/lockbox/models"
	"io"
)

type Encryptor interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type AESCBCEncryptor struct {
	Key string
}

func New(key string) Encryptor {
	return &AESCBCEncryptor{Key: key}
}

func createHMAC(key, message []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// Проверка HMAC
func verifyHMAC(key, message, receivedHMAC []byte) bool {
	expectedHMAC := createHMAC(key, message)
	return hmac.Equal(expectedHMAC, receivedHMAC)
}

func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func unpad(data []byte) ([]byte, error) {
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, fmt.Errorf("ошибка: некорректный паддинг")
	}
	return data[:len(data)-padding], nil
}

func (e *AESCBCEncryptor) Encrypt(plaintext string) (string, error) {
	key := []byte(e.Key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("ошибка создания AES-шифра: %w", err)
	}

	blockSize := block.BlockSize()
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("ошибка генерации IV: %w", err)
	}

	plaintextBytes := pad([]byte(plaintext), blockSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintextBytes))
	mode.CryptBlocks(ciphertext, plaintextBytes)

	hmacValue := createHMAC(key, append(iv, ciphertext...))

	finalData := append(iv, ciphertext...)
	finalData = append(finalData, hmacValue...)

	return base64.StdEncoding.EncodeToString(finalData), nil
}

func (e *AESCBCEncryptor) Decrypt(encryptedText string) (string, error) {
	key := []byte(e.Key)

	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования Base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(data) < blockSize+sha256.Size {
		return "", fmt.Errorf("ошибка: повреждённые данные")
	}

	iv := data[:blockSize]
	hmacValue := data[len(data)-sha256.Size:]
	ciphertext := data[blockSize : len(data)-sha256.Size]

	if !verifyHMAC(key, append(iv, ciphertext...), hmacValue) {
		return "", fmt.Errorf("ошибка: HMAC не совпадает, возможна подмена данных")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = unpad(plaintext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func EncryptStruct(data *models.LockBoxInput, encryptor Encryptor) (*models.LockBoxInput, error) {
	encryptedData := &models.LockBoxInput{
		Description: data.Description,
		Login:       data.Login,
		URL:         data.URL,
		Password:    data.Password,
	}

	var err error
	if encryptedData.Description != "" {
		encryptedData.Description, err = encryptor.Encrypt(encryptedData.Description)
		if err != nil {
			return nil, err
		}
	}
	if encryptedData.Password != "" {
		encryptedData.Password, err = encryptor.Encrypt(encryptedData.Password)
		if err != nil {
			return nil, err
		}
	}
	if encryptedData.URL != "" {
		encryptedData.URL, err = encryptor.Encrypt(encryptedData.URL)
		if err != nil {
			return nil, err
		}
	}
	if encryptedData.Login != "" {
		encryptedData.Login, err = encryptor.Encrypt(encryptedData.Login)
		if err != nil {
			return nil, err
		}
	}
	encryptedData.Name = data.Name
	return encryptedData, nil
}

func DecryptStruct(data *models.LockBoxInput, encryptor Encryptor) (*models.LockBoxInput, error) {
	decryptedData := &models.LockBoxInput{
		Description: data.Description,
		Login:       data.Login,
		URL:         data.URL,
		Password:    data.Password,
	}

	var err error
	if decryptedData.Description != "" {
		decryptedData.Description, err = encryptor.Decrypt(decryptedData.Description)
		if err != nil {
			return nil, err
		}
	}
	if decryptedData.Password != "" {
		decryptedData.Password, err = encryptor.Decrypt(decryptedData.Password)
		if err != nil {
			return nil, err
		}
	}
	if decryptedData.URL != "" {
		decryptedData.URL, err = encryptor.Decrypt(decryptedData.URL)
		if err != nil {
			return nil, err
		}
	}
	if decryptedData.Login != "" {
		decryptedData.Login, err = encryptor.Decrypt(decryptedData.Login)
		if err != nil {
			return nil, err
		}
	}

	return decryptedData, nil
}

func DecryptLockBox(lockBox *models.LockBox, encryptor Encryptor) (*models.LockBox, error) {
	lockInput := models.LockBoxInput{
		Description: lockBox.Description,
		Login:       lockBox.Login,
		URL:         lockBox.URL,
		Password:    lockBox.Password,
	}

	decryptedInput, err := DecryptStruct(&lockInput, encryptor)
	if err != nil {
		return nil, err
	}

	decryptedLockBox := &models.LockBox{
		Name:        lockBox.Name,
		Description: decryptedInput.Description,
		Login:       decryptedInput.Login,
		URL:         decryptedInput.URL,
		Password:    decryptedInput.Password,
		CreatedAt:   lockBox.CreatedAt,
		UpdatedAt:   lockBox.UpdatedAt,
		SyncedAt:    lockBox.SyncedAt,
		DeletedAt:   lockBox.DeletedAt,
	}

	return decryptedLockBox, nil
}

func EncryptLockBox(lockBox *models.LockBox, encryptor Encryptor) (*models.LockBox, error) {
	lockInput := models.LockBoxInput{
		Description: lockBox.Description,
		Login:       lockBox.Login,
		URL:         lockBox.URL,
		Password:    lockBox.Password,
	}

	encryptedInput, err := EncryptStruct(&lockInput, encryptor)
	if err != nil {
		return nil, err
	}

	encryptedLockBox := &models.LockBox{
		Name:        lockBox.Name,
		Description: encryptedInput.Description,
		Login:       encryptedInput.Login,
		URL:         encryptedInput.URL,
		Password:    encryptedInput.Password,
		CreatedAt:   lockBox.CreatedAt,
		UpdatedAt:   lockBox.UpdatedAt,
		SyncedAt:    lockBox.SyncedAt,
		DeletedAt:   lockBox.DeletedAt,
	}
	return encryptedLockBox, nil
}
