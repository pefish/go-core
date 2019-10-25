package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/pefish/go-decimal"
	"github.com/pefish/go-error"
	"github.com/pefish/go-reflect"
	"github.com/pefish/go-string"
	"reflect"
	"regexp"
	"strings"
)

type ValidatorClass struct {
	Validator *validator.Validate
}

var Validator = ValidatorClass{}

func (this *ValidatorClass) Init() {
	this.Validator = validator.New()
	this.Validator.RegisterAlias(`is-email`, `email`) // 有bug，对|运算符没用
	this.Validator.RegisterAlias(`max-length`, `max`)
	this.Validator.RegisterAlias(`min-length`, `min`)
	err := this.Validator.RegisterValidation(`is-mobile`, this.Wrap(this.IsMobile))
	err = this.Validator.RegisterValidation(`contain-alphabet`, this.Wrap(this.ContainAlphabet))
	err = this.Validator.RegisterValidation(`contain-number`, this.Wrap(this.ContainNumber))
	err = this.Validator.RegisterValidation(`str-gte`, this.Wrap(this.StrGte))
	err = this.Validator.RegisterValidation(`str-lte`, this.Wrap(this.StrLte))
	err = this.Validator.RegisterValidation(`str-gt`, this.Wrap(this.StrGt))
	err = this.Validator.RegisterValidation(`str-lt`, this.Wrap(this.StrLt))
	err = this.Validator.RegisterValidation(`start-with`, this.Wrap(this.StartWith))
	err = this.Validator.RegisterValidation(`end-with`, this.Wrap(this.EndWith))
	err = this.Validator.RegisterValidation(`no-sql-inject`, this.Wrap(this.NoSqlInject))
	err = this.Validator.RegisterValidation(`disable-inject-check`, this.EmptyFun)
	err = this.Validator.RegisterValidation(`test`, this.Test)
	if err != nil {
		go_error.ThrowInternal(`validator init error`)
	}
}

func (this *ValidatorClass) Test(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()
	targetField, targetKind, ok := fl.GetStructFieldOK() // 根据指定的字段名获取那个字段
	fmt.Println(`targetField`, targetField, targetKind)
	if !ok || targetKind != kind {
		return false
	}
	return false
}

func (this *ValidatorClass) EmptyFun(fl validator.FieldLevel) bool {
	return true
}

func (this *ValidatorClass) Wrap(method func(val interface{}, target interface{}) bool) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		return method(fl.Field().Interface(), fl.Param())
	}
}

func (this *ValidatorClass) IsMobile(val interface{}, target interface{}) bool {
	return regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`).MatchString(val.(string))
}

func (this *ValidatorClass) ContainAlphabet(val interface{}, target interface{}) bool {
	str := go_reflect.Reflect.MustToString(val)
	allAlphabet := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	for _, charInt := range str {
		if strings.Contains(allAlphabet, string(charInt)) {
			return true
		}
	}
	return false
}

func (this *ValidatorClass) ContainNumber(val interface{}, target interface{}) bool {
	str := go_reflect.Reflect.MustToString(val)
	allNumbers := `0123456789`
	for _, charInt := range str {
		if strings.Contains(allNumbers, string(charInt)) {
			return true
		}
	}
	return false
}

func (this *ValidatorClass) StrGte(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.MustToString(val)).Gte(target)
}

func (this *ValidatorClass) StrLte(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.MustToString(val)).Lte(target)
}

func (this *ValidatorClass) StrGt(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.MustToString(val)).Gt(target)
}

func (this *ValidatorClass) StrLt(val interface{}, target interface{}) bool {
	return go_decimal.Decimal.Start(go_reflect.Reflect.MustToString(val)).Lt(target)
}

func (this *ValidatorClass) StartWith(val interface{}, target interface{}) bool {
	return go_string.String.StartWith(go_reflect.Reflect.MustToString(val), go_reflect.Reflect.MustToString(target))
}

func (this *ValidatorClass) EndWith(val interface{}, target interface{}) bool {
	return go_string.String.EndWith(go_reflect.Reflect.MustToString(val), go_reflect.Reflect.MustToString(target))
}

func (this *ValidatorClass) NoSqlInject(val interface{}, target interface{}) bool {
	if reflect.TypeOf(val).Kind() != reflect.String {
		return true
	}
	err := this.CheckInjectWithErr(go_reflect.Reflect.MustToString(val))
	return err == nil
}

func (this *ValidatorClass) CheckInject(str string) {
	err := this.CheckInjectWithErr(str)
	if err != nil {
		go_error.ThrowInternal(`inject error`)
	}
}

func (this *ValidatorClass) CheckInjectWithErr(str string) error {
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
