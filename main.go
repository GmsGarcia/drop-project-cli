package main

import (
	"cli/auth"
	"fmt"
)

func main() {
  _, err := auth.Auth();
  if (err != nil) {
    fmt.Println(err)
    return
  }

  fmt.Println("ending...")
}
