package util

import (
	"fmt"
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
	"os"
)

func Save2Yaml(keyData interface{}, filePath string) error {
	yamlBytes, err := yaml.Marshal(keyData)
	if err != nil {
		log.Error("YAML marshal error:", err)
		return err
	}
	// Write YAML to file
	err = os.WriteFile(filePath, yamlBytes, 0644)
	if err != nil {
		log.Error("Write file error:", err)
		return err
	}
	log.Info("YAML data has been written successfully.")
	return nil
}

func ReadAccountInfo(filePath string) (*Response, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Error(fmt.Sprintf("%s Read Fail %E", filePath, err))
		return nil, err
	}
	var data Response
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if &data == nil {
		err = fmt.Errorf("an error occurred while reading account information")
		log.Error(err)
		return nil, err
	}
	return &data, nil
}

func ReadKey(filePath string) (*Key, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Error(fmt.Sprintf("%s Read Fail %E", filePath, err))
		return nil, err
	}
	var data Key
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &data, nil
}
