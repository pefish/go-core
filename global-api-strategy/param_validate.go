package global_api_strategy

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-core/validator"
	go_desensitize "github.com/pefish/go-desensitize"
	go_format "github.com/pefish/go-format"
	i_core "github.com/pefish/go-interface/i-core"
	t_error "github.com/pefish/go-interface/t-error"
	go_json "github.com/pefish/go-json"
	go_string "github.com/pefish/go-string"
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

var ParamValidateStrategyInstance = NewParamValidateStrategy()

func NewParamValidateStrategy() *ParamValidateStrategy {
	return &ParamValidateStrategy{}
}

func (pvs *ParamValidateStrategy) Name() string {
	return `ParamValidateStrategy`
}

func (pvs *ParamValidateStrategy) Description() string {
	return `validate params`
}

func (pvs *ParamValidateStrategy) SetErrorCode(code uint64) i_core.IApiStrategy {
	pvs.errorCode = code
	return pvs
}

func (pvs *ParamValidateStrategy) SetErrorMsg(msg string) i_core.IApiStrategy {
	pvs.errorMsg = msg
	return pvs
}

func (pvs *ParamValidateStrategy) ErrorCode() uint64 {
	if pvs.errorCode == 0 {
		return t_error.INTERNAL_ERROR_CODE
	}
	return pvs.errorCode
}

func (pvs *ParamValidateStrategy) ErrorMsg() string {
	if pvs.errorMsg == "" {
		return "Params error."
	}
	return pvs.errorMsg
}

func (pvs *ParamValidateStrategy) recurValidate(out i_core.IApiSession, myValidator validator.ValidatorClass, map_ map[string]interface{}, globalValidator []string, type_ reflect.Type, value_ reflect.Value) (string, error) {
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
				newTag = strings.Join(globalValidator, ",") + "," + newTag
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
				value = go_format.ToString(map_[fieldName])
			case reflect.Uint64:
				if map_[fieldName] == nil {
					map_[fieldName] = 0
					defaultVal := typeField.Tag.Get(`default`)
					if defaultVal != `` {
						map_[fieldName] = defaultVal
					}
				}
				tmpValue, err := go_format.ToUint64(map_[fieldName])
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
				tmpValue, err := go_format.ToInt64(map_[fieldName])
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
				tmpValue, err := go_format.ToFloat64(map_[fieldName])
				if err != nil {
					return fieldName, errors.Errorf("ToFloat64 error - %#v", err)
				}
				value = tmpValue
			case reflect.Map:
				if map_[fieldName] == nil {
					value = nil
				} else {
					v, ok := map_[fieldName].(map[string]interface{})
					if !ok {
						return fieldName, errors.Errorf("Param <%#v> to map error", map_[fieldName])
					}
					value = v
				}
			default:
				return fieldName, errors.Errorf("Param kind error. fieldKind: %s", fieldKind.String())
			}
			out.Params()[fieldName] = value

			logger.LoggerDriverInstance.Logger.DebugF("[global_api_strategy.param_validate]: value: %#v, tag: %s", value, newTag)
			err := myValidator.Validator.Var(value, newTag)
			if err != nil {
				tempStr := go_string.StringUtilInstance.ReplaceAll(err.Error(), `for '' failed`, `for '`+fieldName+`' failed`)
				msg := go_string.StringUtilInstance.ReplaceAll(tempStr, `Key: ''`, `Key: '`+typeField.Name+`';`) + `; ` + newTag
				return fieldName, errors.Errorf(msg)
			}
		}
	}
	return "", nil
}

func (pvs *ParamValidateStrategy) Execute(out i_core.IApiSession) *t_error.ErrorInfo {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s trigger`, pvs.Name())
	if out.Api().Params() == nil {
		logger.LoggerDriverInstance.Logger.DebugF(`Params not be defined in api.`)
		return nil
	}

	myValidator := validator.ValidatorClass{}
	err := myValidator.Init()
	if err != nil {
		logger.LoggerDriverInstance.Logger.ErrorF(`validator init error`)
		return t_error.WrapWithAll(errors.New(pvs.ErrorMsg()), pvs.ErrorCode(), nil)
	}

	tempParam := make(map[string]interface{})

	if out.Method() == `GET` { // +号和%都有特殊含义，+会被替换成空格
		for k, v := range out.UrlParams() {
			tempParam[k] = v
		}
	} else if out.Method() == `POST` {
		requestContentType := out.Header(`content-type`)
		if out.Api().ParamType() != `` && !strings.HasPrefix(requestContentType, out.Api().ParamType()) {
			logger.LoggerDriverInstance.Logger.ErrorF(`content-type error`)
			return t_error.WrapWithAll(errors.New(pvs.ErrorMsg()), pvs.ErrorCode(), nil)
		}

		if strings.HasPrefix(requestContentType, MULTIPART_TYPE) && (out.Api().ParamType() == MULTIPART_TYPE || out.Api().ParamType() == ``) {
			formValues, err := out.FormValues()
			if err != nil {
				panic(err)
			}
			for k, v := range formValues {
				tempParam[k] = v[0]
			}
		} else if strings.HasPrefix(requestContentType, JSON_TYPE) && (out.Api().ParamType() == JSON_TYPE || out.Api().ParamType() == ``) {
			if err := out.ReadJSON(&tempParam); err != nil {
				logger.LoggerDriverInstance.Logger.ErrorF(`parse params error. err: %#v`, err)
				return t_error.WrapWithAll(errors.New(pvs.ErrorMsg()), pvs.ErrorCode(), nil)
			}
		} else if strings.HasPrefix(requestContentType, TEXT_TYPE) && (out.Api().ParamType() == TEXT_TYPE || out.Api().ParamType() == ``) {
			if err := out.ReadJSON(&tempParam); err != nil {
				logger.LoggerDriverInstance.Logger.ErrorF(`parse params error. err: %#v`, err)
				return t_error.WrapWithAll(errors.New(pvs.ErrorMsg()), pvs.ErrorCode(), nil)
			}
		} else {
			logger.LoggerDriverInstance.Logger.ErrorF(`content-type not be supported`)
			return t_error.WrapWithAll(errors.New(pvs.ErrorMsg()), pvs.ErrorCode(), nil)
		}
	} else {
		logger.LoggerDriverInstance.Logger.ErrorF(`scan params not be supported`)
		return t_error.WrapWithAll(errors.New(pvs.ErrorMsg()), pvs.ErrorCode(), nil)
	}
	// 深拷贝
	out.SetOriginalParams(go_json.Json.MustParseToMap(go_json.Json.MustStringify(tempParam)))
	out.SetParams(map[string]interface{}{})
	logger.LoggerDriverInstance.Logger.DebugF(`original params: %s`, go_desensitize.Desensitize.MustDesensitizeToString(out.OriginalParams()))

	globalValidator := make([]string, 0)
	fieldName, err := pvs.recurValidate(out, myValidator, tempParam, globalValidator, reflect.TypeOf(out.Api().Params()), reflect.ValueOf(out.Api().Params()))
	if err != nil {
		logger.LoggerDriverInstance.Logger.ErrorF(`Param validate error. - %#v`, err)
		return t_error.WrapWithAll(errors.New(pvs.ErrorMsg()), pvs.ErrorCode(), map[string]interface{}{
			`field`: fieldName,
		})
	}

	logger.LoggerDriverInstance.Logger.DebugF(`params: %s`, go_desensitize.Desensitize.MustDesensitizeToString(out.Params()))
	util.UpdateSessionErrorMsg(out, `params`, out.Params())

	return nil
}
