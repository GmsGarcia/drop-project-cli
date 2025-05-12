package main

import (
	"cli/auth"
	. "cli/structs"
	"flag"
	"fmt"
)

var config Config

func main() {
  signInFlag := flag.Bool("sign-in", false, "sign-in")
  signOutFlag := flag.Bool("sign-out", false, "sign-out")
  flag.Parse()

  config.LoadConfig("config.yaml")

  auth.Prepare(config)
  authed := auth.IsAuthed()

  if *signInFlag {
    if authed {
      fmt.Println("Already signed in.")
      return
    }

    err := auth.SignIn();
    if err != nil {
      fmt.Println("Failed to sign in: " + err.Error())
    }
    fmt.Println("Successfully signed in.")
    return
  }

  if *signOutFlag {
    if !authed {
      fmt.Println("Already signed out.")
      return
    }

    err := auth.SignOut();
    if err != nil {
      fmt.Println("Failed to sign out: " + err.Error())
    }
    fmt.Println("Successfully signed out.")
    return
  } 

  fmt.Println("end")
}
