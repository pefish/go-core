package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pefish/go-core/global-api-strategy"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-error"
	"github.com/pefish/go-format"
	"github.com/pefish/yaml"
	"io/ioutil"
	"reflect"
	"strings"
)

type SwaggerClass struct {

}

var swagger *SwaggerClass

func GetSwaggerInstance() *SwaggerClass {
	if swagger != nil {
		return swagger
	}
	swagger = &SwaggerClass{}
	return swagger
}

type Yaml_Info struct {
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
	Title       string `json:"title" yaml:"title"`
}

type Yaml_Tag struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

type Yaml_Parameter struct {
	In          string            `json:"in" yaml:"in"`
	Name        string            `json:"name" yaml:"name"`
	Required    bool              `json:"required" yaml:"required"`
	Description string            `json:"description" yaml:"description"`
	Type        string            `json:"type,omitempty" yaml:"type,omitempty"`
	Schema      map[string]string `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type Yaml_Response struct {
	Description string                 `json:"description" yaml:"description"`
	Schema      map[string]interface{} `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type Yaml_Path struct {
	Tags        []string                 `json:"tags" yaml:"tags"`
	Summary     string                   `json:"summary" yaml:"summary"`
	Consumes    []string                 `json:"consumes" yaml:"consumes"`
	Produces    []string                 `json:"produces" yaml:"produces"`
	Parameters  []Yaml_Parameter         `json:"parameters" yaml:"parameters"`
	Responses   map[string]Yaml_Response `json:"responses" yaml:"responses"`
	Description string                   `json:"description" yaml:"description"`
}

type Yaml_Property struct {
	Type        string                   `json:"type,omitempty" yaml:"type,omitempty"`
	Example     interface{}              `json:"example,omitempty" yaml:"example,omitempty"`
	Ref         string                   `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Items       map[string]interface{}   `json:"items,omitempty" yaml:"items,omitempty"`
	Description string                   `json:"description" yaml:"description"`
	Properties  map[string]Yaml_Property `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type Yaml_Definition struct {
	Type       string                   `json:"type,omitempty" yaml:"type,omitempty"`
	Required   []string                 `json:"required,omitempty" yaml:"required,omitempty"`
	Properties map[string]Yaml_Property `json:"properties" yaml:"properties"`
}

type Yaml_Swagger struct {
	Swagger     string                          `json:"swagger" yaml:"swagger"`
	Info        Yaml_Info                       `json:"info" yaml:"info"`
	Host        string                          `json:"host" yaml:"host"`
	BasePath    string                          `json:"basePath" yaml:"basePath"`
	Tags        []Yaml_Tag                      `json:"tags" yaml:"tags"`
	Schemes     []string                        `json:"schemes" yaml:"schemes"`
	Paths       map[string]map[string]Yaml_Path `json:"paths" yaml:"paths"`
	Definitions map[string]Yaml_Definition      `json:"definitions" yaml:"definitions"`
}

func (this *SwaggerClass) recuGetParams(paramsType reflect.Type, paramsVal reflect.Value, properties map[string]Yaml_Property, requiredParams *[]string, parameters *[]Yaml_Parameter) {
	if paramsType.Kind() == reflect.Ptr {
		paramsType = paramsType.Elem()
	}
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		realParamName := field.Tag.Get(`json`)
		properties[realParamName] = Yaml_Property{
			Type:        this.getType(field.Type.Name()),
			Example:     paramsVal.Interface(),
			Description: field.Tag.Get(`desc`),
		}
		parameter := Yaml_Parameter{
			Name:        realParamName,
			In:          `query`,
			Required:    strings.Contains(field.Tag.Get(`validate`), `required`),
			Description: field.Tag.Get(`desc`),
			Type:        this.getType(field.Type.Name()),
		}
		*parameters = append(*parameters, parameter)

		if strings.Contains(field.Tag.Get(`validate`), `required`) {
			*requiredParams = append(*requiredParams, realParamName)
		}
	}
}

func (this *SwaggerClass) recuReturn(map_ map[string]interface{}, properties map[string]Yaml_Property) {
	for k, v := range map_ {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			p := map[string]Yaml_Property{}
			this.recuReturn(v.(map[string]interface{}), p)
			properties[k] = Yaml_Property{
				Description: `no desc`,
				Properties:  p,
			}
		} else {
			properties[k] = Yaml_Property{
				Example:     v,
				Description: `no desc`,
			}
		}
	}
}

func (this *SwaggerClass) recuPostParams(paramsType reflect.Type, paramsVal reflect.Value, properties map[string]Yaml_Property, requiredParams *[]string) {
	if paramsType.Kind() == reflect.Ptr {
		paramsType = paramsType.Elem()
	}
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		fieldVal := paramsVal.Field(i)
		if field.Type.String() == `os.File` {
			realParamName := field.Tag.Get(`json`)
			properties[realParamName] = Yaml_Property{
				Type:        `file`,
				Example:     fieldVal.Interface(),
				Description: field.Tag.Get(`desc`),
			}
			if strings.Contains(field.Tag.Get(`validate`), `required`) {
				*requiredParams = append(*requiredParams, realParamName)
			}
		} else if field.Type.Kind() == reflect.Struct {
			this.recuPostParams(field.Type, fieldVal, properties, requiredParams)
		} else {
			realParamName := field.Tag.Get(`json`)
			properties[realParamName] = Yaml_Property{
				Type:        this.getType(field.Type.Name()),
				Example:     fieldVal.Interface(),
				Description: field.Tag.Get(`desc`),
			}
			if strings.Contains(field.Tag.Get(`validate`), `required`) {
				*requiredParams = append(*requiredParams, realParamName)
			}
		}
	}
}

func (this *SwaggerClass) getType(typeName string) string {
	result := ``
	if typeName == `int` ||
		typeName == `int8` ||
		typeName == `int16` ||
		typeName == `int32` ||
		typeName == `int64` ||
		typeName == `uint8` ||
		typeName == `uint16` ||
		typeName == `uint32` ||
		typeName == `uint64` {
		result = `number`
	} else if typeName == `bool` {
		result = `boolean`
	} else {
		result = typeName
	}
	return result
}

func (this *SwaggerClass) GeneSwagger(hostAndPort string, filename string, type_ string) {
	definitions := map[string]Yaml_Definition{}

	paths := map[string]map[string]Yaml_Path{}

	for _, api := range service.Service.GetApis() {
		temp := map[string]Yaml_Path{}

		desc := api.Description

		parameters := []Yaml_Parameter{}

		description := ``
		if api.Strategies != nil {
			for _, strategy := range api.Strategies {
				if strategy.Disable == false && strategy.Strategy.GetName() == `jwtAuth` {
					// 添加 jwt header
					parameters = append(parameters, Yaml_Parameter{
						Name:        `Json-Web-Token`,
						In:          `header`,
						Required:    true,
						Description: `jwt token`,
						Type:        this.getType(`string`),
					})
				}
				description += strategy.Strategy.GetName() + ": " + strategy.Strategy.GetDescription() + "\n"
			}
		}

		if api.Params != nil {
			paramsType := reflect.TypeOf(api.Params)
			if paramsType.Kind() == reflect.Ptr {
				paramsType = paramsType.Elem()
			}
			paramsTypeName := paramsType.Name()
			requiredParams := []string{}
			// 解析 properties
			properties := map[string]Yaml_Property{}
			if api.Method == `POST` {
				this.recuPostParams(paramsType, reflect.ValueOf(api.Params), properties, &requiredParams)
				parameter := Yaml_Parameter{
					In:       `body`,
					Name:     `body`,
					Required: true,
					Schema: map[string]string{
						`$ref`: fmt.Sprintf(`#/definitions/%s`, paramsTypeName),
					},
				}
				parameters = append(parameters, parameter)
			} else if api.Method == `GET` {
				this.recuGetParams(paramsType, reflect.ValueOf(api.Params), properties, &requiredParams, &parameters)
			} else {
				go_error.ThrowInternal(errors.New(`method error`))
			}
			definitions[paramsTypeName] = Yaml_Definition{
				Type:       `object`,
				Properties: properties,
				Required:   requiredParams,
			}
		}

		responses := map[string]Yaml_Response{}
		if api.Return != nil {
			type_ := reflect.TypeOf(api.Return)
			returnTypeName := type_.Name()
			kind := type_.Kind()
			properties := map[string]Yaml_Property{}
			if kind == reflect.Struct {
				this.recuReturn(go_format.Format.StructToMap(api.Return), properties)
			} else {
				go_error.ThrowInternal(errors.New(`return config type error`))
			}
			definitions[api.Path+`_`+returnTypeName] = Yaml_Definition{
				Type:       `object`,
				Properties: properties,
			}
			responses[`200`] = Yaml_Response{
				Description: `正确返回`,
				Schema: map[string]interface{}{
					`$ref`: fmt.Sprintf(`#/definitions/%s`, api.Path+`_`+returnTypeName),
				},
			}
		}

		paramTypes := []string{}
		if api.ParamType == global_api_strategy.ALL_TYPE {
			paramTypes = append(paramTypes, `application/json`, `multipart/form-data`)
		} else {
			paramTypes = append(paramTypes, api.ParamType)
		}

		temp[strings.ToLower(string(api.Method))] = Yaml_Path{
			Tags:        []string{service.Service.GetName()},
			Summary:     desc,
			Consumes:    paramTypes,
			Produces:    []string{`application/json`},
			Parameters:  parameters,
			Responses:   responses,
			Description: description,
		}
		paths[service.Service.GetPath()+api.Path] = temp
	}

	swagger := Yaml_Swagger{
		`2.0`,
		Yaml_Info{
			Title:       service.Service.GetName(),
			Description: service.Service.GetDescription(),
			Version:     `1.0.0`,
		},
		hostAndPort,
		service.Service.GetPath(),
		[]Yaml_Tag{
			{
				Name:        service.Service.GetName(),
				Description: service.Service.GetDescription(),
			},
		},
		[]string{`http`},
		paths,
		definitions,
	}

	if type_ == `yaml` {
		bytes, _ := yaml.Marshal(&swagger)
		err := ioutil.WriteFile(filename, bytes, 0777)
		if err != nil {
			panic(err)
		}
	} else if type_ == `json` {
		bytes, _ := json.Marshal(&swagger)
		err := ioutil.WriteFile(filename, bytes, 0777)
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New(`type 指定有误`))
	}

}
