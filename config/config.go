package config

import (
	"errors"
	"github.com/pefish/go-file"
	"github.com/pefish/go-json"
	"github.com/pefish/go-map"
	"gopkg.in/yaml.v2"
	"os"
)

type ConfigClass struct {
	configs map[string]interface{}
}

var Config = ConfigClass{}

type Configuration struct {
	ConfigFilepath string
	SecretFilepath string
}

func (this *ConfigClass) LoadYamlConfig(config Configuration) {
	configFile := ``
	configMap := map[string]interface{}{}
	if config.ConfigFilepath == `` {
		configFile = os.Getenv(`GO_CONFIG`)
	} else {
		configFile = config.ConfigFilepath
	}
	if configFile != `` {
		bytes, err := p_file.File.ReadFileWithErr(configFile)
		if err == nil {
			err = yaml.Unmarshal(bytes, &configMap)
			if err != nil {
				panic(err)
			}
		}
	}

	secretFile := ``
	secretMap := map[string]interface{}{}
	if config.SecretFilepath == `` {
		secretFile = os.Getenv(`GO_SECRET`)
	} else {
		secretFile = config.SecretFilepath
	}
	if secretFile != `` {
		bytes, err := p_file.File.ReadFileWithErr(secretFile)
		if err == nil {
			err = yaml.Unmarshal(bytes, &secretMap)
			if err != nil {
				panic(err)
			}
		}
	}

	if configFile == `` && secretFile == `` {
		panic(errors.New(`unspecified config file and secret file`))
	}
	this.configs = p_map.Map.Append(configMap, secretMap)
}

func (this *ConfigClass) LoadJsonConfig(config Configuration) {
	configFile := ``
	configMap := map[string]interface{}{}
	if config.ConfigFilepath == `` {
		configFile = os.Getenv(`GO_CONFIG`)
	} else {
		configFile = config.ConfigFilepath
	}
	if configFile != `` {
		configMap = p_json.Json.ParseBytes(p_file.File.ReadFile(configFile)).(map[string]interface{})
	}

	secretFile := ``
	secretMap := map[string]interface{}{}
	if config.SecretFilepath == `` {
		secretFile = os.Getenv(`GO_SECRET`)
	} else {
		secretFile = config.SecretFilepath
	}
	if secretFile != `` {
		bytes, err := p_file.File.ReadFileWithErr(secretFile)
		if err == nil {
			secretMap = p_json.Json.ParseBytes(bytes).(map[string]interface{})
		}
	}

	if configFile == `` && secretFile == `` {
		panic(errors.New(`unspecified config file and secret file`))
	}
	this.configs = p_map.Map.Append(configMap, secretMap)
}

func (this *ConfigClass) GetString(str string) string {
	return this.configs[str].(string)
}

func (this *ConfigClass) GetInt64(str string) int64 {
	return int64(this.configs[str].(float64))
}

func (this *ConfigClass) GetBool(str string) bool {
	return this.configs[str].(bool)
}

func (this *ConfigClass) GetFloat64(str string) float64 {
	return this.configs[str].(float64)
}

func (this *ConfigClass) GetMap(str string) map[string]interface{} {
	return this.configs[str].(map[string]interface{})
}

func (this *ConfigClass) GetSlice(str string) []interface{} {
	return this.configs[str].([]interface{})
}

func (this *ConfigClass) GetSliceString(str string) []string {
	return this.configs[str].([]string)
}

func (this *ConfigClass) GetSliceWithErr(str string) ([]interface{}, error) {
	result, ok := this.configs[str].([]interface{})
	if !ok {
		return nil, errors.New(`type assert error`)
	}
	return result, nil
}
