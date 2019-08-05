package config

import (
	"errors"
	"github.com/pefish/go-file"
	"github.com/pefish/go-json"
	"github.com/pefish/go-map"
	"github.com/pefish/go-reflect"
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
		bytes, err := go_file.File.ReadFileWithErr(configFile)
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
		bytes, err := go_file.File.ReadFileWithErr(secretFile)
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
	this.configs = go_map.Map.Append(configMap, secretMap)
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
		configMap = go_json.Json.ParseBytes(go_file.File.ReadFile(configFile)).(map[string]interface{})
	}

	secretFile := ``
	secretMap := map[string]interface{}{}
	if config.SecretFilepath == `` {
		secretFile = os.Getenv(`GO_SECRET`)
	} else {
		secretFile = config.SecretFilepath
	}
	if secretFile != `` {
		bytes, err := go_file.File.ReadFileWithErr(secretFile)
		if err == nil {
			secretMap = go_json.Json.ParseBytes(bytes).(map[string]interface{})
		}
	}

	if configFile == `` && secretFile == `` {
		panic(errors.New(`unspecified config file and secret file`))
	}
	this.configs = go_map.Map.Append(configMap, secretMap)
}

func (this *ConfigClass) GetString(str string) string {
	return go_reflect.Reflect.ToString(this.configs[str])
}

func (this *ConfigClass) GetInt(str string) int {
	return go_reflect.Reflect.ToInt(this.configs[str])
}

func (this *ConfigClass) GetInt64(str string) int64 {
	return go_reflect.Reflect.ToInt64(this.configs[str])
}

func (this *ConfigClass) GetUint64(str string) uint64 {
	return go_reflect.Reflect.ToUint64(this.configs[str])
}

func (this *ConfigClass) GetBool(str string) bool {
	return go_reflect.Reflect.ToBool(this.configs[str])
}

func (this *ConfigClass) GetFloat64(str string) float64 {
	return go_reflect.Reflect.ToFloat64(this.configs[str])
}

func (this *ConfigClass) Get(str string) interface{} {
	return this.configs[str]
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
