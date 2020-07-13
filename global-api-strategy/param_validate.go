package global_api_strategy

import (
	_type "github.com/pefish/go-core/api-session/type"
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

func (paramValidate *ParamValidateStrategyClass) GetName() string {
	return `paramValidate`
}

func (paramValidate *ParamValidateStrategyClass) GetDescription() string {
	return `validate params`
}

func (paramValidate *ParamValidateStrategyClass) SetErrorCode(code uint64) {
	paramValidate.errorCode = code
}

func (paramValidate *ParamValidateStrategyClass) GetErrorCode() uint64 {
	if paramValidate.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return paramValidate.errorCode
}

func (paramValidate *ParamValidateStrategyClass) processGlobalValidators(fieldValue reflect.Value, globalValidator []string, oldTag string) string {
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

func (paramValidate *ParamValidateStrategyClass) recurValidate(out _type.IApiSession, myValidator validator.ValidatorClass, map_ map[string]interface{}, globalValidator []string, type_ reflect.Type, value_ reflect.Value) *go_error.ErrorInfo {
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

			err := myValidator.Validator.Var(map_[fieldName], newTag)
			if err != nil {
				tempStr := go_string.String.ReplaceAll(err.Error(), `for '' failed`, `for '`+fieldName+`' failed`)
				msg := go_string.String.ReplaceAll(tempStr, `Key: ''`, `Key: '`+typeField.Name+`';`)+`; `+newTag
				return &go_error.ErrorInfo{
					InternalErrorMessage: msg,
					ErrorMessage: msg,
					ErrorCode: paramValidate.errorCode,
					Data: map[string]interface{}{
						`field`: fieldName,
					},
					Err: err,
				}
			}
		}
	}
	return nil
}

func (paramValidate *ParamValidateStrategyClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, paramValidate.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, paramValidate.GetName())
}

func (paramValidate *ParamValidateStrategyClass) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, paramValidate.GetName())
	myValidator := validator.ValidatorClass{}
	err := myValidator.Init()
	if err != nil {
		return &go_error.ErrorInfo{
			InternalErrorMessage: `validator init error`,
			ErrorMessage: go_error.INTERNAL_ERROR,
			ErrorCode: paramValidate.errorCode,
		}
	}

	tempParam := map[string]interface{}{}

	if out.Method() == `GET` { // +号和%都有特殊含义，+会被替换成空格
		for k, v := range out.UrlParams() {
			tempParam[k] = v
		}
	} else if out.Method() == `POST` {
		requestContentType := out.Header(`content-type`)
		if out.Api().GetParamType() != `` && !strings.HasPrefix(requestContentType, out.Api().GetParamType()) {
			return &go_error.ErrorInfo{
				InternalErrorMessage: `content-type error`,
				ErrorMessage: `content-type error`,
				ErrorCode: paramValidate.errorCode,
			}
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
				return &go_error.ErrorInfo{
					InternalErrorMessage: `parse params error`,
					ErrorMessage: `parse params error`,
					ErrorCode: paramValidate.errorCode,
				}
			}
		} else {
			return &go_error.ErrorInfo{
				InternalErrorMessage: `content-type not be supported`,
				ErrorMessage: `content-type not be supported`,
				ErrorCode: paramValidate.errorCode,
			}
		}
	} else {
		return &go_error.ErrorInfo{
			InternalErrorMessage: `scan params not be supported`,
			ErrorMessage: `scan params not be supported`,
			ErrorCode: paramValidate.errorCode,
		}
	}
	// 深拷贝
	out.SetOriginalParams(go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam)))
	out.SetParams(go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam)))
	paramsStr := go_desensitize.Desensitize.DesensitizeToString(tempParam)
	logger.LoggerDriver.Logger.InfoF(`params: %s`, paramsStr)
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
