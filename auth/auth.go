package auth

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cli/utils"
	. "cli/structs"
)

var config Config
var keyFile string
var tokenFile string

func Auth(config Config) error {
  err := config.LoadConfig("config.yaml")
  if err != nil {
    panic("Error loading config: " + err.Error())
  }

	keyFile = filepath.Join(config.Headless.KeyFilePath, config.Headless.KeyFileName)
	tokenFile = filepath.Join(config.Headless.TokenFilePath, config.Headless.TokenFileName)

  if !utils.IsOsHeadless() && !config.Dev.ForceFallbackMethod {
    // use go-keyring here...
      panic("Error loading key: " + err.Error())
  	return nil
  } else {
    key, err := loadOrCreateKey()
    if err != nil {
      panic("Error loading key: " + err.Error())
    }

    credsExist := utils.FileExists(keyFile) && utils.FileExists(tokenFile);

    if (credsExist) {
      creds, err := loadEncryptedCredentials(key)
      if err != nil {
        panic("Error loading token: " + err.Error())
      }

      fmt.Println("Decrypted credentials: username", creds.Username)
      fmt.Println("Decrypted credentials: token", creds.Username)
    } else {
      reader := bufio.NewReader(os.Stdin)

      var username string;
      var token string;

      for {
        fmt.Println("Enter your name:")
        u, err := reader.ReadString('\n')

        if (err == nil) {
          username = strings.TrimSpace(u)
          break
        }
      }
      
      for {
        fmt.Println("Enter your token:")
        t, err := reader.ReadString('\n')

        if (err == nil) {
          token = strings.TrimSpace(t)
          break
        }
      }

      var creds Credentials = Credentials{username, token}
      
      err = saveEncryptedCredentials(creds, key)
      if err != nil {
        panic("Error saving token: " + err.Error())
      }
    }
	
  	return nil
  }
}

func loadOrCreateKey() ([]byte, error) {
	if data, err := os.ReadFile(keyFile); err == nil {
		return data, nil
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

  err := os.WriteFile(keyFile, key, 0600)
	return key, err
}

func saveEncryptedCredentials(creds Credentials, key []byte) error {
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
