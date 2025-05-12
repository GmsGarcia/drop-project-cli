package structs

import (
	"os"
  "strings"
	"gopkg.in/yaml.v3"
  "cli/utils"
)

type Config struct {
  Headless HeadlessConfig `yaml:"headless"`
  Api ApiConfig `yaml:"api"`
  Dev DevConfig `yaml:"dev"` // TODO: remove this... 
}

type HeadlessConfig struct {
  KeyFilePath string `yaml:"keyFilePath"`
  KeyFileName string `yaml:"keyFileName"`
  TokenFilePath string `yaml:"tokenFilePath"`
  TokenFileName string `yaml:"tokenFileName"`
}

type DevConfig struct {
  ForceFallbackMethod bool `yaml:"forceFallbackMethod"`
  EnableTokenEncryption bool `yaml:"enableTokenEncryption"` 
}

type ApiConfig struct {
  Server string `yaml:"server"`
  Endpoints EndpointsConfig `yaml:"endpoints"`
}

type EndpointsConfig struct {
  Assignments string `yaml:"assignments"`
  CurrentAssignment string `yaml:"current_assignment"`
  Submissions string `yaml:"submissions"`
  NewSubmission string `yaml:"new_submission"`
}

func (c *Config) LoadConfig(file string) {
  data, err := os.ReadFile(file)
  if err != nil {
    panic("Failed to load config: " + err.Error())
  }
  err = yaml.Unmarshal(data, &c)
  if err != nil {
    panic("Failed to load config: " + err.Error())
  }

  if utils.IsOsHeadless() {
    if strings.Contains(c.Headless.KeyFilePath, "{USER}") {
      username := utils.GetUsername() 
      c.Headless.KeyFilePath = strings.Replace(c.Headless.KeyFilePath, "{USER}", username, -1);
    }

    if strings.Contains(c.Headless.TokenFilePath, "{USER}") {
      username:= utils.GetUsername() 
      c.Headless.TokenFilePath = strings.Replace(c.Headless.TokenFilePath, "{USER}", username, -1);
    }
  }
}
