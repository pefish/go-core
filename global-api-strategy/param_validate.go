package global_api_strategy

import (
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-core/validator"
	"github.com/pefish/go-desensitize"
	"github.com/pefish/go-error"
	"github.com/pefish/go-json"
	"github.com/pefish/go-string"
	"reflect"
	"strings"
)

const (
	ALL_TYPE       = ``
	MULTIPART_TYPE = `multipart/form-data`
	JSON_TYPE      = `application/json`
	TEXT_TYPE      = `text/plain`
)

// 默认自带
type ParamValidateStrategyClass struct {
	errorCode uint64
}

var ParamValidateStrategy = ParamValidateStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

func (this *ParamValidateStrategyClass) GetName() string {
	return `paramValidate`
}

func (this *ParamValidateStrategyClass) GetDescription() string {
	return `validate params`
}

func (this *ParamValidateStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *ParamValidateStrategyClass) GetErrorCode() uint64 {
	if this.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return this.errorCode
}

func (this *ParamValidateStrategyClass) processGlobalValidators(fieldValue reflect.Value, globalValidator []string, oldTag string) string {
	result := ``
	for _, validatorName := range globalValidator {
		if validatorName == `no-sql-inject` && (strings.Contains(oldTag, `disable-inject-check`) || fieldValue.Type().Kind() != reflect.String) {
			// 不是string类型 或者 有disable-inject-check tag，就不校验no-sql-inject
			continue
		}
		result += validatorName + `,`
	}
	if oldTag != `` {
		result += oldTag
	} else if len(result) > 0 {
		result = go_string.String.RemoveLast(result, 1)
	}
	return result
}

func (this *ParamValidateStrategyClass) recurValidate(out *api_session.ApiSessionClass, myValidator validator.ValidatorClass, map_ map[string]interface{}, globalValidator []string, type_ reflect.Type, value_ reflect.Value) {
	for i := 0; i < value_.NumField(); i++ {
		typeField := type_.Field(i)
		typeFieldType := typeField.Type
		fieldKind := typeFieldType.Kind()
		fieldValue := value_.Field(i)
		if fieldKind == reflect.Struct {
			this.recurValidate(out, myValidator, map_, globalValidator, typeFieldType, fieldValue)
		} else {
			tagVal := typeField.Tag.Get(`validate`)
			newTag := tagVal
			if len(globalValidator) != 0 {
				newTag = this.processGlobalValidators(value_.Field(i), globalValidator, tagVal)
			}
			jsonTag := typeField.Tag.Get(`json`)
			fieldName := strings.Split(jsonTag, `,`)[0]
			if map_[fieldName] == nil { // map_[fieldName] 为nil的话，后面任何检查都不通过，不合理，所以这样处理
				typeName := typeField.Type.String()
				if typeName == `string` {
					map_[fieldName] = ``
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
						out.Params[fieldName] = defaultVal
					}
				} else if strings.Contains(typeName, `int`) || strings.Contains(typeName, `float`) {
					map_[fieldName] = 0
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
						out.Params[fieldName] = defaultVal
					}
				}
			}

			err := myValidator.Validator.Var(map_[fieldName], newTag)
			if err != nil {
				tempStr := go_string.String.ReplaceAll(err.Error(), `for '' failed`, `for '`+fieldName+`' failed`)
				go_error.ThrowErrorWithData(go_string.String.ReplaceAll(tempStr, `Key: ''`, `Key: '`+typeField.Name+`';`)+`; `+newTag, this.errorCode, map[string]interface{}{
					`field`: fieldName,
				}, err)
			}
		}
	}
}

func (this *ParamValidateStrategyClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())
}

func (this *ParamValidateStrategyClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, this.GetName())
	myValidator := validator.ValidatorClass{}
	myValidator.Init()

	tempParam := map[string]interface{}{}

	if out.GetMethod() == `GET` { // +号和%都有特殊含义，+会被替换成空格
		for k, v := range out.GetUrlParams() {
			tempParam[k] = v
		}
	} else if out.GetMethod() == `POST` {
		requestContentType := out.GetHeader(`content-type`)
		if out.Api.GetParamType() != `` && !strings.HasPrefix(requestContentType, out.Api.GetParamType()) {
			go_error.Throw(`content-type error`, this.errorCode)
		}

		if strings.HasPrefix(requestContentType, MULTIPART_TYPE) && (out.Api.GetParamType() == MULTIPART_TYPE || out.Api.GetParamType() == ``) {
			formValues, err := out.GetFormValues()
			if err != nil {
				panic(err)
			}
			for k, v := range formValues {
				tempParam[k] = v[0]
			}
		} else if strings.HasPrefix(requestContentType, JSON_TYPE) && (out.Api.GetParamType() == JSON_TYPE || out.Api.GetParamType() == ``) {
			if err := out.ReadJSON(&tempParam); err != nil {
				go_error.ThrowError(`parse params error`, this.errorCode, err)
			}
		} else {
			go_error.Throw(`content-type not be supported`, this.errorCode)
		}
	} else {
		go_error.Throw(`scan params not be supported`, this.errorCode)
	}
	// 深拷贝
	out.OriginalParams = go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam))
	out.Params = go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam))
	paramsStr := go_desensitize.Desensitize.DesensitizeToString(tempParam)
	logger.LoggerDriver.Logger.InfoF(`Params: %s`, paramsStr)
	util.UpdateSessionErrorMsg(out, `params`, paramsStr)
	glovalValdator := []string{`no-sql-inject`}
	if out.Api.GetParams() != nil {
		this.recurValidate(out, myValidator, tempParam, glovalValdator, reflect.TypeOf(out.Api.GetParams()), reflect.ValueOf(out.Api.GetParams()))
	}
}
