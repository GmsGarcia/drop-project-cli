package main

import (
	"cli/auth"
	"cli/structs"
	"fmt"
)

var config structs.Config

func main() {
  err := config.LoadConfig("config.yaml")
  if err != nil {
    fmt.Println(err)
    return
  }

  err = auth.Auth(config);
  if (err != nil) {
    fmt.Println(err)
    return
  }

  fmt.Println("ending...")
}
