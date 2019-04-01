package config

import (
	"errors"
	"github.com/pefish/go-file"
	"github.com/pefish/go-json"
	"github.com/pefish/go-map"
	"os"
)

type ConfigClass struct {
	configs map[string]interface{}
}

var Config = ConfigClass{}

func (this *ConfigClass) LoadConfig(configFilePtr *string, secretFilePtr *string) {
	configFile := ``
	configMap := map[string]interface{}{}
	if configFilePtr == nil {
		configFile = os.Getenv(`GO_CONFIG`)
	} else {
		configFile = *configFilePtr
	}
	if configFile != `` {
		configMap = p_json.Json.ParseBytes(p_file.File.ReadFile(configFile)).(map[string]interface{})
	}

	secretFile := ``
	secretMap := map[string]interface{}{}
	if secretFilePtr == nil {
		secretFile = os.Getenv(`GO_SECRET`)
	} else {
		secretFile = *secretFilePtr
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
