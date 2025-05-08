package utils

import (
  "os"
)

func IsOsHeadless() bool {
  if os.Getenv("DISPLAY") == "" {
    return true
  } else {
    return false
  }
}

func FileExists(path string) bool {
  _, err := os.ReadFile(path);

  if err != nil {
    return false;
  } 
  return true;
}
