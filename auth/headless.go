package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
  . "cli/structs"
)

func loadOrCreateKey() ([]byte, error) {
	if data, err := os.ReadFile(keyFile); err == nil {
		return data, nil
	}

  dir := filepath.Dir(keyFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

  err := os.WriteFile(keyFile, key, 0600)
	return key, err
}

func saveEncryptedCredentials(creds Credentials, key []byte) error {
  dir := filepath.Dir(tokenFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

  data, err := creds.AsBytes()
  if err != nil {
    return err
  }

  if (!config.Dev.EnableTokenEncryption) {
	  return os.WriteFile(tokenFile, data, 0600)
  }

  block, err := aes.NewCipher(key)
  if err != nil {
    return err
  }

  gcm, err := cipher.NewGCM(block)
  if err != nil {
    return err
  }

  nonce := make([]byte, gcm.NonceSize())
  if _, err = rand.Read(nonce); err != nil {
    return err
  }

  ciphertext := gcm.Seal(nonce, nonce, data, nil)
  
  return os.WriteFile(tokenFile, ciphertext, 0600)
}

func loadEncryptedCredentials(key []byte) (Credentials, error) {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return Credentials{}, err
	}

  if (!config.Dev.EnableTokenEncryption) {
    var creds Credentials
    err = creds.FromBytes(data)

    if err != nil {
	    return Credentials{}, err
    }

	  return creds, nil
  }

	block, err := aes.NewCipher(key)
	if err != nil {
		return Credentials{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Credentials{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return Credentials{}, fmt.Errorf("data too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return Credentials{}, err
	}
  
  var creds Credentials = Credentials{}
  err = creds.FromBytes(plaintext);
  if err != nil {
		return Credentials{}, err
	}

	return creds, nil
}


func clearCredentials() error {
  return nil
}
