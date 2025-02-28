package conf

import "gopkg.in/yaml.v3"   

func ReadYamlContent(content []byte, aim any) error {
	err := yaml.Unmarshal(content, aim)
	if err != nil {
		return err
	}
	return nil
}