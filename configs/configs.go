package configs

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
)

type AppConfig struct {
	DB    DBConf     `yaml:"db"`
	Cache CacheConf  `yaml:"cache"`
	Users []UserConf `yaml:"users"`
}

func (ac *AppConfig) LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteBuff := new(bytes.Buffer)
	//defer byteBuff.Reset()

	_, err = byteBuff.ReadFrom(file)
	if err != nil {
		return err
	}
	err = ReadYamlContent(byteBuff.Bytes(), ac)
	if err != nil {
		return err
	}
	return nil
}

type DBConf struct {
	Addr string `yaml:"addr"`
}

type CacheConf struct {
	Addr string `yaml:"addr"`
}

type UserConf struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

func ReadYamlContent(content []byte, aim any) error {
	err := yaml.Unmarshal(content, aim)
	if err != nil {
		return err
	}
	return nil
}
