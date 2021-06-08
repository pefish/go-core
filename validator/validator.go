package validator

import (
	"errors"
	"github.com/go-playground/validator"
	"github.com/pefish/go-decimal"
	"github.com/pefish/go-reflect"
	"github.com/pefish/go-string"
	"reflect"
	"regexp"
	"strings"
)

type ValidatorClass struct {
	Validator *validator.Validate
}

const (
	SQL_INJECT_CHECK = "sql-inject-check"
	DISABLE_SQL_INJECT_CHECK = "disable-inject-check"
)

var Validator = ValidatorClass{}

func (validatorInstance *ValidatorClass) Init() error {
	validatorInstance.Validator = validator.New()
	err := validatorInstance.Validator.RegisterValidation(`is-mobile`, validatorInstance.Wrap(validatorInstance.IsMobile))
	err = validatorInstance.Validator.RegisterValidation(`contain-alphabet`, validatorInstance.Wrap(validatorInstance.ContainAlphabet))
	err = validatorInstance.Validator.RegisterValidation(`contain-number`, validatorInstance.Wrap(validatorInstance.ContainNumber))
	err = validatorInstance.Validator.RegisterValidation(`str-gte`, validatorInstance.Wrap(validatorInstance.StrGte))
	err = validatorInstance.Validator.RegisterValidation(`str-lte`, validatorInstance.Wrap(validatorInstance.StrLte))
	err = validatorInstance.Validator.RegisterValidation(`str-gt`, validatorInstance.Wrap(validatorInstance.StrGt))
	err = validatorInstance.Validator.RegisterValidation(`str-lt`, validatorInstance.Wrap(validatorInstance.StrLt))
	err = validatorInstance.Validator.RegisterValidation(`start-with`, validatorInstance.Wrap(validatorInstance.StartWith))
	err = validatorInstance.Validator.RegisterValidation(`end-with`, validatorInstance.Wrap(validatorInstance.EndWith))
	err = validatorInstance.Validator.RegisterValidation(SQL_INJECT_CHECK, validatorInstance.Wrap(validatorInstance.NoSqlInject))
	err = validatorInstance.Validator.RegisterValidation(DISABLE_SQL_INJECT_CHECK, validatorInstance.EmptyFun)
	if err != nil {
		return err
	}
	return nil
}

func (validatorInstance *ValidatorClass) EmptyFun(fl validator.FieldLevel) bool {
	return true
}

func (validatorInstance *ValidatorClass) Wrap(method func(val interface{}, target interface{}) bool) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		return method(fl.Field().Interface(), fl.Param())
	}
}

func (validatorInstance *ValidatorClass) IsMobile(val interface{}, target interface{}) bool {
	return regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`).MatchString(val.(string))
}

func (validatorInstance *ValidatorClass) ContainAlphabet(val interface{}, target interface{}) bool {
	str := go_reflect.Reflect.ToString(val)
	allAlphabet := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	for _, charInt := range str {
		if strings.Contains(allAlphabet, string(charInt)) {
			return true
		}
	}
	return false
}

func (validatorInstance *ValidatorClass) ContainNumber(val interface{}, target interface{}) bool {
	str := go_reflect.Reflect.ToString(val)
	allNumbers := `0123456789`
	for _, charInt := range str {
		if strings.Contains(allNumbers, string(charInt)) {
			return true
		}
	}
	return false
}

func (validatorInstance *ValidatorClass) StrGte(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.ToString(val)).Gte(target)
}

func (validatorInstance *ValidatorClass) StrLte(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.ToString(val)).Lte(target)
}

func (validatorInstance *ValidatorClass) StrGt(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.ToString(val)).Gt(target)
}

func (validatorInstance *ValidatorClass) StrLt(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.ToString(val)).Lt(target)
}

func (validatorInstance *ValidatorClass) StartWith(val interface{}, target interface{}) bool {
	return go_string.String.StartWith(go_reflect.Reflect.ToString(val), go_reflect.Reflect.ToString(target))
}

func (validatorInstance *ValidatorClass) EndWith(val interface{}, target interface{}) bool {
	return go_string.String.EndWith(go_reflect.Reflect.ToString(val), go_reflect.Reflect.ToString(target))
}

func (validatorInstance *ValidatorClass) NoSqlInject(val interface{}, target interface{}) bool {
	if reflect.TypeOf(val).Kind() != reflect.String {
		return true
	}
	err := validatorInstance.checkInjectWithErr(go_reflect.Reflect.ToString(val))
	return err == nil
}

func (validatorInstance *ValidatorClass) checkInjectWithErr(str string) error {
	arr := []string{
		`=`, `{`, `}`, `;`, `|`, `>`, `<`, `"`, `[`, `]`, `\`, `/`, `?`, `%`, `1 = 1`, `1=1`, `1 =1`, `1= 1`,
	}
	for _, char := range arr {
		if strings.Contains(str, char) {
			return errors.New(`inject error`)
		}
	}
	return nil
}
