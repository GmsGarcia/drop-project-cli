package utils

import (
	"os"
	"os/user"
	"runtime"
)

func IsOsHeadless() bool {
  switch runtime.GOOS {
  case "windows":
    return false 
  case "darwin":
    return false 
  case "linux":
    return os.Getenv("DISPLAY") == ""
  default:
    return true // unknown OS, assume headless
  }
}

func FileExists(path string) bool {
  _, err := os.ReadFile(path);

  if err != nil {
    return false;
  } 
  return true;
}

func GetUsername() string {
  usr, err := user.Current()
  if err != nil {
    panic("Failed to get user: " + err.Error())
  }

  if runtime.GOOS == "linux" && usr.Username == "root" {
    name, found := os.LookupEnv("SUDO_USER")
    if !found {
      panic("Failed to get user: command executed by root but no 'SUDO_USER' variable is set")
    }

    if _, err := user.LookupId(name); err != nil {
      panic("Failed to get user: command executed by root but 'SUDO_USER' is not a valid user")
    }

    return name
  }

  return usr.Username
}
