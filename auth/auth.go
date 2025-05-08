package auth

import (
	// "bufio"
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
  keyFilePath = "/etc/dp-cli/"
	keyFileName = "dp-cli.key"

	tokenFilePath = "token"
	tokenFileName = "dp-cli.token"

  appDirName  = ".dp-cli" // TODO: remove this
	tokenFile   = "token"

  forceFallbackMethod = true // disable this in prod to enable keyring storage method
  enableTokenEncryption = true // dev stuff! set this to true in prod
)

func Auth() error {
  if !isOsHeadless() && !forceFallbackMethod {
    // use go-keyring here...
  	return nil
  } else {
    key, err := loadOrCreateKey()
    if err != nil {
      panic("Error loading key: " + err.Error())
    }

    credsExist := areCredentialsStored();

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

func isOsHeadless() bool {
  if os.Getenv("DISPLAY") == "" {
    return true
  } else {
    return false
  }
}

func getAppDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, appDirName)
	os.MkdirAll(appDir, 0700)
	return appDir, nil
}

func loadOrCreateKey() ([]byte, error) {
	appDir, err := getAppDir()
	if err != nil {
		return nil, err
	}

	keyPath := filepath.Join(appDir, keyFileName)

	if data, err := os.ReadFile(keyPath); err == nil {
		return data, nil
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	err = os.WriteFile(keyPath, key, 0600)
	return key, err
}

func saveEncryptedCredentials(creds Credentials, key []byte) error {
  data, err := creds.AsBytes()
  if err != nil {
    return err
  }

  if (!enableTokenEncryption) {
	  appDir, _ := getAppDir()
	  return os.WriteFile(filepath.Join(appDir, tokenFile), data, 0600)
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
  appDir, _ := getAppDir()
  return os.WriteFile(filepath.Join(appDir, tokenFile), ciphertext, 0600)
}

func loadEncryptedCredentials(key []byte) (Credentials, error) {
	appDir, _ := getAppDir()
	data, err := os.ReadFile(filepath.Join(appDir, tokenFile))
	if err != nil {
		return Credentials{}, err
	}

  if (!enableTokenEncryption) {
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

func areCredentialsStored() bool {
	appDir, _ := getAppDir()

  _, errK := os.ReadFile(filepath.Join(appDir, keyFileName));
  _, errT := os.ReadFile(filepath.Join(appDir, tokenFile));

  if (errK != nil || errT != nil) {
    return false;
  } 
  return true;
}
