package structs

import "encoding/json"

type Assignment struct {
  Id string `json:"id"`
  Name string `json:"name"`
  PackageName string `json:"packageName"`
  DueDate string `json:"dueDate"`
  SubmissionMethod string `json:"submissionMethod"`
  Language string `json:"language"`
  Active string `json:"active"`
  Intructions AssignmentInstructions `json:"instructions"`
}

type AssignmentInstructions struct {
  Format string `json:"format"`
  Body string `json:"body"`
}

func (a *Assignment) FromBytes(data []byte) error {
  return json.Unmarshal(data, &a);
}
