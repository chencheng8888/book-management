package configs

import "os"


type Configs interface {
	ReadConfig(data []byte) error
}



func LoadConfigs(filePath string,conf ...Configs) error{
	data,err := os.ReadFile(filePath)
	if err!= nil {
		return err
	}
	for _,c := range conf {
		err = c.ReadConfig(data)
		if err != nil {
			return err
		}
	}
	return nil
}