package api_session

import (
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
	"github.com/pefish/go-desensitize"
	"github.com/pefish/go-error"
	"github.com/pefish/go-format"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-string"
	"reflect"
	"strings"
)

type ApiHandlerType func(apiContext *ApiSessionClass) interface{}


type ApiSessionClass struct {
	Ctx iris.Context

	Validator *validator.Validate

	JwtHeaderName string
	JwtPayload    map[string]interface{}
	UserId        *uint64

	RouteName string

	Lang       string
	ClientType string // web、android、ios

	Options map[string]interface{}
}

func NewApiSession() *ApiSessionClass {
	return &ApiSessionClass{
		Options: map[string]interface{}{},
	}
}

func (this *ApiSessionClass) processGlobalValidators(fieldValue reflect.Value, globalValidator []string, oldTag string) string {
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
		result = p_string.String.RemoveLast(result, 1)
	}
	return result
}

func (this *ApiSessionClass) recurValidate(globalValidator []string, type_ reflect.Type, value_ reflect.Value) {
	for i := 0; i < value_.NumField(); i++ {
		typeField := type_.Field(i)
		typeFieldType := typeField.Type
		fieldKind := typeFieldType.Kind()
		fieldValue := value_.Field(i)
		if fieldKind == reflect.Struct {
			this.recurValidate(globalValidator, typeFieldType, fieldValue)
		} else {
			tagVal := typeField.Tag.Get(`validate`)
			newTag := tagVal
			if len(globalValidator) != 0 {
				newTag = this.processGlobalValidators(value_.Field(i), globalValidator, tagVal)
			}
			err := this.Validator.Var(fieldValue.Interface(), newTag)
			if err != nil {
				fieldName := typeField.Tag.Get(`json`)
				tempStr := p_string.String.ReplaceAll(err.Error(), `for '' failed`, `for '`+fieldName+`' failed`)
				p_error.ThrowErrorWithData(p_string.String.ReplaceAll(tempStr, `Key: ''`, `Key: '`+typeField.Name+`';`)+`; `+newTag, 0, map[string]interface{}{
					`field`: fieldName,
				}, err)
			}
		}
	}
}

func (this *ApiSessionClass) ScanParams(dest interface{}) {
	type_ := reflect.TypeOf(dest)
	if type_.Kind() != reflect.Ptr {
		p_error.ThrowInternal(`must be ptr`)
	}
	if type_.Elem().Kind() == reflect.Map {
		if this.Ctx.Method() == `GET` {
			p_format.Format.MapStringToStruct(dest, this.Ctx.URLParams())
		} else if this.Ctx.Method() == `POST` {
			if err := this.Ctx.ReadJSON(dest); err != nil {
				p_error.ThrowInternalError(`parse params error`, err)
			}
		} else {
			p_error.ThrowInternal(`scan params not support`)
		}
		this.logParams(dest)
	} else if type_.Elem().Kind() == reflect.Struct {
		if this.Ctx.Method() == `GET` {
			p_format.Format.MapStringToStruct(dest, this.Ctx.URLParams()) // +号和%都有特殊含义，+会被替换成空格
		} else if this.Ctx.Method() == `POST` {
			if err := this.Ctx.ReadJSON(dest); err != nil {
				p_error.ThrowInternalError(`parse params error`, err)
			}
		} else {
			p_error.ThrowInternal(`scan params not support`)
		}

		this.logParams(dest)
		if this.Validator != nil {
			glovalValdator := []string{`no-sql-inject`}
			type_ := reflect.TypeOf(dest).Elem()
			value_ := reflect.ValueOf(dest).Elem()
			valueKind := value_.Kind()
			if valueKind == reflect.Map {
				// validator不支持map
			} else {
				this.recurValidate(glovalValdator, type_, value_)
			}
		}
	} else {
		p_error.ThrowInternal(`ScanParams do not support this type`)
	}
}

func (this *ApiSessionClass) logParams(struct_ interface{}) {
	p_logger.Logger.Info(go_desensitize.Desensitize.DesensitizeToString(struct_))
}
