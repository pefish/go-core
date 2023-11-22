package external_service

import (
	go_test_ "github.com/pefish/go-test"
	"testing"
)

func TestExternalServiceDriver_Register(t *testing.T) {
	ExternalServiceDriverInstance.Register("go_test_", &TestExternalServiceInstance)

	ExternalServiceDriverInstance.Startup()

	svcs := ExternalServiceDriverInstance.ExternalServices()
	_, ok := svcs["go_test_"]
	go_test_.Equal(t, true, ok)
}
