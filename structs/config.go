package structs

import (
	"os"
	"os/user"
  "strings"
	"gopkg.in/yaml.v3"
)

type Config struct {
  Headed HeadedConfig `yaml:"headed"`
  Headless HeadlessConfig `yaml:"headless"`
  Dev DevConfig `yaml:"dev"`
}

type HeadedConfig struct {
  Todo string `yaml:"todo"`
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

func (c *Config) LoadConfig(file string) error {
  data, err := os.ReadFile(file)
  if err != nil {
    return err
  }
  err = yaml.Unmarshal(data, &c)
  if err != nil {
    return err
  }

  if strings.Contains(c.Headless.KeyFilePath, "{USER}") {
    user, err := user.Current()
    if err != nil {
      return err
    }

    c.Headless.KeyFilePath = strings.Replace(c.Headless.KeyFilePath, "{USER}", user.Username, -1);
  }

  if strings.Contains(c.Headless.TokenFilePath, "{USER}") {
    user, err := user.Current()
    if err != nil {
      return err
    }

    c.Headless.TokenFilePath = strings.Replace(c.Headless.TokenFilePath, "{USER}", user.Username, -1);
  }

  return nil
}
