package structs

import "encoding/json"

type Credentials struct {
  Username string `json:"username"`
  Token string `json:"token"`
}

func (c *Credentials) FromBytes(data []byte) error {
  return json.Unmarshal(data, &c);
}

func (c Credentials) AsBytes() ([]byte, error) {
  return json.Marshal(c)
}
