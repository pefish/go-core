package validator

import (
	"github.com/go-playground/validator"
	"github.com/pefish/go-decimal"
	go_format "github.com/pefish/go-format"
	"github.com/pefish/go-string"
	"regexp"
	"strings"
)

type ValidatorClass struct {
	Validator *validator.Validate
}

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
	str := go_format.FormatInstance.ToString(val)
	allAlphabet := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	for _, charInt := range str {
		if strings.Contains(allAlphabet, string(charInt)) {
			return true
		}
	}
	return false
}

func (validatorInstance *ValidatorClass) ContainNumber(val interface{}, target interface{}) bool {
	str := go_format.FormatInstance.ToString(val)
	allNumbers := `0123456789`
	for _, charInt := range str {
		if strings.Contains(allNumbers, string(charInt)) {
			return true
		}
	}
	return false
}

func (validatorInstance *ValidatorClass) StrGte(val interface{}, target interface{}) bool {
	dc, err := go_decimal.Decimal.Start(val)
	if err != nil {
		return false
	}
	result, err := dc.Gte(target)
	if err != nil {
		return false
	}
	return result
}

func (validatorInstance *ValidatorClass) StrLte(val interface{}, target interface{}) bool {
	dc, err := go_decimal.Decimal.Start(val)
	if err != nil {
		return false
	}
	result, err := dc.Lte(target)
	if err != nil {
		return false
	}
	return result
}

func (validatorInstance *ValidatorClass) StrGt(val interface{}, target interface{}) bool {
	dc, err := go_decimal.Decimal.Start(val)
	if err != nil {
		return false
	}
	result, err := dc.Gt(target)
	if err != nil {
		return false
	}
	return result
}

func (validatorInstance *ValidatorClass) StrLt(val interface{}, target interface{}) bool {
	dc, err := go_decimal.Decimal.Start(val)
	if err != nil {
		return false
	}
	result, err := dc.Lt(target)
	if err != nil {
		return false
	}
	return result
}

func (validatorInstance *ValidatorClass) StartWith(val interface{}, target interface{}) bool {
	return go_string.StringUtilInstance.StartWith(go_format.FormatInstance.ToString(val), go_format.FormatInstance.ToString(target))
}

func (validatorInstance *ValidatorClass) EndWith(val interface{}, target interface{}) bool {
	return go_string.StringUtilInstance.EndWith(go_format.FormatInstance.ToString(val), go_format.FormatInstance.ToString(target))
}
