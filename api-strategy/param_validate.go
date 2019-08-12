package api_strategy

import (
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/validator"
	"github.com/pefish/go-desensitize"
	"github.com/pefish/go-error"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-string"
	"reflect"
	"strings"
)

// 默认自带
type ParamValidateStrategyClass struct {
	errorCode uint64
	Validator validator.ValidatorClass
}

type ParamValidateParam struct {
	Param interface{}
}

var ParamValidateApiStrategy = ParamValidateStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
	Validator: validator.ValidatorClass{},
}

func (this *ParamValidateStrategyClass) GetName() string {
	return `paramValidate`
}

func (this *ParamValidateStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
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

func (this *ParamValidateStrategyClass) recurValidate(map_ map[string]interface{}, globalValidator []string, type_ reflect.Type, value_ reflect.Value) {
	for i := 0; i < value_.NumField(); i++ {
		typeField := type_.Field(i)
		typeFieldType := typeField.Type
		fieldKind := typeFieldType.Kind()
		fieldValue := value_.Field(i)
		if fieldKind == reflect.Struct {
			this.recurValidate(map_, globalValidator, typeFieldType, fieldValue)
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
				} else if strings.Contains(typeName, `int`) || strings.Contains(typeName, `float`) {
					map_[fieldName] = 0
				}
			}

			err := this.Validator.Validator.Var(map_[fieldName], newTag)
			if err != nil {
				tempStr := go_string.String.ReplaceAll(err.Error(), `for '' failed`, `for '`+fieldName+`' failed`)
				go_error.ThrowErrorWithData(go_string.String.ReplaceAll(tempStr, `Key: ''`, `Key: '`+typeField.Name+`';`)+`; `+newTag, this.errorCode, map[string]interface{}{
					`field`: fieldName,
				}, err)
			}
		}
	}
}

func (this *ParamValidateStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	newParam := param.(ParamValidateParam)
	this.Validator.Init()

	tempParam := map[string]interface{}{}

	if ctx.Method() == `GET` { // +号和%都有特殊含义，+会被替换成空格
		for k, v := range ctx.URLParams() {
			tempParam[k] = v
		}
	} else if ctx.Method() == `POST` {
		requestContentType := ctx.GetHeader(`content-type`)
		if strings.HasPrefix(requestContentType, `application/json`) {
			if err := ctx.ReadJSON(&tempParam); err != nil {
				go_error.ThrowError(`parse params error`, this.errorCode, err)
			}
		} else if strings.HasPrefix(requestContentType, `multipart/form-data`) {
			for k, v := range ctx.FormValues() {
				tempParam[k] = v[0]
			}
		} else {
			go_error.Throw(`content-type not be supported`, this.errorCode)
		}
	} else {
		go_error.Throw(`scan params not be supported`, this.errorCode)
	}
	out.Params = tempParam
	go_logger.Logger.Info(go_desensitize.Desensitize.DesensitizeToString(tempParam))
	glovalValdator := []string{`no-sql-inject`}
	if newParam.Param != nil {
		this.recurValidate(tempParam, glovalValdator, reflect.TypeOf(newParam.Param), reflect.ValueOf(newParam.Param))
	}
}
