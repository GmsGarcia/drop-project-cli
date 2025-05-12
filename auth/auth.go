package auth

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
  "syscall"
	"github.com/zalando/go-keyring"
  "golang.org/x/term"
	"cli/structs"
	"cli/utils"
  "net/http"
)

var config structs.Config
var keyFile string
var tokenFile string

func Prepare(cfg structs.Config) {
  config = cfg
	keyFile = filepath.Join(config.Headless.KeyFilePath, config.Headless.KeyFileName)
	tokenFile = filepath.Join(config.Headless.TokenFilePath, config.Headless.TokenFileName)
}

func SignIn() error {
  authed := IsAuthed()
  
  if authed {
    return nil
  }

  if !utils.IsOsHeadless() && !config.Dev.ForceFallbackMethod {
    username := utils.GetUsername()

    for {
      var credentials structs.Credentials = promptCredentials()

      bytes, err := credentials.AsBytes()
      if err != nil {
        return err
      } 

      valid, err := checkCredentials(credentials)
      if err != nil {
        return err
      }

      if valid {
        keyring.Set("dp-cli", username, string(bytes))
        break
      } else {
        fmt.Println("Invalid credentials.")
      }
    }
  } else {
    key, err := loadOrCreateKey()
    if err != nil {
      return err
    }

    for {
      var credentials structs.Credentials = promptCredentials()
    
      valid, err := checkCredentials(credentials)
      if err != nil {
        return err
      }

      if valid {
        err = saveEncryptedCredentials(credentials, key)
        if err != nil {
          return err
        }
        break
      } else {
        fmt.Println("Invalid credentials.")
      }
    }
  }

  return nil
}

func SignOut() error {
  authed := IsAuthed()

  if !authed {
    return nil
  }

  if !utils.IsOsHeadless() && !config.Dev.ForceFallbackMethod {
    username := utils.GetUsername()

    err := keyring.Delete("dp-cli", username)
    if err != nil {
      return fmt.Errorf("%s", "Failed to clear credentials: " + err.Error())
    }
  } else {
    err := clearCredentials()
    if err != nil {
      return fmt.Errorf("%s", "Failed to clear credentials: " + err.Error())
    }
  }
  return nil
}

func GetCredentials() *structs.Credentials {
  authed := IsAuthed()

  if !authed {
    return nil
  }

  if !utils.IsOsHeadless() && !config.Dev.ForceFallbackMethod {
    username := utils.GetUsername()
    data, err := keyring.Get("dp-cli", username)
    
    var credentials structs.Credentials
    
    err = credentials.FromBytes([]byte(data))
    if err != nil {
      panic("Failed to load credentials: " + err.Error())
    }

  	return &credentials
  } else {
    key, err := loadOrCreateKey()
    if err != nil {
      panic("Failed to load key: " + err.Error())
    }

    credentials, err := loadEncryptedCredentials(key)
    if err != nil {
      panic("Failed to load credentials: " + err.Error())
    }

  	return &credentials
  }
}

func IsAuthed() bool {
  if !utils.IsOsHeadless() && !config.Dev.ForceFallbackMethod {
    username:= utils.GetUsername()

    _, err := keyring.Get("dp-cli", username)
    if err != nil {
      if err != keyring.ErrNotFound {
        panic("Failed to check auth: " + err.Error())
      }

      return false
    } 

    return true
  } 

  return utils.FileExists(keyFile) && utils.FileExists(tokenFile);
}

func promptCredentials() structs.Credentials {
  reader := bufio.NewReader(os.Stdin)

  var username string;
  var token string;

  for {
    fmt.Print("Enter your name: ")
    u, err := reader.ReadString('\n')

    if err == nil {
      username = strings.TrimSpace(u)
      break
    }
  }
  
  for {
    fmt.Print("Enter your token: ")
    t, err := term.ReadPassword(int(syscall.Stdin))
    fmt.Println(" ")
    if err == nil {
      token = strings.TrimSpace(string(t))
      break
    }
  }

  return structs.Credentials{Username: username, Token: token}
}

func checkCredentials(credentials structs.Credentials) (bool, error) {
  req, err := http.NewRequest("GET", "https://" + config.Api.Server + config.Api.Endpoints.CurrentAssignment, nil)
  if err != nil {
    return false, err
  }

  req.SetBasicAuth(credentials.Username, credentials.Token)

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    return false, err
  }
  defer res.Body.Close()

  if res.StatusCode == 200 {
    return true, nil
  } else if res.StatusCode == 401 {
    return false, nil
  } else {
    return false, fmt.Errorf("Received %d from the server. Expected 200 or 401", res.StatusCode)
  }
}
