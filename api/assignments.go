package api

import (
	"cli/structs"
	"net/http"
)

// WIP
func GetCurrentAssignment(config structs.Config, credentials structs.Credentials) (int, error) {
  req, err := http.NewRequest("GET", "https://" + config.Api.Server + config.Api.Endpoints.CurrentAssignment, nil)
  if err != nil {
    return -1, err
  }

  req.SetBasicAuth(credentials.Username, credentials.Token)

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    return -1, err
  }
  defer res.Body.Close()

  return res.StatusCode, nil
}

// WIP
func GetAssignmentById(config structs.Config, credentials structs.Credentials, id string) (int, error) {
  req, err := http.NewRequest("GET", "https://" + config.Api.Server + config.Api.Endpoints.Assignments + "/" + id, nil)
  if err != nil {
    return -1, err
  }

  req.SetBasicAuth(credentials.Username, credentials.Token)

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    return -1, err
  }
  defer res.Body.Close()

  return res.StatusCode, nil
}
