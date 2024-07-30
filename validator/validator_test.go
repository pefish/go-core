package validator

import (
	go_test_ "github.com/pefish/go-test"
	"testing"
)

func TestValidatorClass_Init(t *testing.T) {
	err := Validator.Init()
	go_test_.Equal(t, nil, err)
}

func TestValidatorClass_IsMobile(t *testing.T) {
	result := Validator.IsMobile("17299647732", nil)
	go_test_.Equal(t, true, result)

	result = Validator.IsMobile("1729964773a", nil)
	go_test_.Equal(t, false, result)
}

func TestValidatorClass_ContainAlphabet(t *testing.T) {
	result := Validator.ContainAlphabet("17299647732", nil)
	go_test_.Equal(t, false, result)

	result = Validator.ContainAlphabet("1729964773a", nil)
	go_test_.Equal(t, true, result)
}

func TestValidatorClass_ContainNumber(t *testing.T) {
	result := Validator.ContainNumber("17299647732", nil)
	go_test_.Equal(t, true, result)

	result = Validator.ContainNumber("sfdafha", nil)
	go_test_.Equal(t, false, result)
}

func TestValidatorClass_StrGte(t *testing.T) {
	result := Validator.StrGte("123", 123)
	go_test_.Equal(t, true, result)

	result = Validator.StrGte("123", 124)
	go_test_.Equal(t, false, result)

	result = Validator.StrGte("123", 121)
	go_test_.Equal(t, true, result)
}

func TestValidatorClass_StrLte(t *testing.T) {
	result := Validator.StrLte("123", 123)
	go_test_.Equal(t, true, result)

	result = Validator.StrLte("123", 124)
	go_test_.Equal(t, true, result)

	result = Validator.StrLte("123", 121)
	go_test_.Equal(t, false, result)
}

func TestValidatorClass_StrGt(t *testing.T) {
	result := Validator.StrGt("123", 123)
	go_test_.Equal(t, false, result)

	result = Validator.StrGt("123", 124)
	go_test_.Equal(t, false, result)

	result = Validator.StrGt("123", 121)
	go_test_.Equal(t, true, result)
}

func TestValidatorClass_StrLt(t *testing.T) {
	result := Validator.StrLt("123", 123)
	go_test_.Equal(t, false, result)

	result = Validator.StrLt("123", 124)
	go_test_.Equal(t, true, result)

	result = Validator.StrLt("123", 121)
	go_test_.Equal(t, false, result)
}

func TestValidatorClass_StartWith(t *testing.T) {
	result := Validator.StartWith("abc635462yh", "abc")
	go_test_.Equal(t, true, result)

	result = Validator.StartWith("bdjyetsrdn", "abc")
	go_test_.Equal(t, false, result)
}

func TestValidatorClass_EndWith(t *testing.T) {
	result := Validator.EndWith("635462yhabc", "abc")
	go_test_.Equal(t, true, result)

	result = Validator.EndWith("bdjyetsrdn", "abc")
	go_test_.Equal(t, false, result)
}
