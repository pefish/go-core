package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-error"
	"github.com/pefish/go-file"
	"github.com/pefish/yaml"
	"reflect"
	"strings"
)

type SwaggerClass struct {
	service service.InterfaceService
}

var swagger *SwaggerClass

func GetSwaggerInstance() *SwaggerClass {
	if swagger != nil {
		return swagger
	}
	swagger = &SwaggerClass{}
	return swagger
}

func (this *SwaggerClass) SetService(service service.InterfaceService) *SwaggerClass {
	this.service = service
	return this
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

func (this *SwaggerClass) recuGetParams(paramsType reflect.Type, properties map[string]Yaml_Property, requiredParams *[]string, parameters *[]Yaml_Parameter) {
	if paramsType.Kind() == reflect.Ptr {
		paramsType = paramsType.Elem()
	}
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		realParamName := field.Tag.Get(`json`)
		properties[realParamName] = Yaml_Property{
			Type:        this.getType(field.Type.Name()),
			Example:     field.Tag.Get(`example`),
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

func (this *SwaggerClass) recuPostParams(paramsType reflect.Type, properties map[string]Yaml_Property, requiredParams *[]string) {
	if paramsType.Kind() == reflect.Ptr {
		paramsType = paramsType.Elem()
	}
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		if field.Type.String() == `os.File` {
			realParamName := field.Tag.Get(`json`)
			properties[realParamName] = Yaml_Property{
				Type:        `file`,
				Example:     field.Tag.Get(`example`),
				Description: field.Tag.Get(`desc`),
			}
			if strings.Contains(field.Tag.Get(`validate`), `required`) {
				*requiredParams = append(*requiredParams, realParamName)
			}
		} else if field.Type.Kind() == reflect.Struct {
			this.recuPostParams(field.Type, properties, requiredParams)
		} else {
			realParamName := field.Tag.Get(`json`)
			properties[realParamName] = Yaml_Property{
				Type:        this.getType(field.Type.Name()),
				Example:     field.Tag.Get(`example`),
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
	if typeName == `int64` {
		result = `integer`
	} else if typeName == `bool` {
		result = `boolean`
	} else {
		result = typeName
	}
	return result
}

func (this *SwaggerClass) GeneSwagger(hostAndPort string, filename string, type_ string) {
	definitions := map[string]Yaml_Definition{}
	definitions[`Failed`] = Yaml_Definition{
		Type: `object`,
		Properties: map[string]Yaml_Property{
			`succeed`: {
				Type:        this.getType(`bool`),
				Example:     `false`,
				Description: `表示请求是否成功`,
			},
			`error_code`: {
				Type:        this.getType(`int64`),
				Example:     `2000`,
				Description: `错误码`,
			},
			`error_message`: {
				Type:        this.getType(`string`),
				Example:     `jwt verify error`,
				Description: `错误原因（开发者看的）`,
			},
			`show_message`: {
				Type:        this.getType(`string`),
				Example:     `登录失败`,
				Description: `错误原因（用户看的）`,
			},
		},
	}

	paths := map[string]map[string]Yaml_Path{}

	for key, route := range this.service.GetRoutes() {
		temp := map[string]Yaml_Path{}

		desc := route.Description

		parameters := []Yaml_Parameter{}

		// 添加 lang header
		//parameters = append(parameters, Yaml_Parameter{
		//	Name: `lang`,
		//	In: `header`,
		//	Required: false,
		//	Description: `客户端语言`,
		//	Type: this.getType(`string`),
		//})

		description := ``
		if route.Strategies != nil {
			for _, strategy := range route.Strategies {
				if strategy[0].(string) == `jwt_auth` {
					// 添加 jwt header
					parameters = append(parameters, Yaml_Parameter{
						Name:        `Json-Web-Token`,
						In:          `header`,
						Required:    true,
						Description: `jwt token`,
						Type:        this.getType(`string`),
					})
				}
				description += strategy[0].(string) + `     `
			}
		}

		if route.Params != nil {
			paramsType := reflect.TypeOf(route.Params)
			if paramsType.Kind() == reflect.Ptr {
				paramsType = paramsType.Elem()
			}
			paramsTypeName := paramsType.Name()
			requiredParams := []string{}
			// 解析 properties
			properties := map[string]Yaml_Property{}
			if route.Method == `POST` {
				this.recuPostParams(paramsType, properties, &requiredParams)
				parameter := Yaml_Parameter{
					In:       `body`,
					Name:     `body`,
					Required: true,
					Schema: map[string]string{
						`$ref`: fmt.Sprintf(`#/definitions/%s`, paramsTypeName),
					},
				}
				parameters = append(parameters, parameter)
			} else if route.Method == `GET` {
				this.recuGetParams(paramsType, properties, &requiredParams, &parameters)
			} else {
				p_error.Throw(`method error`, 0)
			}
			definitions[paramsTypeName] = Yaml_Definition{
				Type:       `object`,
				Properties: properties,
				Required:   requiredParams,
			}
		}

		succeedDefinitionName := fmt.Sprintf(`%sSucceed`, key)
		succeedData := Yaml_Property{}

		if route.Return != nil {
			kind := reflect.TypeOf(route.Return).Kind()
			if kind == reflect.Map {
				properties := map[string]Yaml_Property{}
				for key, val1 := range route.Return.(map[string]map[string]interface{}) {
					properties[key] = Yaml_Property{
						Description: val1[`desc`].(string),
						Example:     val1[`example`],
					}
				}
				succeedData.Properties = properties
			} else if kind == reflect.Slice {
				properties := map[string]Yaml_Property{}
				eleType := reflect.TypeOf(route.Return).Elem()
				eleKind := eleType.Kind()
				if eleKind == reflect.Struct {
					for i := 0; i < eleType.NumField(); i++ {
						field := eleType.Field(i)
						realParamName := field.Tag.Get(`json`)
						properties[realParamName] = Yaml_Property{
							Example:     field.Tag.Get(`example`),
							Description: field.Tag.Get(`desc`),
						}
					}
				} else if eleKind == reflect.Map {
					map_ := route.Return.([]map[string]map[string]interface{})[0]
					for key, val1 := range map_ {
						properties[key] = Yaml_Property{
							Description: val1[`desc`].(string),
							Example:     val1[`example`],
						}
					}
				}
				succeedData.Items = map[string]interface{}{
					`properties`: properties,
				}
			} else if kind == reflect.Struct {
				properties := map[string]Yaml_Property{}
				paramsType := reflect.TypeOf(route.Return)
				for i := 0; i < paramsType.NumField(); i++ {
					field := paramsType.Field(i)
					fieldType := field.Type
					fieldKind := fieldType.Kind()
					if fieldKind == reflect.Struct {
						for i := 0; i < fieldType.NumField(); i++ {
							field := fieldType.Field(i)
							realParamName := field.Tag.Get(`json`)
							properties[realParamName] = Yaml_Property{
								Example:     field.Tag.Get(`example`),
								Description: field.Tag.Get(`desc`),
							}
						}
					} else {
						realParamName := field.Tag.Get(`json`)
						properties[realParamName] = Yaml_Property{
							Example:     field.Tag.Get(`example`),
							Description: field.Tag.Get(`desc`),
						}
					}
				}
				succeedData.Properties = properties
			}
		} else {
			succeedData.Type = this.getType(`object`)
			succeedData.Example = map[string]interface{}{}
		}
		succeedData.Description = `请求收到的内容`
		definitions[succeedDefinitionName] = Yaml_Definition{
			Type: `object`,
			Properties: map[string]Yaml_Property{
				`succeed`: {
					Type:        this.getType(`bool`),
					Example:     `true`,
					Description: `表示请求是否成功`,
				},
				`data`: succeedData,
			},
		}

		responses := map[string]Yaml_Response{
			`200`: {
				Description: `正确返回`,
				Schema: map[string]interface{}{
					`$ref`: fmt.Sprintf(`#/definitions/%s`, succeedDefinitionName),
				},
			},
			`400`: {
				Description: `错误返回`,
				Schema: map[string]interface{}{
					`$ref`: `#/definitions/Failed`,
				},
			},
		}

		paramType := `application/json`
		if route.ParamType != `` {
			paramType = route.ParamType
		}

		temp[strings.ToLower(route.Method)] = Yaml_Path{
			Tags:        []string{this.service.GetName()},
			Summary:     desc,
			Consumes:    []string{paramType},
			Produces:    []string{`application/json`},
			Parameters:  parameters,
			Responses:   responses,
			Description: description,
		}
		paths[this.service.GetPath()+route.Path] = temp
	}

	swagger := Yaml_Swagger{
		`2.0`,
		Yaml_Info{
			Title:       this.service.GetName(),
			Description: this.service.GetDescription(),
			Version:     `1.0.0`,
		},
		hostAndPort,
		this.service.GetPath(),
		[]Yaml_Tag{
			{
				Name:        this.service.GetName(),
				Description: this.service.GetDescription(),
			},
		},
		[]string{`http`},
		paths,
		definitions,
	}

	if type_ == `yaml` {
		bytes, _ := yaml.Marshal(&swagger)
		p_file.File.WriteFile(filename, bytes)
	} else if type_ == `json` {
		bytes, _ := json.Marshal(&swagger)
		p_file.File.WriteFile(filename, bytes)
	} else {
		panic(errors.New(`type 指定有误`))
	}

}
