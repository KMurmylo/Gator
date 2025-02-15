package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {

	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}
	}
	dat, err := os.ReadFile(filepath.Join(home, configFileName))
	if err != nil {
		return Config{}
	}
	var result Config
	err = json.Unmarshal(dat, &result)
	if err != nil {
		return Config{}
	}
	return result

}
func (c *Config) SetUser(user string) error {

	c.UserName = user
	json, err := json.Marshal(c)
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(home, configFileName), json, 0644)
	if err != nil {
		return err
	}
	return nil

}
