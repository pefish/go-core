package global_api_strategy

import (
	"errors"
	"fmt"
	_type "github.com/pefish/go-core-type/api-session"
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
type ParamValidateStrategy struct {
	errorCode uint64
}

var ParamValidateStrategyInstance = ParamValidateStrategy{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

func (paramValidate *ParamValidateStrategy) GetName() string {
	return `paramValidate`
}

func (paramValidate *ParamValidateStrategy) GetDescription() string {
	return `validate params`
}

func (paramValidate *ParamValidateStrategy) SetErrorCode(code uint64) {
	paramValidate.errorCode = code
}

func (paramValidate *ParamValidateStrategy) GetErrorCode() uint64 {
	if paramValidate.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return paramValidate.errorCode
}

func (paramValidate *ParamValidateStrategy) processGlobalValidators(fieldValue reflect.Value, globalValidator []string, oldTag string) string {
	result := ``
	for _, validatorName := range globalValidator {
		if validatorName == validator.SQL_INJECT_CHECK && (strings.Contains(oldTag, validator.DISABLE_SQL_INJECT_CHECK) || fieldValue.Type().Kind() != reflect.String) {
			// 不是string类型 或者 有 DISABLE_SQL_INJECT_CHECK tag，就不校验 SQL_INJECT_CHECK
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

func (paramValidate *ParamValidateStrategy) recurValidate(out _type.IApiSession, myValidator validator.ValidatorClass, map_ map[string]interface{}, globalValidator []string, type_ reflect.Type, value_ reflect.Value) *go_error.ErrorInfo {
	for i := 0; i < value_.NumField(); i++ {
		typeField := type_.Field(i)
		typeFieldType := typeField.Type
		fieldKind := typeFieldType.Kind()
		fieldValue := value_.Field(i)
		if fieldKind == reflect.Struct {
			err := paramValidate.recurValidate(out, myValidator, map_, globalValidator, typeFieldType, fieldValue)
			if err != nil {
				return err
			}
		} else {
			tagVal := typeField.Tag.Get(`validate`)
			newTag := tagVal
			if len(globalValidator) != 0 {
				newTag = paramValidate.processGlobalValidators(value_.Field(i), globalValidator, tagVal)
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
						out.Params()[fieldName] = defaultVal
					}
				} else if strings.Contains(typeName, `int`) || strings.Contains(typeName, `float`) {
					map_[fieldName] = 0
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
						out.Params()[fieldName] = defaultVal
					}
				}
			}

			logger.LoggerDriverInstance.Logger.DebugF("[global_api_strategy.param_validate]: value: %#v, tag: %s", map_[fieldName], newTag)
			err := myValidator.Validator.Var(map_[fieldName], newTag)
			if err != nil {
				tempStr := go_string.String.ReplaceAll(err.Error(), `for '' failed`, `for '`+fieldName+`' failed`)
				msg := go_string.String.ReplaceAll(tempStr, `Key: ''`, `Key: '`+typeField.Name+`';`) + `; ` + newTag
				logger.LoggerDriverInstance.Logger.Error(msg)
				return go_error.WrapWithAll(
					fmt.Errorf("Params error."),
					paramValidate.errorCode, map[string]interface{}{
						`field`: fieldName,
					},
				)
			}
		}
	}
	return nil
}

func (paramValidate *ParamValidateStrategy) Init(param interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init`, paramValidate.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init defer`, paramValidate.GetName())
}

func (paramValidate *ParamValidateStrategy) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s trigger`, paramValidate.GetName())
	myValidator := validator.ValidatorClass{}
	err := myValidator.Init()
	if err != nil {
		return go_error.WrapWithAll(errors.New(`validator init error`), paramValidate.errorCode, nil)
	}

	tempParam := map[string]interface{}{}

	if out.Method() == `GET` { // +号和%都有特殊含义，+会被替换成空格
		for k, v := range out.UrlParams() {
			tempParam[k] = v
		}
	} else if out.Method() == `POST` {
		requestContentType := out.Header(`content-type`)
		if out.Api().GetParamType() != `` && !strings.HasPrefix(requestContentType, out.Api().GetParamType()) {
			return go_error.WrapWithAll(errors.New(`content-type error`), paramValidate.errorCode, nil)
		}

		if strings.HasPrefix(requestContentType, MULTIPART_TYPE) && (out.Api().GetParamType() == MULTIPART_TYPE || out.Api().GetParamType() == ``) {
			formValues, err := out.FormValues()
			if err != nil {
				panic(err)
			}
			for k, v := range formValues {
				tempParam[k] = v[0]
			}
		} else if strings.HasPrefix(requestContentType, JSON_TYPE) && (out.Api().GetParamType() == JSON_TYPE || out.Api().GetParamType() == ``) {
			if err := out.ReadJSON(&tempParam); err != nil {
				return go_error.WrapWithAll(fmt.Errorf(`parse params error. err: %#v`, err), paramValidate.errorCode, nil)
			}
		} else if strings.HasPrefix(requestContentType, TEXT_TYPE) && (out.Api().GetParamType() == TEXT_TYPE || out.Api().GetParamType() == ``) {
			if err := out.ReadJSON(&tempParam); err != nil {
				return go_error.WrapWithAll(fmt.Errorf(`parse params error. err: %#v`, err), paramValidate.errorCode, nil)
			}
		} else {
			return go_error.WrapWithAll(fmt.Errorf(`content-type not be supported`), paramValidate.errorCode, nil)
		}
	} else {
		return go_error.WrapWithAll(errors.New(`scan params not be supported`), paramValidate.errorCode, nil)
	}
	// 深拷贝
	out.SetOriginalParams(go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam)))
	out.SetParams(go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam)))
	paramsStr := go_desensitize.Desensitize.DesensitizeToString(tempParam)
	logger.LoggerDriverInstance.Logger.DebugF(`params: %s`, paramsStr)
	util.UpdateSessionErrorMsg(out, `params`, paramsStr)
	globalValidator := []string{validator.SQL_INJECT_CHECK}
	if out.Api().GetParams() != nil {
		err := paramValidate.recurValidate(out, myValidator, tempParam, globalValidator, reflect.TypeOf(out.Api().GetParams()), reflect.ValueOf(out.Api().GetParams()))
		if err != nil {
			return err
		}
	}

	return nil
}
