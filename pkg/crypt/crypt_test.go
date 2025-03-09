// crypt_test.go
package crypt

import (
	"testing"
	"time"

	"gophKeeper/internal/client/services/lockbox/models"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "superSecretKey19"
	encryptor := New(key)
	plaintext := "Hello, World!"

	encrypted, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original. Got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecryptStruct(t *testing.T) {
	key := "superSecretKey19"
	encryptor := New(key)
	input := models.LockBoxInput{
		Description: "My secret description",
		Login:       "user123",
		URL:         "https://example.com",
		Password:    "passw0rd",
	}
	original := input

	err := EncryptStruct(&input, encryptor)
	if err != nil {
		t.Fatalf("EncryptStruct failed: %v", err)
	}

	if input.Description == original.Description {
		t.Error("Description was not encrypted")
	}
	if input.Login == original.Login {
		t.Error("Login was not encrypted")
	}
	if input.URL == original.URL {
		t.Error("URL was not encrypted")
	}
	if input.Password == original.Password {
		t.Error("Password was not encrypted")
	}

	err = DecryptStruct(&input, encryptor)
	if err != nil {
		t.Fatalf("DecryptStruct failed: %v", err)
	}

	if input.Description != original.Description {
		t.Errorf("After decryption, Description mismatch: got %q, want %q", input.Description, original.Description)
	}
	if input.Login != original.Login {
		t.Errorf("After decryption, Login mismatch: got %q, want %q", input.Login, original.Login)
	}
	if input.URL != original.URL {
		t.Errorf("After decryption, URL mismatch: got %q, want %q", input.URL, original.URL)
	}
	if input.Password != original.Password {
		t.Errorf("After decryption, Password mismatch: got %q, want %q", input.Password, original.Password)
	}
}

func TestEncryptDecryptLockBox(t *testing.T) {
	key := "superSecretKey19"
	encryptor := New(key)
	box := models.LockBox{
		Name:        "TestBox",
		Description: "Confidential info",
		Login:       "testuser",
		URL:         "https://test.com",
		Password:    "s3cr3t",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	original := box

	err := EncryptLockBox(&box, encryptor)
	if err != nil {
		t.Fatalf("EncryptLockBox failed: %v", err)
	}

	if box.Description == original.Description {
		t.Error("LockBox Description was not encrypted")
	}
	if box.URL == original.URL {
		t.Error("LockBox URL was not encrypted")
	}
	if box.Password == original.Password {
		t.Error("LockBox Password was not encrypted")
	}

	err = DecryptLockBox(&box, encryptor)
	if err != nil {
		t.Fatalf("DecryptLockBox failed: %v", err)
	}
	if box.Description != original.Description {
		t.Errorf("After decryption, Description mismatch: got %q, want %q", box.Description, original.Description)
	}
	if box.Login != original.Login {
		t.Errorf("After decryption, Login mismatch: got %q, want %q", box.Login, original.Login)
	}
	if box.URL != original.URL {
		t.Errorf("After decryption, URL mismatch: got %q, want %q", box.URL, original.URL)
	}
	if box.Password != original.Password {
		t.Errorf("After decryption, Password mismatch: got %q, want %q", box.Password, original.Password)
	}
}
