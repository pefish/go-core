package global_api_strategy

import (
	"errors"
	"fmt"
	api_session "github.com/pefish/go-core-type/api-session"
	global_api_strategy "github.com/pefish/go-core-type/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-core/validator"
	"github.com/pefish/go-desensitize"
	"github.com/pefish/go-error"
	"github.com/pefish/go-json"
	go_reflect "github.com/pefish/go-reflect"
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

type ParamValidateStrategy struct {
	errorCode uint64
	errorMsg  string
}

var ParamValidateStrategyInstance = ParamValidateStrategy{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

func (pvs *ParamValidateStrategy) GetName() string {
	return `ParamValidateStrategy`
}

func (pvs *ParamValidateStrategy) GetDescription() string {
	return `validate params`
}

func (pvs *ParamValidateStrategy) SetErrorCode(code uint64) global_api_strategy.IGlobalApiStrategy {
	pvs.errorCode = code
	return pvs
}

func (pvs *ParamValidateStrategy) SetErrorMsg(msg string) global_api_strategy.IGlobalApiStrategy {
	pvs.errorMsg = msg
	return pvs
}

func (pvs *ParamValidateStrategy) GetErrorCode() uint64 {
	if pvs.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return pvs.errorCode
}

func (pvs *ParamValidateStrategy) GetErrorMsg() string {
	if pvs.errorMsg == "" {
		return "Params error."
	}
	return pvs.errorMsg
}

func (pvs *ParamValidateStrategy) processGlobalValidators(fieldValue reflect.Value, globalValidator []string, oldTag string) string {
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

func (pvs *ParamValidateStrategy) recurValidate(out api_session.IApiSession, myValidator validator.ValidatorClass, map_ map[string]interface{}, globalValidator []string, type_ reflect.Type, value_ reflect.Value) (string, error) {
	logger.LoggerDriverInstance.Logger.DebugF("[global_api_strategy.param_validate]: map_: %#v", map_)

	for i := 0; i < value_.NumField(); i++ {
		typeField := type_.Field(i)
		typeFieldType := typeField.Type
		fieldKind := typeFieldType.Kind()
		fieldValue := value_.Field(i)
		if fieldKind == reflect.Struct {
			fieldName, err := pvs.recurValidate(out, myValidator, map_, globalValidator, typeFieldType, fieldValue)
			if err != nil {
				return fieldName, err
			}
		} else {
			tagVal := typeField.Tag.Get(`validate`)
			newTag := tagVal
			if len(globalValidator) != 0 {
				newTag = pvs.processGlobalValidators(value_.Field(i), globalValidator, tagVal)
			}
			jsonTag := typeField.Tag.Get(`json`)
			fieldName := strings.Split(jsonTag, `,`)[0]

			var value interface{}
			// correct type
			switch fieldKind {
			case reflect.String:
				if map_[fieldName] == nil {
					map_[fieldName] = ``
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
					}
				}
				value = go_reflect.Reflect.ToString(map_[fieldName])
			case reflect.Uint64:
				if map_[fieldName] == nil {
					map_[fieldName] = 0
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
					}
				}
				tmpValue, err := go_reflect.Reflect.ToUint64(map_[fieldName])
				if err != nil {
					return fieldName, fmt.Errorf("ToUint64 error - %#v", err)
				}
				value = tmpValue
			case reflect.Int64:
				if map_[fieldName] == nil {
					map_[fieldName] = 0
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
					}
				}
				tmpValue, err := go_reflect.Reflect.ToInt64(map_[fieldName])
				if err != nil {
					return fieldName, fmt.Errorf("ToInt64 error - %#v", err)
				}
				value = tmpValue
			case reflect.Float64:
				if map_[fieldName] == nil {
					map_[fieldName] = 0
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
					}
				}
				tmpValue, err := go_reflect.Reflect.ToFloat64(map_[fieldName])
				if err != nil {
					return fieldName, fmt.Errorf("ToFloat64 error - %#v", err)
				}
				value = tmpValue
			default:
				return fieldName, fmt.Errorf("param kind error. fieldKind: %#v", fieldKind)
			}
			out.Params()[fieldName] = value

			logger.LoggerDriverInstance.Logger.DebugF("[global_api_strategy.param_validate]: value: %#v, tag: %s", value, newTag)
			err := myValidator.Validator.Var(value, newTag)
			if err != nil {
				tempStr := go_string.String.ReplaceAll(err.Error(), `for '' failed`, `for '`+fieldName+`' failed`)
				msg := go_string.String.ReplaceAll(tempStr, `Key: ''`, `Key: '`+typeField.Name+`';`) + `; ` + newTag
				return fieldName, fmt.Errorf(msg)
			}
		}
	}
	return "", nil
}

func (pvs *ParamValidateStrategy) Init(param interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init`, pvs.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init defer`, pvs.GetName())
}

func (pvs *ParamValidateStrategy) Execute(out api_session.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s trigger`, pvs.GetName())
	myValidator := validator.ValidatorClass{}
	err := myValidator.Init()
	if err != nil {
		logger.LoggerDriverInstance.Logger.ErrorF(`validator init error`)
		return go_error.WrapWithAll(fmt.Errorf(pvs.GetErrorMsg()), pvs.GetErrorCode(), nil)
	}

	tempParam := map[string]interface{}{}

	if out.Method() == `GET` { // +号和%都有特殊含义，+会被替换成空格
		for k, v := range out.UrlParams() {
			tempParam[k] = v
		}
	} else if out.Method() == `POST` {
		requestContentType := out.Header(`content-type`)
		if out.Api().GetParamType() != `` && !strings.HasPrefix(requestContentType, out.Api().GetParamType()) {
			logger.LoggerDriverInstance.Logger.ErrorF(`content-type error`)
			return go_error.WrapWithAll(fmt.Errorf(pvs.GetErrorMsg()), pvs.GetErrorCode(), nil)
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
				logger.LoggerDriverInstance.Logger.ErrorF(`parse params error. err: %#v`, err)
				return go_error.WrapWithAll(fmt.Errorf(pvs.GetErrorMsg()), pvs.GetErrorCode(), nil)
			}
		} else if strings.HasPrefix(requestContentType, TEXT_TYPE) && (out.Api().GetParamType() == TEXT_TYPE || out.Api().GetParamType() == ``) {
			if err := out.ReadJSON(&tempParam); err != nil {
				logger.LoggerDriverInstance.Logger.ErrorF(`parse params error. err: %#v`, err)
				return go_error.WrapWithAll(fmt.Errorf(pvs.GetErrorMsg()), pvs.GetErrorCode(), nil)
			}
		} else {
			logger.LoggerDriverInstance.Logger.ErrorF(`content-type not be supported`)
			return go_error.WrapWithAll(fmt.Errorf(pvs.GetErrorMsg()), pvs.GetErrorCode(), nil)
		}
	} else {
		logger.LoggerDriverInstance.Logger.ErrorF(`scan params not be supported`)
		return go_error.WrapWithAll(errors.New(pvs.GetErrorMsg()), pvs.GetErrorCode(), nil)
	}
	// 深拷贝
	out.SetOriginalParams(go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam)))
	out.SetParams(map[string]interface{}{})
	logger.LoggerDriverInstance.Logger.DebugF(`original params: %s`, go_desensitize.Desensitize.DesensitizeToString(out.OriginalParams()))

	globalValidator := []string{validator.SQL_INJECT_CHECK}
	if out.Api().GetParams() != nil {
		fieldName, err := pvs.recurValidate(out, myValidator, tempParam, globalValidator, reflect.TypeOf(out.Api().GetParams()), reflect.ValueOf(out.Api().GetParams()))
		if err != nil {
			logger.LoggerDriverInstance.Logger.ErrorF(`param validate error. - %#v`, err)
			return go_error.WrapWithAll(errors.New(pvs.GetErrorMsg()), pvs.GetErrorCode(), map[string]interface{}{
				`field`: fieldName,
			})
		}
	}

	logger.LoggerDriverInstance.Logger.DebugF(`params: %s`, go_desensitize.Desensitize.DesensitizeToString(out.Params()))
	util.UpdateSessionErrorMsg(out, `params`, out.Params())

	return nil
}
